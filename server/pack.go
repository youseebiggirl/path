package server

import (
	"bytes"
	"encoding/binary"
	"log"
)

// 封包格式：
// ====================== head ====================================================
// +-----------+-----------+-----------+------------------------+-------------------+-----------+
// |  datalen  | msgType   |  dataType |    from[userId]        |     to[userid]    |   data    |
// +-----------+-----------+-----------+------------------------+-------------------+-----------+
//	   uint32      int32       int32          uint64                    uint64          []byte

const MaxPackSize = 25535

type Pack interface {
	HeadSize() uint32
	Packet(msg Message) ([]byte, error)
	UnPack([]byte) (Message, error)
}

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) HeadSize() uint32 {
	// head = datalen + msgType() + dataType + from + to
	//      = 4 + 4 + 4 + 8 + 8 = 28
	return 28
}

func (d *DataPack) Packet(msg Message) ([]byte, error) {
	var (
		dataLen  = msg.Len()
		msgType  = msg.MsgType()
		dataType = msg.DataType()
		from     = msg.From()
		to       = msg.To()
		data     = msg.Data()
	)

	var b bytes.Buffer

	// 按照封包格式依次写入，注意顺序不能出错，否则拆包时会出现不可预计的错误
	if err := binary.Write(&b, binary.BigEndian, &dataLen); err != nil {
		log.Println("packet write [dataLen] error: ", err)
		return nil, err
	}

	if err := binary.Write(&b, binary.BigEndian, &msgType); err != nil {
		log.Println("packet write [msgType] error: ", err)
		return nil, err
	}

	if err := binary.Write(&b, binary.BigEndian, &dataType); err != nil {
		log.Println("packet write [dataType] error: ", err)
		return nil, err
	}

	if err := binary.Write(&b, binary.BigEndian, &from); err != nil {
		log.Println("packet write [from] error: ", err)
		return nil, err
	}

	if err := binary.Write(&b, binary.BigEndian, to); err != nil {
		log.Println("packet write [to] error: ", err)
		return nil, err
	}

	if err := binary.Write(&b, binary.BigEndian, &data); err != nil {
		log.Println("packet write [data] error: ", err)
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *DataPack) UnPack(pkg []byte) (Message, error) {
	r := bytes.NewReader(pkg)

	var m privateChatMessage

	// 拆包同样也按照格式顺序
	if err := binary.Read(r, binary.BigEndian, &m.len); err != nil {
		log.Println("packet read dataLen error: ", err)
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &m.msgType); err != nil {
		log.Println("packet read dataLen error: ", err)
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &m.dataType); err != nil {
		log.Println("packet read dataLen error: ", err)
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &m.from); err != nil {
		log.Println("packet read dataLen error: ", err)
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &m.to); err != nil {
		log.Println("packet read dataLen error: ", err)
		return nil, err
	}

	if err := binary.Read(r, binary.BigEndian, &m.data); err != nil {
		log.Println("packet read dataLen error: ", err)
		return nil, err
	}

	//if m.dataLen > MaxPackSize {
	//	return nil, errors.New("unpack error: too large message")
	//}

	return &m, nil
}
