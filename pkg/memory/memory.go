package memory

import (
	"encoding/binary"
	"fmt"
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

func ReadInt(handle windows.Handle, address uintptr) (int32, error) {
	bytes, err := ReadBytes(handle, address, 4)
	if err != nil {
		return 0, fmt.Errorf("memory read error at address %x: %w", address, err)
	}
	if len(bytes) < 4 {
		return 0, fmt.Errorf("expected to read 4 bytes, got %d", len(bytes))
	}
	val := int32(binary.LittleEndian.Uint32(bytes))
	return val, nil
}

func ReadUint(handle windows.Handle, address uintptr) (uint32, error) {
	bytes, err := ReadBytes(handle, address, 4)
	if err != nil {
		return 0, fmt.Errorf("memory read error at address %x: %w", address, err)
	}
	if len(bytes) < 4 {
		return 0, fmt.Errorf("expected to read 4 bytes, got %d", len(bytes))
	}
	val := binary.LittleEndian.Uint32(bytes)
	return val, nil
}
