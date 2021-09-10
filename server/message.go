package server

import (
	"time"
)

type DataType int32

const (
	PhotoMsg DataType = iota + 1
	VideoMsg
	TextMsg
)

type MsgType int32

const (
	PrivateMsg MsgType = iota + 1 // 私聊
	GroupMsg                      // 群聊
)

// To 如果是 private message，则设置 To 中的 user，
// 如果是 group message 则设置 To 中的 groupId
type To struct {
	user    User
	groupId uint64
}

type Message interface {
	From() uint64       // 发送方（使用 id 作为标识，便于传输）
	To() uint64         // 接收方（使用 id 作为标识，如果是群发消息，则 id 为群号，如果是私聊信息，则 id 为接受者的 id）
	MsgType() MsgType   // 消息类型（私聊、群聊）
	DataType() DataType // 数据类型（视频、文本、图片）
	Data() []byte
	Len() uint32
	Time() time.Time
}

// -------------------- 私聊消息 --------------------
type privateChatMessage struct {
	basicMessage
}

func NewPrivateChatMessage(from, to uint64, dataType DataType, data []byte) *privateChatMessage {
	m := &privateChatMessage{basicMessage{
		from:     from,
		to:       to,
		dataType: dataType,
		msgType:  PrivateMsg,
		data:     data,
		time:     time.Now(),
		len:      (uint32)(len(data)),
	}}
	return m
}

// -------------------- 群聊消息 --------------------

type groupChatMessage struct {
	basicMessage
}

func NewGroupChatMessage(from, groupId uint64, dataType DataType, data []byte) *groupChatMessage {
	return &groupChatMessage{basicMessage{
		from:     from,
		to:       groupId,
		dataType: dataType,
		msgType:  GroupMsg,
		data:     data,
		time:     time.Now(),
		len:      (uint32)(len(data)),
	}}
}

// 通用类
type basicMessage struct {
	from     uint64
	to       uint64
	msgType  MsgType
	dataType DataType
	data     []byte
	time     time.Time
	len      uint32
}

func (b *basicMessage) From() uint64 {
	return b.from
}

func (b *basicMessage) To() uint64 {
	return b.to
}

func (b *basicMessage) MsgType() MsgType {
	return b.msgType
}

func (b *basicMessage) DataType() DataType {
	return b.dataType
}

func (b *basicMessage) Data() []byte {
	return b.data
}

func (b *basicMessage) Len() uint32 {
	return b.len
}

func (b *basicMessage) Time() time.Time {
	return b.time
}
