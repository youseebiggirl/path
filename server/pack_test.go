package server

import (
	"fmt"
	"log"
	"testing"
	"time"
	"unsafe"
)

func TestSizeOfPacket(t *testing.T) {
	msg := NewPrivateChatMessage(123, 456, TextMsg, []byte("123"))
	packet, err := NewDataPack().Packet(msg)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println(packet)
	sizeof := unsafe.Sizeof(packet)
	fmt.Println(sizeof) // 24
}

func TestUnpack(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	msg := NewPrivateChatMessage(123, 456, TextMsg, []byte("123"))
	pack := NewDataPack()

	packet, err := pack.Packet(msg)
	if err != nil {
		log.Fatalln(err)
	}

	unPack, err := pack.UnPack(packet)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v \n", unPack)
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
	fmt.Println(sizeof) // 24
}
