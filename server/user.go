package server

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"time"
)

type User interface {
	Id() uint64   // 用户的唯一标识
	Online() bool // 用户当前是否在线
	SetOnline(bool)

	RemoteAddr() string
	Name() string
	Conn() *websocket.Conn
	MsgChan() chan Message // 缓存用户接收到的信息

	Send(Hub, Message) // 将 Message 发送到 Hub
	Receive(Hub)       // 从 conn 中读取消息，并发送到 Hub 中
	Logout()           // 用户要退出时，调用该方法
	ServerWS(hub Hub)  //
}

type user struct {
	conn       *websocket.Conn
	remoteAddr string
	id         uint64
	name       string
	online     bool
	msgChan    chan Message
}

func NewUser(conn *websocket.Conn, addr string, name string) *user {
	rand.Seed(time.Now().Unix())

	return &user{
		remoteAddr: addr,
		// TODO id 应该唯一标识而不是随机，为此应该存储在数据库中
		id:      rand.Uint64(),
		name:    name,
		conn:    conn,
		msgChan: make(chan Message, DefaultChanSize),
	}
}

func (u *user) Logout() {

}

func (u *user) Send(hub Hub, msg Message) {
	switch msg.(type) {
	case *privateChatMessage:
		hub.PrivateMsgChan() <- msg
		//// 如果用户当前不在线
		//if !msg.To().user.Online() {
		//	// TODO 暂时保存到队列中
		//	return
		//}
		//if err := u.conn.WriteMessage(msg.Type(), msg.Data()); err != nil {
		//	log.Println(err)
		//	return
		//}
	case *groupChatMessage:
		hub.GroupMsgChan() <- msg
		//for _, conn := range hub.Users() {
		//	if err := conn.WriteMessage(msg.Type(), msg.Data()); err != nil {
		//		log.Println(err)
		//		continue
		//	}
		//}
	}
}

func (u *user) Receive(hub Hub) {
	for {
		_, msg, err := u.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		hub.GroupMsgChan() <- NewGroupChatMessage(u.id, 0, 0/*FIXME*/, msg)
	}
}

func (u *user) ServerWS(hub Hub) {
	// 上线信息
	hub.OnlineChan() <- u
	log.Printf("%p \n", hub)
	go u.Receive(hub)
}

func (u *user) Online() bool {
	return u.online
}

func (u *user) Id() uint64 {
	return u.id
}

func (u *user) SetOnline(b bool) {
	u.online = b
}

func (u *user) RemoteAddr() string {
	return u.remoteAddr
}

func (u *user) Name() string {
	return u.name
}

func (u *user) Conn() *websocket.Conn {
	return u.conn
}

func (u *user) MsgChan() chan Message {
	return u.msgChan
}
