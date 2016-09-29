// test
package main

import (
	"fmt"

	"go/cmkj_server_go/util"
)

func main() {
	s := []byte{0x00, 0x00, 0x03, 0xe8}
	b := util.NewByteBufferWith(s)
	fmt.Println(b.ReadInt())

	//b := util.NewByteBufer()

	b.WriteByte(11)
	b.WriteByte(2)
	b.WriteInt(10)
	b.WriteUTF("c额测试")
	b.WriteByte(9)
	b.WriteDouble(0.09)
	b.WriteShort(1000)

	/*
		g := b.ReadByte()
		t := b.ReadByte()
		fmt.Println(10*g + t)
		fmt.Println(b.ReadInt())
		fmt.Println(b.ReadUTF())
	*/
	b.Buf.Next(2)
	fmt.Println(b.Bytes())
	fmt.Println(b.ReadInt())
	fmt.Println(b.ReadUTF())
	fmt.Println(b.ReadByte())
	fmt.Println(b.ReadDouble())
	fmt.Println(b.ReadShort())
	fmt.Println(b.Bytes())
}
