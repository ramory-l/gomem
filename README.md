# Gomem
A Go library to manipulate Windows processes

## Documentation
In progress...

### Example
``` go
type Data struct {
	DwEntityList      uintptr
	DwLocalPlayerPawn uintptr
	M_iIDEntIndex     uintptr
	M_iTeamNum        uintptr
	M_iHealth         uintptr
	DwForceAttack     uintptr
}

func main() {
	gm, err := memory.NewGomem("cs2.exe", nil)
	if err != nil {
		fmt.Printf("failed to create go mem, err: %v", err)
		return
	}

	clientDll, err := process.ModuleFromName(gm.ProcessHandle, "client.dll")
	if err != nil {
		fmt.Printf("failed to get module from name, err: %v", err)
		return
	}

	data := &Data{
		DwEntityList:      0x18C2D58,
		DwLocalPlayerPawn: 0x17371A8,
		DwForceAttack:     0x1730020,
		M_iIDEntIndex:     0x15A4,
		M_iTeamNum:        0x3CB,
		M_iHealth:         0x334,
	}

	triggerBot(gm, clientDll.BaseOfDll, data)

}

func triggerBot(gm *memory.Gomem, clientDllBase uintptr, data *Data) {
	for {
		player, err := gm.ReadLongLong(clientDllBase + data.DwLocalPlayerPawn)
		if err != nil {
			fmt.Printf("failed to read dw local player pawn, err: %v", err)
			time.Sleep(time.Millisecond)
			continue
		}
		entityId, err := gm.ReadInt(uintptr(*player) + data.M_iIDEntIndex)
		if err != nil {
			fmt.Printf("failed to read M_iIDEntIndex, err: %v", err)
			time.Sleep(time.Millisecond)
			continue
		}

		if *entityId > 0 {
			entList, err := gm.ReadLongLong(clientDllBase + data.DwEntityList)
			if err != nil {
				fmt.Printf("failed to read dw entity list, err: %v", err)
				time.Sleep(time.Millisecond)
				continue
			}

			entEntry, err := gm.ReadLongLong(uintptr(*entList) + 0x8*(uintptr(*entityId)>>9) + 0x10)
			if err != nil {
				fmt.Printf("failed to read ent entry, err: %v", err)
				time.Sleep(time.Millisecond)
				continue
			}

			entity, err := gm.ReadLongLong(uintptr(*entEntry) + 120*(uintptr(*entityId)&0x1FF))
			if err != nil {
				fmt.Printf("failed to read dw entity list, err: %v", err)
				time.Sleep(time.Millisecond)
				continue
			}

			entityTeam, err := gm.ReadInt(uintptr(*entity) + data.M_iTeamNum)
			if err != nil {
				fmt.Printf("failed to read entity team, err: %v", err)
				time.Sleep(time.Millisecond)
				continue
			}

			playerTeam, err := gm.ReadInt(uintptr(*player) + data.M_iTeamNum)
			if err != nil {
				fmt.Printf("failed to read entity team, err: %v", err)
				time.Sleep(time.Millisecond)
				continue
			}

			if entityTeam != playerTeam {
				entityHp, err := gm.ReadInt(uintptr(*entity) + data.M_iHealth)
				if err != nil {
					fmt.Printf("failed to read entity hp, err: %v", err)
					time.Sleep(time.Millisecond)
					continue
				}
				if *entityHp > 0 {
					// time.Sleep(...)
					_, err := gm.WriteUint(clientDllBase+data.DwForceAttack, 65537)
					if err != nil {
						fmt.Printf("failed to force attack, err: %v", err)
						time.Sleep(time.Millisecond)
						continue
					}
					_, err = gm.WriteUint(clientDllBase+data.DwForceAttack, 256)
					if err != nil {
						fmt.Printf("failed to force attack, err: %v", err)
						time.Sleep(time.Millisecond)
						continue
					}
				}
			}
		}
	}
}

```