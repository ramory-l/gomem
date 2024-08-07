package memory

import (
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/ramory-l/gomem/pkg/process"
	"golang.org/x/sys/windows"
)

type Gomem struct {
	ProcessId     *uint32
	ProcessHandle *windows.Handle
	ThreadHandle  uintptr
	IsWow64       bool
}

func NewGomem(processId interface{}, desiredAccess *uint32) (*Gomem, error) {
	var err error
	gm := &Gomem{}

	switch processId := any(processId).(type) {
	case string:
		gm.ProcessId, gm.ProcessHandle, err = process.OpenProcessFromName(processId, desiredAccess)
	case uint32:
		gm.ProcessId, gm.ProcessHandle, err = process.OpenProcessFromId(processId, desiredAccess)
	}

	return gm, err
}

func VirtualQuery(handle *windows.Handle, address uintptr) (*windows.MemoryBasicInformation, error) {
	var mbi windows.MemoryBasicInformation
	err := windows.VirtualQueryEx(*handle, address, (*windows.MemoryBasicInformation)(unsafe.Pointer(&mbi)), uintptr(unsafe.Sizeof(mbi)))
	if err != nil {
		return nil, fmt.Errorf("failed to VirtualQueryEx, err: %w", err)
	}
	return &mbi, nil
}

func (gm *Gomem) ReadBytes(address uintptr, size uintptr) ([]byte, error) {
	var bytesRead uintptr
	buffer := make([]byte, size)
	err := windows.ReadProcessMemory(*gm.ProcessHandle, address, &buffer[0], size, &bytesRead)
	if err != nil {
		return nil, fmt.Errorf("failed to read process memory, err: %w", err)
	}
	return buffer[:bytesRead], nil
}

func (gm *Gomem) ReadInt(address uintptr) (*int32, error) {
	bytes, err := gm.ReadBytes(address, 4)
	if err != nil {
		return nil, fmt.Errorf("memory read error at address %x: %w", address, err)
	}
	if len(bytes) < 4 {
		return nil, fmt.Errorf("expected to read 4 bytes, got %d", len(bytes))
	}
	val := int32(binary.LittleEndian.Uint32(bytes))
	return &val, nil
}

func (gm *Gomem) ReadUint(address uintptr) (*uint32, error) {
	bytes, err := gm.ReadBytes(address, 4)
	if err != nil {
		return nil, fmt.Errorf("memory read error at address %x: %w", address, err)
	}
	if len(bytes) < 4 {
		return nil, fmt.Errorf("expected to read 4 bytes, got %d", len(bytes))
	}
	val := binary.LittleEndian.Uint32(bytes)
	return &val, nil
}

func (gm *Gomem) WriteBytes(address uintptr, data []byte, size uintptr) (bool, error) {
	var bytesWritten uintptr
	err := windows.WriteProcessMemory(*gm.ProcessHandle, address, &data[0], size, &bytesWritten)
	if err != nil {
		return false, fmt.Errorf("write process memory failed: %w", err)
	}
	return bytesWritten == size, nil
}

func (gm *Gomem) WriteUint(address uintptr, value uint32) (bool, error) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, value)
	return gm.WriteBytes(address, data, uintptr(len(data)))
}

func (gm *Gomem) ReadLongLong(address uintptr) (*int64, error) {
	bytes, err := gm.ReadBytes(address, 8)
	if err != nil {
		return nil, fmt.Errorf("failed to read int64 bytes: %w", err)
	}
	if len(bytes) < 8 {
		return nil, fmt.Errorf("expected to read 8 bytes, got %d", len(bytes))
	}
	val := int64(binary.LittleEndian.Uint64(bytes))
	return &val, nil
}
