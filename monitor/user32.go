package monitor

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	winuser             = windows.NewLazyDLL("user32.dll")
	enumDisplayMonitors = winuser.NewProc("EnumDisplayMonitors")
	enumDisplayDevices  = winuser.NewProc("EnumDisplayDevicesW")
	getMonitorInfo      = winuser.NewProc("GetMonitorInfoW")
)

type Monitor struct {
	Handle          windows.Handle
	PhysicalHandles []windows.Handle
	Rect            WindowRect
}

type WindowRect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

func GetAllMonitors() ([]Monitor, error) {

	var monitors []Monitor

	callback := windows.NewCallback(func(handle windows.Handle, hdc windows.Handle, rect *WindowRect, param uintptr) uintptr {

		monitor := Monitor{
			Handle: handle,
			Rect:   *rect,
		}

		monitors = append(monitors, monitor)

		return 1
	})

	if ret, _, err := enumDisplayMonitors.Call(0, 0, callback, 0); ret == 0 {
		fmt.Printf("enum display monitors failed: %v", err)
		return nil, err
	}

	return monitors, nil
}

type displayDeviceInternal struct {
	CB           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

type DisplayDevice struct {
	DeviceName   string
	DeviceString string
	StateFlags   uint32
	DeviceID     string
	DeviceKey    string
}

const (
	EDD_GET_DEVICE_INTERFACE_NAME = 0x00000001
)

const (
	DISPLAY_DEVICE_ATTACHED_TO_DESKTOP = 0x00000001
	DISPLAY_DEVICE_MULTI_DRIVER        = 0x00000002
	DISPLAY_DEVICE_PRIMARY_DEVICE      = 0x00000004
	DISPLAY_DEVICE_MIRRORING_DRIVER    = 0x00000008
	DISPLAY_DEVICE_VGA_COMPATIBLE      = 0x00000010
	DISPLAY_DEVICE_REMOVABLE           = 0x00000020
	DISPLAY_DEVICE_ACC_DRIVER          = 0x00000040
	DISPLAY_DEVICE_MODESPRUNED         = 0x08000000
	DISPLAY_DEVICE_RDPUDD              = 0x01000000
	DISPLAY_DEVICE_REMOTE              = 0x04000000
	DISPLAY_DEVICE_DISCONNECT          = 0x02000000
	DISPLAY_DEVICE_TS_COMPATIBLE       = 0x00200000
	DISPLAY_DEVICE_UNSAFE_MODES_ON     = 0x00080000
)

func GetDisplayDevice(name string, idx int, flags uint32) (*DisplayDevice, bool) {
	var dd displayDeviceInternal
	dd.CB = uint32(unsafe.Sizeof(dd))

	namePtr := uintptr(0)
	if name != "" {
		bname := windows.StringToUTF16(name)
		namePtr = uintptr(unsafe.Pointer(&bname[0]))
	}

	ret, _, _ := enumDisplayDevices.Call(namePtr, uintptr(idx), uintptr(unsafe.Pointer(&dd)), uintptr(flags))
	if ret == 0 {
		return nil, false
	}

	result := &DisplayDevice{}
	result.DeviceName = windows.UTF16ToString(dd.DeviceName[:])
	result.DeviceString = windows.UTF16ToString(dd.DeviceString[:])
	result.StateFlags = dd.StateFlags
	result.DeviceID = windows.UTF16ToString(dd.DeviceID[:])
	result.DeviceKey = windows.UTF16ToString(dd.DeviceKey[:])

	return result, true
}

type monitorInfo struct {
	cbSize    uint32
	rcMonitor WindowRect
	rcWork    WindowRect
	dwFlags   uint32
}

type monitorInfoEx struct {
	monitorInfo
	szDevice [32]uint16
}

type MonitorInfo struct {
	MonitorRect WindowRect
	WorkRect    WindowRect
	Flags       uint32
	DeviceName  string
}

func GetMonitorInfo(handle windows.Handle) (*MonitorInfo, error) {
	var info monitorInfoEx
	info.cbSize = uint32(unsafe.Sizeof(info))

	ret, _, err := getMonitorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&info)))
	if ret == 0 {
		fmt.Println("get monitor info failed:", err)
		return nil, err
	}

	result := &MonitorInfo{}
	result.MonitorRect = info.rcMonitor
	result.WorkRect = info.rcWork
	result.Flags = info.dwFlags
	result.DeviceName = windows.UTF16ToString(info.szDevice[:])

	return result, nil
}
