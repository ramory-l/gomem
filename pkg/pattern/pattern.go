package pattern

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/ramory-l/gomem/pkg/memory"
	"github.com/ramory-l/gomem/pkg/resources"
	"golang.org/x/sys/windows"
)

func PatternScanModule(handle windows.Handle, module resources.MODULEINFO, pattern string, returnMultiple bool) ([]uintptr, error) {
	var results []uintptr
	var err error
	baseAddress := module.BaseOfDll
	maxAddress := module.BaseOfDll + uintptr(module.SizeOfImage)
	pageAddress := baseAddress
	var found []uintptr

	if !returnMultiple {
		for pageAddress < maxAddress {
			pageAddress, found, err = scanPatternPage(handle, pageAddress, pattern, returnMultiple)
			if err != nil {
				return nil, fmt.Errorf("failed to scan page, err: %w", err)
			}

			if found != nil {
				results = append(results, found...)
				if !returnMultiple {
					break
				}
			}
		}
	}

	return results, nil
}

func scanPatternPage(handle windows.Handle, address uintptr, pattern string, returnMultiple bool) (uintptr, []uintptr, error) {
	mbi, err := memory.VirtualQuery(handle, address)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to query virtual, err: %w", err)
	}

	nextRegion := mbi.BaseAddress + mbi.RegionSize
	allowedProtections := []uint32{
		resources.MEMORY_PROTECTION_PAGE_EXECUTE_READ,
		resources.MEMORY_PROTECTION_PAGE_EXECUTE_READWRITE,
		resources.MEMORY_PROTECTION_PAGE_READWRITE,
		resources.MEMORY_PROTECTION_PAGE_READONLY,
	}

	if mbi.State != uint32(resources.MEMORY_STATE_MEM_COMMIT) || slices.Contains(allowedProtections, mbi.Protect) {
		return nextRegion, nil, nil
	}

	pageBytes, err := memory.ReadBytes(handle, address, mbi.RegionSize)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read bytes, err: %w", err)
	}

	var foundAddresses []uintptr

	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to compile regexp, err: %w", err)
	}

	if !returnMultiple {
		match := re.FindIndex(pageBytes)

		if match != nil {
			foundAddresses = append(foundAddresses, address+uintptr(match[0]))
		}
	} else {
		matches := re.FindAllIndex(pageBytes, -1)
		for _, match := range matches {
			foundAddresses = append(foundAddresses, address+uintptr(match[0]))
		}
	}

	return nextRegion, foundAddresses, nil
}
