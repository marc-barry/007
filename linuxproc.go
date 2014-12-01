package main

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
)

func readCPUInfo() *linuxproc.CPUInfo {
	cpuInfo, err := linuxproc.ReadCPUInfo("/proc/cpuinfo")

	if err != nil {
		Log.WithField("error", err).Error("Error reading cpuinfo.")
		return nil
	}

	return cpuInfo
}
