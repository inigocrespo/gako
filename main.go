package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type InputEvent struct {
	Time  syscall.Timeval // 16 bytes
	Type  uint16          // 2 bytes
	Code  uint16          // 2 bytes
	Value int32           // 4 bytes
}

var size uintptr = unsafe.Sizeof(InputEvent{})

func main() {

	kbd, err := os.Open("/dev/input/event0")
	if err != nil {
		panic(err)
		return
	}
	defer kbd.Close()

	for {
		b := make([]byte, size)
		n, err := kbd.Read(b)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ie := &InputEvent{}
		err = binary.Read(bytes.NewBuffer(b), binary.LittleEndian, ie)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(ie)
	}

}
