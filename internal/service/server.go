package service

import (
	"flag"
	"fmt"
	"github.com/authentication_hub/global"
	"github.com/authentication_hub/internal/controller"
	"net/http"
)

var hub = newHub()
var ctrl = controller.New()

func Start() {
	// 注册回调
	go hub.run()
	address := flag.String("addr", fmt.Sprintf("0.0.0.0:%d", global.Cfg.ServerCfg.Port), "http service address")
	http.HandleFunc("/authentication_hub/connect", func(w http.ResponseWriter, r *http.Request) {
		connect(hub, w, r)
	})
	err := http.ListenAndServe(*address, nil)
	if err != nil {
		panic(err)
	}
}
