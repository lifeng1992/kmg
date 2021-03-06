package kmgSys

import (
	"errors"
	"io"
)

type DeviceType string

func (s DeviceType) String() string {
	return string(s)
}

var DeviceTypeTap DeviceType = "tap"
var DeviceTypeTun DeviceType = "tun"

var ErrAllDeviceBusy = errors.New("tun/tap: all dev is busy.")

// Interface is a TUN/TAP interface.
type TunTapInterface interface {
	io.ReadWriteCloser
	GetDeviceType() DeviceType
	Name() string
}

/*
//set tun p2p ip and up this device
// mtu default to 1500
func SetP2PIpAndUp(ifac Interface, srcIp string, destIp string, mtu int) error {
	if mtu == 0 {
		mtu = 1500
	}
	switch runtime.GOOS {
	case "darwin":
		return kmgCmd.StdioSliceRun([]string{"ifconfig", ifac.Name(), srcIp, destIp, "mtu", strconv.Itoa(mtu), "up"})
	case "linux":
		return kmgCmd.StdioSliceRun([]string{"ifconfig", ifac.Name(), srcIp, "pointopoint", destIp, "mtu", strconv.Itoa(mtu), "up"})
	default:
		return ErrPlatformNotSupport
	}
}

//set mtu on a device
func SetMtu(ifac Interface, mtu int) error {
	return kmgCmd.StdioSliceRun([]string{"ifconfig", ifac.Name(), "mtu", strconv.Itoa(mtu)})
}
*/
