package main

import (
	"net"
	"strings"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

func getInterfaceIPAddressesString(iface net.Interface) string {
	addrs, err := iface.Addrs()
	if err != nil {
		return err.Error()
	}

	var addrStrings = make([]string, 0, len(addrs))
	for _, addr := range addrs {
		addrStrings = append(addrStrings, addr.String())
	}

	return strings.Join(addrStrings, ", ")
}

func getCPUInfo() *linuxproc.CPUInfo {
	info, err := readCPUInfo()
	if err != nil {
		Log.WithField("error", err).Error("Error reading cpuinfo.")
	}
	return info
}

func getNetworkDeviceStats() []linuxproc.NetworkStat {
	stats, err := readNetworkDeviceStats()
	if err != nil {
		Log.WithField("error", err).Error("Error reading network device stats.")
	}
	return stats
}
