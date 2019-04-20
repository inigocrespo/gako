package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type RawKeyEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

var size uintptr = unsafe.Sizeof(RawKeyEvent{})

func main() {

	f, err := os.Open("/dev/input/event0")
	if err != nil {
		panic(err)
		return
	}
	defer f.Close()

	fd := int(f.Fd())

	var lshift bool
	var rshift bool
	var altgr bool

	for {

		b := make([]byte, size)
		_, err := syscall.Read(fd, b)
		if err != nil {
			fmt.Println(err)
			continue
		}

		e := &RawKeyEvent{}
		err = binary.Read(bytes.NewBuffer(b), binary.LittleEndian, e)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if e.Type == 1 && e.Code == 42 && (e.Value == 1 || e.Value == 2) {
			lshift = true
		}

		if e.Type == 1 && e.Code == 42 && e.Value == 0 {
			lshift = false
			continue
		}

		if e.Type == 1 && e.Code == 54 && (e.Value == 1 || e.Value == 2) {
			rshift = true
		}

		if e.Type == 1 && e.Code == 54 && e.Value == 0 {
			rshift = false
			continue
		}

		if e.Type == 1 && e.Code == 100 && (e.Value == 1 || e.Value == 2) {
			altgr = true
		}

		if e.Type == 1 && e.Code == 100 && e.Value == 0 {
			altgr = false
			continue
		}

		if e.Type == 1 && (e.Value == 1 || e.Value == 2) {
			if int(e.Code) < len(keys) {
				if lshift || rshift {
					fmt.Print(shiftKeys[e.Code])
				} else if altgr {
					fmt.Print(altgrtKeys[e.Code])
				} else {
					fmt.Print(keys[e.Code])
				}
			} else {
				fmt.Print("<Unknonw ", e.Code, ">")
			}
		}

		//fmt.Println("Type ", e.Type, "Code ", e.Code, "Value ", e.Value)
	}

}
