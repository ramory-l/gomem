# Gomem
A Go library to manipulate Windows processes

__Library not tested yet, but you can try__

## Documentation
In progress...

### Example
``` go
func main() {
	gameHandle, err := gomem.NewGomem("cs2.exe", nil)

	if err != nil {
		fmt.Printf("failed to create go mem, err: %v", err)
		return
	}

	fmt.Printf("%+v\n", gameHandle)

	clientDll, err := gomem.ModuleFromName(*gameHandle.ProcessHandle, "client.dll")
	if err != nil {
		fmt.Printf("failed to get module from name, err: %v", err)
		return
	}

	fmt.Printf("%+v\n", clientDll)

	engineDll, err := gomem.ModuleFromName(*gameHandle.ProcessHandle, "engine2.dll")
	if err != nil {
		fmt.Printf("failed to get module from name, err: %v", err)
		return
	}

	fmt.Printf("%+v", engineDll)
}
```