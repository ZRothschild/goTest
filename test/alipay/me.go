package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	var num uint64
	num = 0x11
	fmt.Printf("num = %x\n", num)

	enc := make([]byte, 1)
	//0000000000001234
	// 转化为大端
	binary.BigEndian.PutUint64(enc, num)
	fmt.Printf("bigendian enc = %x\n", enc)

	// 转化为小端
	binary.LittleEndian.PutUint64(enc, num)
	fmt.Printf("littleendian enc = %x\n", enc)
}
