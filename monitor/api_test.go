package monitor

import (
	"fmt"
	"testing"
)

func TestGetAllMonitors(t *testing.T) {
	monitors, err := GetAllMonitors()
	if err != nil {
		t.Errorf("get all monitors failed: %v", err)
		return
	}

	fmt.Println(monitors)

	for _, m := range monitors {
		n, err := GetNumberOfPhysicalMonitors(m.Handle)
		if err != nil {
			t.Errorf("get number of physical monitors failed: %v", err)
			return
		}
		fmt.Println("n=", n)

		pms, err := GetPhysicalMonitors(m.Handle, n)
		if err != nil {
			t.Errorf("get physical monitors failed: %v", err)
			return
		}

		for _, pm := range pms {
			fmt.Println(pm.Handle, pm.Description)
		}

		info, err := GetMonitorInfo(m.Handle)
		if err != nil {
			t.Errorf("get monitor info failed: %v", err)
			return
		}

		fmt.Println(info)
	}

}

func TestGetDisplayDevice(t *testing.T) {

	idx := 0
	for {
		device, ok := GetDisplayDevice("", idx, EDD_GET_DEVICE_INTERFACE_NAME)
		if !ok {
			break
		}

		fmt.Println(device)

		idx += 1
	}
}

func TestGetAllPhysicalMonitors(t *testing.T) {

	result, err := GetAllPhysicalMonitors()
	if err != nil {
		t.Errorf("test failed: %v", err)
		return
	}

	fmt.Println("result is ===", result)

	for name, handle := range result {
		brightness, err := GetMonitorBrightness(handle)
		if err != nil {
			t.Errorf("get brightness failed: %v", err)
			return
		}

		fmt.Println(name, ":", brightness)
	}
}
