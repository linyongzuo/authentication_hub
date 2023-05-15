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

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./conf/home.html")
}

func Start() {
	// 注册回调
	go hub.run()
	address := flag.String("addr", fmt.Sprintf("0.0.0.0:%d", global.Cfg.ServerCfg.Port), "http service address")
	http.HandleFunc("/authentication_hub/connect", func(w http.ResponseWriter, r *http.Request) {
		connect(hub, w, r)
	})
	http.HandleFunc("/", serveHome)
	err := http.ListenAndServe(*address, nil)
	if err != nil {
		panic(err)
	}
}
