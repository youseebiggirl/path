package server

import (
	"time"
)

type DataType int32

const (
	PhotoMsg DataType = iota
	VideoMsg
	TextMsg
)

type MsgType int32

const (
	PrivateMsg MsgType = iota // 私聊
	GroupMsg                  // 群聊
)

// To 如果是 private message，则设置 To 中的 user，
// 如果是 group message 则设置 To 中的 groupId
type To struct {
	user    User
	groupId uint64
}

type Message interface {
	From() User
	To() *To
	MsgType() MsgType
	DataType() DataType
	Data() []byte
	Len() uint64
	Time() time.Time
}

// -------------------- 私聊消息 --------------------
type privateChatMessage struct {
	from     User     `json:"from"`
	to       *To      `json:"to"`
	msgType  MsgType  `json:"msg_type"`
	dataType DataType `json:"data_type"`
	data     []byte   `json:"data"`
	time     time.Time
}

func NewPrivateChatMessage(from, to User, dataType DataType, data []byte) *privateChatMessage {
	return &privateChatMessage{
		from:     from,
		to:       &To{user: to},
		dataType: dataType,
		msgType:  PrivateMsg,
		data:     data,
		time:     time.Now(),
	}
}

func (m *privateChatMessage) From() User {
	return m.from
}

func (m *privateChatMessage) To() *To {
	return m.to
}

func (m *privateChatMessage) DataType() DataType {
	return m.dataType
}

func (m *privateChatMessage) MsgType() MsgType {
	return m.msgType
}

func (m *privateChatMessage) Data() []byte {
	return m.data
}

func (m *privateChatMessage) Len() uint64 {
	return (uint64)(len(m.data))
}

func (m *privateChatMessage) Time() time.Time {
	return m.time
}

// -------------------- 群聊消息 --------------------

type groupChatMessage struct {
	from     User
	to       *To
	data     []byte
	time     time.Time
	msgType  MsgType
	dataType DataType
}

func NewGroupChatMessage(from User, groupId uint64, dataType DataType, data []byte) *groupChatMessage {
	return &groupChatMessage{
		from:     from,
		to:       &To{groupId: groupId},
		dataType: dataType,
		msgType:  GroupMsg,
		data:     data,
		time:     time.Now(),
	}
}

func (m *groupChatMessage) From() User {
	return m.from
}

func (m *groupChatMessage) To() *To {
	return m.to
}

func (m *groupChatMessage) MsgType() MsgType {
	return m.msgType
}

func (m *groupChatMessage) DataType() DataType {
	return m.dataType
}

func (m *groupChatMessage) Data() []byte {
	return m.data
}

func (m *groupChatMessage) Len() uint64 {
	return (uint64)(len(m.data))
}

func (m *groupChatMessage) Time() time.Time {
	return m.time
}
