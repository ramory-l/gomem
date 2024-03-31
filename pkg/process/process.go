package process

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/ramory-l/gomem/pkg/resources"
	"golang.org/x/sys/windows"
)

var SYNCHRONIZE uint32 = 0x00100000
var STANDARD_RIGHTS_REQUIRED uint32 = 0x000F0000
var PROCESS_ALL_ACCESS uint32 = STANDARD_RIGHTS_REQUIRED | SYNCHRONIZE | 0xFFF

func OpenProcessFromName(processName string, desiredAccess *uint32) (*uint32, *windows.Handle, error) {
	process32, err := processFromName(processName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get process by name, err: %w", err)
	}

	return OpenProcessFromId(process32.ProcessID, desiredAccess)
}

func OpenProcessFromId(processId uint32, desiredAccess *uint32) (*uint32, *windows.Handle, error) {
	desiredAccessTemp := desiredAccess
	if desiredAccess == nil {
		desiredAccessTemp = &PROCESS_ALL_ACCESS
	}
	processHandle, err := windows.OpenProcess(*desiredAccessTemp, false, processId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open process, err: %w", err)
	}
	return &processId, &processHandle, nil
}

func processFromName(processId string) (*windows.ProcessEntry32, error) {
	name := strings.ToLower(processId)
	processes, err := listProcesses()

	if err != nil {
		return nil, fmt.Errorf("failed to get processes list, err: %w", err)
	}

	for _, process := range processes {
		processLower := strings.ToLower(windows.UTF16ToString(process.ExeFile[:]))
		if strings.Contains(name, processLower) {
			return &process, nil
		}
	}

	return nil, fmt.Errorf("failed to find process with name: %s", processId)
}

func listProcesses() ([]windows.ProcessEntry32, error) {
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	processes := make([]windows.ProcessEntry32, 0)

	if err := windows.Process32First(handle, &entry); err != nil {
		return nil, err
	}
	processes = append(processes, entry)

	for windows.Process32Next(handle, &entry) == nil {
		processes = append(processes, entry)
	}

	return processes, nil
}

func ModuleFromName(processHandle *windows.Handle, moduleName string) (*resources.MODULEINFO, error) {
	name := strings.ToLower(moduleName)
	modules, err := enumProcessModule(processHandle)
	if err != nil {
		return nil, fmt.Errorf("failed to enum process module, err: %w", err)
	}

	for _, module := range modules {
		moduleNameTemp, err := module.Name(processHandle)
		if err != nil {
			return nil, fmt.Errorf("failed to get module name, err: %w", err)
		}
		if strings.EqualFold(strings.ToLower(moduleNameTemp), name) {
			return &module, nil
		}
	}

	return nil, fmt.Errorf("failed to find module with name: %s", name)
}

func enumProcessModule(handle *windows.Handle) ([]resources.MODULEINFO, error) {
	var modules [1024]windows.Handle
	var needed uint32

	err := windows.EnumProcessModulesEx(
		*handle,
		&modules[0],
		uint32(unsafe.Sizeof(modules)),
		&needed,
		windows.LIST_MODULES_ALL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to enum process modules, err: %w", err)
	}

	count := int(needed) / int(unsafe.Sizeof(modules[0]))
	var moduleInfos []resources.MODULEINFO

	for _, mod := range modules[:count] {
		var mi resources.MODULEINFO
		err = windows.GetModuleInformation(
			*handle,
			mod,
			&mi.ModuleInfo,
			uint32(unsafe.Sizeof(mi)))
		if err != nil {
			return nil, fmt.Errorf("failed to get module information, err: %w", err)
		}
		moduleInfos = append(moduleInfos, mi)
	}

	return moduleInfos, nil
}
