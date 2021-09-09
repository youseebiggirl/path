package main

import (
	"chat/server"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

// error:
// 使用 http.ServeFile 会产生如下错误：
// websocket: the client is not using the websocket protocol:
// 'upgrade' token not found in 'Connection' header
func main() {
	hub := server.NewHub()
	log.Printf("%p \n", hub)
	go hub.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//http.ServeFile(w, r, "./html/home.html")
		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			//return
		}
		//log.Printf("%p \n", &hub)
		user := server.NewUser(conn, r.RemoteAddr, "test")
		user.ServerWS(hub)
	})

	//http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
	//	if r.Method != "GET" {
	//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//		return
	//	}
	//	http.ServeFile(w, r, "./html/home.html")
	//})

	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalln(err)
	}
}
