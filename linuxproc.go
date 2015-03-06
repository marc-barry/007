package main

import (
	"reflect"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

const (
	CPUInfoPath     = "/proc/cpuinfo"
	NetworkStatPath = "/proc/net/dev"
	NetstatStatPath = "/proc/net/netstat"
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

func readNetstatStats() (*linuxproc.Netstat, error) {
	return linuxproc.ReadNetstat(NetstatStatPath)
}

func getNetstatStatsList() []string {
	stat := linuxproc.Netstat{}

	elem := reflect.ValueOf(&stat).Elem()
	typeOfElem := elem.Type()

	list := make([]string, 0)

	for i := 0; i < elem.NumField(); i++ {
		list = append(list, typeOfElem.Field(i).Name)
	}

	return list
}
