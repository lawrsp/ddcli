package monitor

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	dxva2                      = windows.NewLazySystemDLL("dxva2.dll")
	getNumberOfPhysicalMonitor = dxva2.NewProc("GetNumberOfPhysicalMonitorsFromHMONITOR")
	getPhysicalMonitors        = dxva2.NewProc("GetPhysicalMonitorsFromHMONITOR")
	getMonitorBrightness       = dxva2.NewProc("GetMonitorBrightness")
	setMonitorBrightness       = dxva2.NewProc("SetMonitorBrightness")
	getMonitorContrast         = dxva2.NewProc("GetMonitorContrast")
	setMonitorContrast         = dxva2.NewProc("SetMonitorContrast")
)

type PhysicalMonitor struct {
	Handle      windows.Handle
	Description string
}

func GetNumberOfPhysicalMonitors(handle windows.Handle) (int, error) {

	var result uint32
	if ret, _, err := getNumberOfPhysicalMonitor.Call(uintptr(handle), uintptr(unsafe.Pointer(&result))); ret == 0 {
		fmt.Printf("get monitor physical numbers failed: %v", err)

		return 0, err
	}

	return int(result), nil
}

type physicalMonitorInternal struct {
	Handle      windows.Handle
	Description [128]uint16
}

func GetPhysicalMonitors(handle windows.Handle, num int) ([]PhysicalMonitor, error) {

	pms := make([]physicalMonitorInternal, num)

	if ret, _, err := getPhysicalMonitors.Call(uintptr(handle), uintptr(num), uintptr(unsafe.Pointer(&pms[0]))); ret == 0 {
		fmt.Printf("get physical monitors failed: %v", err)
		return nil, err
	}

	var result []PhysicalMonitor

	for _, pm := range pms {
		desc := windows.UTF16ToString(pm.Description[:])
		result = append(result, PhysicalMonitor{
			Handle:      pm.Handle,
			Description: desc,
		})
	}

	return result, nil
}

type MonitorBrightness struct {
	Min uint
	Max uint
	Cur uint
}

type MonitorContrast struct {
	Min uint
	Max uint
	Cur uint
}

func GetMonitorBrightness(handle windows.Handle) (*MonitorBrightness, error) {
	var min, max, cur uint

	ret, _, err := getMonitorBrightness.Call(uintptr(handle),
		uintptr(unsafe.Pointer(&min)),
		uintptr(unsafe.Pointer(&cur)),
		uintptr(unsafe.Pointer(&max)))
	if ret == 0 {
		return nil, fmt.Errorf("failed to get monitor brightness %w", err)
	}

	return &MonitorBrightness{
		Max: max,
		Min: min,
		Cur: cur,
	}, nil

}

func GetMonitorContrast(handle windows.Handle) (*MonitorContrast, error) {
	var min, max, cur uint

	ret, _, err := getMonitorContrast.Call(uintptr(handle),
		uintptr(unsafe.Pointer(&min)),
		uintptr(unsafe.Pointer(&cur)),
		uintptr(unsafe.Pointer(&max)))
	if ret == 0 {
		return nil, fmt.Errorf("failed to get monitor contrast %w", err)
	}

	return &MonitorContrast{
		Max: max,
		Min: min,
		Cur: cur,
	}, nil

}

func SetMonitorBrightness(handle windows.Handle, level uint) error {
	brightness, err := GetMonitorBrightness(handle)
	if err != nil {
		return err
	}

	if level > brightness.Max {
		return fmt.Errorf("brightness level exceeds maximum")
	}

	ret, _, err := setMonitorBrightness.Call(uintptr(handle), uintptr(level))
	if ret == 0 {

		return fmt.Errorf("failed to set monitor brightness %w", err)
	}

	return nil
}

func SetMonitorContrast(handle windows.Handle, level uint) error {
	contrast, err := GetMonitorContrast(handle)
	if err != nil {
		return err
	}

	if level > contrast.Max {
		return fmt.Errorf("contrast level exceeds maximum")
	}

	ret, _, err := setMonitorContrast.Call(uintptr(handle), uintptr(level))
	if ret == 0 {
		return fmt.Errorf("failed to set monitor contrast %w", err)
	}

	return nil
}
