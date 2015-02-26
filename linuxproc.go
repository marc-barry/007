package main

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
)

func readCPUInfo() (*linuxproc.CPUInfo, error) {
	return linuxproc.ReadCPUInfo("/proc/cpuinfo")
}

func readNetworkDeviceStats() ([]linuxproc.NetworkStat, error) {
	return linuxproc.ReadNetworkStat("/proc/net/dev")
}
