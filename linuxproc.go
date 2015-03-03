package main

import (
	"reflect"

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

func getNetworkDeviceStatsList() []string {
	stat := linuxproc.NetworkStat{}

	elem := reflect.ValueOf(&stat).Elem()
	typeOfElem := elem.Type()

	list := make([]string, 0)

	for i := 0; i < elem.NumField(); i++ {
		if field := typeOfElem.Field(i); field.Name != "Iface" {
			list = append(list, field.Name)
		}
	}

	return list
}
