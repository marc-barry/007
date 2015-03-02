package main

import (
	linuxproc "github.com/marc-barry/goprocinfo/linux"
)

const (
	CPUInfoPath     = "/proc/cpuinfo"
	NetworkStatPath = "/proc/net/dev"
)

func readCPUInfo() (*linuxproc.CPUInfo, error) {
	return linuxproc.ReadCPUInfo(CPUInfoPath)
}

func readNetworkDeviceStats() ([]linuxproc.NetworkStat, error) {
	return linuxproc.ReadNetworkStat(NetworkStatPath)
}
