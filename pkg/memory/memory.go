package memory

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func VirtualQuery(handle windows.Handle, address uintptr) (*windows.MemoryBasicInformation, error) {
	var mbi windows.MemoryBasicInformation
	err := windows.VirtualQueryEx(handle, address, (*windows.MemoryBasicInformation)(unsafe.Pointer(&mbi)), uintptr(unsafe.Sizeof(mbi)))
	if err != nil {
		return nil, err
	}
	return &mbi, nil
}

func ReadBytes(handle windows.Handle, address uintptr, size uintptr) ([]byte, error) {
	var bytesRead uintptr
	buffer := make([]byte, size)
	err := windows.ReadProcessMemory(handle, address, &buffer[0], size, &bytesRead)
	if err != nil {
		return nil, err
	}
	return buffer[:bytesRead], nil
}
