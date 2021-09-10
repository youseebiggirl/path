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
	Receive(Hub) error      // 从 conn 中读取用户输入的消息，并发送到 Hub 中，进行后续处理
	Logout()           // 用户要退出时，调用该方法
	Server(hub Hub) error   // 为该 user 开启聊天服务
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
	case *groupChatMessage:
		hub.GroupMsgChan() <- msg
	}
}

// Receive 从 conn 中读取用户输入的消息，并发送到 Hub 中，让 Hub 进行后续处理
func (u *user) Receive(hub Hub) error {
	pack := NewDataPack()
	for {
		_, msg, err := u.conn.ReadMessage()
		if err != nil {
			return err
		}

		m, err := pack.UnPack(msg)
		if err != nil {
			return err
		}
		log.Println(m)

		// 根据消息的类型，发送到对应的 chan 上
		switch m.MsgType() {
		case PrivateMsg:
			hub.PrivateMsgChan() <- m
		case GroupMsg:
			hub.GroupMsgChan() <- m
		}
	}
}

func (u *user) Server(hub Hub) (err error) {
	// 上线信息
	hub.OnlineChan() <- u

	go func() {
		err = u.Receive(hub)
	}()

	return
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
