package server

import (
	"fmt"
	"log"
	"testing"
	"time"
	"unsafe"
)

func TestPacket(t *testing.T) {
	msg := NewPrivateChatMessage(
		NewUser(nil, "", ""),
		NewUser(nil, "", ""),
		TextMsg, []byte("123"))
	packet, err := NewDataPack().Packet(msg)
	if err != nil {
		log.Fatalln(err)
	}
	sizeof := unsafe.Sizeof(packet)
	fmt.Println(sizeof)
}

func TestSizeofUser(t *testing.T) {
	u := NewUser(nil, "", "")
	sizeof := unsafe.Sizeof(u)
	fmt.Println(sizeof)
}

type ABC interface {
	a()
}

type abc struct {

}

func (a *abc) a() {

}

func TestSizeofInterface(t *testing.T) {
	aaa := &abc{}
	println(unsafe.Sizeof(aaa))
}

func TestSizeOfTo(t *testing.T) {
	to := To{}.user
	sizeof := unsafe.Sizeof(to)
	fmt.Println(sizeof) // 24
}

func TestSizeOfTimer(t *testing.T) {
	now := time.Now()
	sizeof := unsafe.Sizeof(now)
	fmt.Println(sizeof)	// 24
}
