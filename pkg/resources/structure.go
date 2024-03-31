package resources

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// : Enables execute access to the committed region of pages.
// An attempt to write to the committed region results in an access violation.

var MEMORY_PROTECTION_PAGE_EXECUTE uint32 = 0x10

// : Enables execute or read-only access to the committed region of pages.
// An attempt to write to the committed region results in an access violation.
var MEMORY_PROTECTION_PAGE_EXECUTE_READ uint32 = 0x20

// : Enables execute, read-only, or read/write access to the committed region of pages.
var MEMORY_PROTECTION_PAGE_EXECUTE_READWRITE uint32 = 0x40

// : Enables execute, read-only, or copy-on-write access to a mapped view of a file mapping object.
// An attempt to write to a committed copy-on-write page results in a private
// copy of the page being made for the process.
// The private page is marked as PAGE_EXECUTE_READWRITE, and the change is written to the new page.
var MEMORY_PROTECTION_PAGE_EXECUTE_WRITECOPY uint32 = 0x80

// : Disables all access to the committed region of pages.
// An attempt to read from, write to, or execute the committed region results in an access violation.
var MEMORY_PROTECTION_PAGE_NOACCESS uint32 = 0x01

// : Enables read-only access to the committed region of pages.
// An attempt to write to the committed region results in an access violation.
// If Data Execution Prevention is enabled,
// an attempt to execute code in the committed region results in an access violation.
var MEMORY_PROTECTION_PAGE_READONLY uint32 = 0x02

// : Enables read-only or read/write access to the committed region of pages.
// If Data Execution Prevention is enabled, attempting to execute code in the
// committed region results in an access violation.
var MEMORY_PROTECTION_PAGE_READWRITE uint32 = 0x04

// : Enables read-only or copy-on-write access to a mapped view of a file mapping object.
// An attempt to write to a committed copy-on-write page results in a private copy of
// the page being made for the process. The private page is marked as PAGE_READWRITE,
// and the change is written to the new page. If Data Execution Prevention is enabled,
// attempting to execute code in the committed region results in an access violation.
var MEMORY_PROTECTION_PAGE_WRITECOPY uint32 = 0x08

// : Pages in the region become guard pages.
// Any attempt to access a guard page causes the system to
// raise a STATUS_GUARD_PAGE_VIOLATION exception and turn off the guard page status.
// Guard pages thus act as a one-time access alarm. For more information, see Creating Guard Pages.
var MEMORY_PROTECTION_PAGE_GUARD uint32 = 0x100

// : Sets all pages to be non-cachable.
// Applications should not use this attribute except when explicitly required for a device.
// Using the interlocked functions with memory that is mapped with
// SEC_NOCACHE can result in an EXCEPTION_ILLEGAL_INSTRUCTION exception.
var MEMORY_PROTECTION_PAGE_NOCACHE uint32 = 0x200

// : Sets all pages to be write-combined.
// : Applications should not use this attribute except when explicitly required for a device.
// Using the interlocked functions with memory that is mapped as write-combined can result in an
// EXCEPTION_ILLEGAL_INSTRUCTION exception.
var MEMORY_PROTECTION_PAGE_WRITECOMBINE uint32 = 0x400

var (
	//: Allocates memory charges (from the overall size of memory and the paging files on disk)
	// for the specified reserved memory pages. The function also guarantees that when the caller later
	// initially accesses the memory, the contents will be zero.
	// Actual physical pages are not allocated unless/until the virtual addresses are actually accessed.
	MEMORY_STATE_MEM_COMMIT = 0x1000
	//: XXX
	MEMORY_STATE_MEM_FREE = 0x10000
	//: XXX
	MEMORY_STATE_MEM_RESERVE = 0x2000
	//: Decommits the specified region of committed pages. After the operation, the pages are in the reserved state.
	//: https://msdn.microsoft.com/en-us/library/windows/desktop/aa366894(v=vs.85).aspx
	MEMORY_STATE_MEM_DECOMMIT = 0x4000
	//: Releases the specified region of pages. After the operation, the pages are in the free state.
	//: https://msdn.microsoft.com/en-us/library/windows/desktop/aa366894(v=vs.85).aspx
	MEMORY_STATE_MEM_RELEASE = 0x8000
)

type MODULEINFO struct {
	windows.ModuleInfo
}

func (mi *MODULEINFO) Name(processHandle *windows.Handle) (string, error) {
	var moduleName [windows.MAX_PATH]uint16
	moduleNamePtr := &moduleName[0]
	size := uint32(len(moduleName) * int(unsafe.Sizeof(moduleName[0])))
	err := windows.GetModuleBaseName(*processHandle, windows.Handle(mi.BaseOfDll), moduleNamePtr, size)
	if err != nil {
		return "", fmt.Errorf("failed to get module base name, err: %w", err)
	}
	return windows.UTF16ToString(moduleName[:]), nil
}
