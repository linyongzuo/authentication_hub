package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/authentication_hub/global"
	"github.com/authentication_hub/internal/constants"
	"github.com/authentication_hub/internal/domain/request"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	heartbeatTime = 10 * time.Second
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Buffered channel of outbound messages.
	receive chan []byte

	offlineChan chan []byte

	stopChan chan struct{}

	online atomic.Bool

	ip string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.offline()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		message := make([]byte, 0)
		_, message, err = c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.receive <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.offline()
	}()
	for {
		select {
		case message, ok := <-c.receive:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if !ok {
				// The hub closed the channel.
				err = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(handlerMessage(message))
			if err != nil {
				return
			}

			if err = w.Close(); err != nil {
				return
			}
		case message, ok := <-c.offlineChan:
			{
				if !ok {
					logrus.Info("读取离线消息失败")
					return
				}
				handlerMessage(message)
				if c.online.Load() {
					err := c.conn.Close()
					if err != nil {
						logrus.Errorf("关闭链接失败")
					}
					logrus.Info("关闭当前链接")
					c.stopChan <- struct{}{}
					c.hub.unregister <- c
					c.online.Swap(false)
				}
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err = c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) heartbeatCheck() {
	ticker := time.NewTicker(heartbeatTime)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			logrus.Info("检测心跳,检测地址:%s", c.ip)
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()
			// 按前缀扫描key
			stringCmd := global.Rdb.HGet(ctx, c.ip, constants.KActiveTime)
			lastTime, err := stringCmd.Result()
			if err != nil {
				continue
			}
			activeTime, err := time.ParseInLocation(constants.KTimeTemplate, lastTime, time.Local)
			if err != nil {
				logrus.Errorf("解析时间出错:%s", err.Error())
				continue
			}
			now := time.Now()
			logrus.Infof("当前时间:%s,最后一次心跳时间:%s", now.Format(constants.KTimeTemplate), lastTime)
			if now.Sub(activeTime).Seconds() > 60 {
				logrus.Info("客户端掉线")
				c.offline()
				return
			}

		}
	}
}
func (c *Client) offline() {
	if !c.online.Load() {
		return
	}
	logrus.Infof("客户端下线:%s", c.ip)
	req := request.UserLogoutReq{
		Header: request.Header{
			Version:     "",
			MessageType: request.MessageUserLogout,
		},
		Mac:    "",
		Ip:     c.ip,
		System: true,
	}
	message, _ := json.Marshal(req)
	c.offlineChan <- message
}

// serveWs handles websocket requests from the peer.
func connect(hub *Hub, w http.ResponseWriter, r *http.Request) {

	conn, err := global.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	ip := conn.RemoteAddr().String()
	address := strings.Split(ip, ":")
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	global.Rdb.HSet(ctx, address[0], constants.KActiveTime, time.Now().Format(constants.KTimeTemplate))
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 1024), receive: make(chan []byte, 1024), stopChan: make(chan struct{}), offlineChan: make(chan []byte, 1024), online: atomic.Bool{}, ip: address[0]}
	client.online.Store(true)
	logrus.Infof("有客户端连接:%s", address[0])
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.heartbeatCheck()
	go client.readPump()
}
