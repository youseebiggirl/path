package server

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

const (
	DefaultChanSize = 100
)

type Hub interface {
	OnlineChan() chan User        // 缓存所有用户的上线信息
	LogoutChan() chan User        // 缓存所有用户的下线信息
	PrivateMsgChan() chan Message // 缓存所有用户要发送的信息
	GroupMsgChan() chan Message   //

	Users() map[uint64]*websocket.Conn // 保存当前在线的所有用户

	HandlerOnline(User)        // 当有用户上线时，相应的处理方法
	HandlerLogout(User)        // 当有用户下线时，相应的处理方法
	HandlerPrivateMsg(Message) // 处理用户发送的私聊信息
	HandlerGroupMsg(Message)   // 处理用户发送的群聊消息
	Run()
}

type hub struct {
	mu             sync.RWMutex
	privateMsgChan chan Message // 缓存所有用户要发送的信息
	groupMsgChan   chan Message
	onlineChan     chan User
	logoutChan     chan User
	users          map[uint64]*websocket.Conn
}

func NewHub() *hub {
	return &hub{
		privateMsgChan: make(chan Message, DefaultChanSize),
		groupMsgChan:   make(chan Message, 0),
		onlineChan:     make(chan User, DefaultChanSize),
		logoutChan:     make(chan User, DefaultChanSize),
		users:          make(map[uint64]*websocket.Conn),
	}
}

func (h *hub) HandlerPrivateMsg(msg Message) {

}

// HandlerGroupMsg 处理用户发送的群聊消息
func (h *hub) HandlerGroupMsg(msg Message) {
	// 发送给全部用户
	// TODO 发送给指定群用户
	for _, conn := range h.users {
		if err := conn.WriteMessage(websocket.BinaryMessage, msg.Data()); err != nil {
			log.Println(err)
			continue
		}
	}
}

func (h *hub) HandlerOnline(usr User) {
	log.Printf("%+v online \n", usr)
	h.mu.Lock()
	h.users[usr.Id()] = usr.Conn()
	usr.SetOnline(true)
	h.mu.Unlock()
}

func (h *hub) HandlerLogout(usr User) {
	h.mu.Lock()
	delete(h.users, usr.Id())
	usr.SetOnline(false)
	h.mu.Unlock()
}

func (h *hub) Run() {
	for {
		select {
		case msg := <-h.privateMsgChan:
			h.HandlerPrivateMsg(msg)
		case msg := <-h.groupMsgChan:
			h.HandlerGroupMsg(msg)
		case usr := <-h.onlineChan:
			h.HandlerOnline(usr)
		case usr := <-h.logoutChan:
			h.HandlerLogout(usr)
		}
	}
}

func (h *hub) OnlineChan() chan User {
	return h.onlineChan
}

func (h *hub) LogoutChan() chan User {
	return h.logoutChan
}

func (h *hub) PrivateMsgChan() chan Message {
	return h.privateMsgChan
}

func (h *hub) GroupMsgChan() chan Message {
	return h.groupMsgChan
}

func (h *hub) Users() map[uint64]*websocket.Conn {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.users
}
