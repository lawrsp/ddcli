package monitor

import (
	"fmt"

	"golang.org/x/sys/windows"
)

type MonitorWithInfo struct {
	Monitor
	Info *MonitorInfo
}

func GetAllPhysicalMonitors() (map[string]windows.Handle, error) {

	var monitors []MonitorWithInfo

	mts, err := GetAllMonitors()
	if err != nil {
		return nil, fmt.Errorf("get all monitors failed: %w", err)
	}

	for _, m := range mts {
		n, err := GetNumberOfPhysicalMonitors(m.Handle)
		if err != nil {
			return nil, fmt.Errorf("get number of physical monitors failed: %w", err)
		}

		pms, err := GetPhysicalMonitors(m.Handle, n)
		if err != nil {
			return nil, fmt.Errorf("get physical monitors failed: %w", err)
		}
		for _, pm := range pms {
			m.PhysicalHandles = append(m.PhysicalHandles, pm.Handle)
		}

		info, err := GetMonitorInfo(m.Handle)
		if err != nil {
			return nil, fmt.Errorf("get monitor info failed: %w", err)
		}

		monitors = append(monitors, MonitorWithInfo{
			m,
			info,
		})
	}

	result := map[string]windows.Handle{}

	// Loop through adapters
	var adapterDevIndex int
	for {
		device, ok := GetDisplayDevice("", adapterDevIndex, 0)
		if !ok {
			break
		}
		adapterDevIndex += 1

		var displayDevIndex int
		for {

			displayDev, ok := GetDisplayDevice(device.DeviceName, displayDevIndex, EDD_GET_DEVICE_INTERFACE_NAME)
			if !ok {
				break
			}
			displayDevIndex += 1

			if (displayDev.StateFlags&DISPLAY_DEVICE_ATTACHED_TO_DESKTOP) == 0 ||
				displayDev.StateFlags&DISPLAY_DEVICE_MIRRORING_DRIVER != 0 {
				continue
			}

			for _, m := range monitors {

				for i, handle := range m.PhysicalHandles {

					monitorName := m.Info.DeviceName + "\\Monitor" + fmt.Sprintf("%d", i)

					if monitorName == displayDev.DeviceName {
						result[monitorName] = handle
					}
				}
			}
		}
	}

	return result, nil
}
