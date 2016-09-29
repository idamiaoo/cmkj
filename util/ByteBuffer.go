// ByteBufer
package util

//package main

import (
	"bytes"
	"encoding/binary"
	//"fmt"
	//"strings"
)

type ByteBuffer struct {
	Buf *bytes.Buffer
	//b_buf := bytes.NewBuffer([]byte{})

}

func NewByteBuffer() *ByteBuffer {
	b := new(ByteBuffer)
	b.Buf = bytes.NewBuffer([]byte{})
	return b

}

func NewByteBufferWith(v []byte) *ByteBuffer {
	b := new(ByteBuffer)
	b.Buf = bytes.NewBuffer(v)
	return b

}

func (b *ByteBuffer) Bytes() []byte {
	return b.Buf.Bytes()
}

func (b *ByteBuffer) Next(n int) {
	b.Buf.Next(n)
}

func (b *ByteBuffer) ReadByte() byte {
	var v byte
	binary.Read(b.Buf, binary.BigEndian, &v)
	return v
}

func (b *ByteBuffer) ReadShort() int16 {
	var v int16
	binary.Read(b.Buf, binary.BigEndian, &v)
	return v
}

func (b *ByteBuffer) ReadInt() int32 {
	var v int32
	binary.Read(b.Buf, binary.BigEndian, &v)
	return v
}

func (b *ByteBuffer) ReadDouble() float64 {
	var v float64
	binary.Read(b.Buf, binary.BigEndian, &v)
	return v
}

func (b *ByteBuffer) ReadUTF() string {
	var v uint16
	binary.Read(b.Buf, binary.BigEndian, &v)
	if v > 0 {
		buf := make([]byte, v)
		b.Buf.Read(buf)
		//binary.Read(b.b_buf, binary.BigEndian, &buf)
		str := string(buf)
		return str
	}
	return ""
}

func (b *ByteBuffer) WriteByte(v byte) {
	binary.Write(b.Buf, binary.BigEndian, v)

}

func (b *ByteBuffer) WriteShort(v int16) {
	binary.Write(b.Buf, binary.BigEndian, v)

}

func (b *ByteBuffer) WriteInt(v int32) {
	binary.Write(b.Buf, binary.BigEndian, v)

}

func (b *ByteBuffer) WriteDouble(v float64) {
	binary.Write(b.Buf, binary.BigEndian, v)

}

func (b *ByteBuffer) WriteUTF(v string) {
	l := uint16(len(v))

	binary.Write(b.Buf, binary.BigEndian, l)
	//binary.Write(b.b_buf, binary.BigEndian, v)
	b.Buf.WriteString(v)

}

//func main() {
//	s := []byte{0x00, 0x00, 0x03, 0xe8}
//	b := NewByteBuferWith(s)
//	fmt.Println(b.ReadInt())

//	//b := NewByteBufer()

//	b.WriteInt(10)
//	b.WriteUTF("c额测试")
//	b.WriteByte(9)
//	b.WriteDouble(0.09)

//	fmt.Println(b.ReadInt())
//	fmt.Println(b.ReadUTF())

//	fmt.Println(b.ReadByte())
//	fmt.Println(b.ReadDouble())
//}
