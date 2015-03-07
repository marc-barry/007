package main

import (
	"encoding/json"
	"reflect"
	"strings"
)

const (
	networkDeviceStatsMetricPrefix = "net.dev."
	netstatStatsMetricPrefix       = "net.netstat."
)

func collectNetworkDeviceStats() {
	stats, err := readNetworkDeviceStats()

	if err != nil {
		Log.WithField("error", err).Error("Error reading network device stats. Can't collect network device stats.")
	}

	for _, stat := range stats {
		jsonBytes, err := json.Marshal(stat)
		if err != nil {
			Log.WithField("error", err).Errorf("Error JSON marshalling: %+v.", stat)
		}

		Log.WithField("stat", NetworkStatPath).Infof(string(jsonBytes))

		if ifaceRegExp == nil || !ifaceRegExp.MatchString(stat.Iface) {
			elem := reflect.ValueOf(&stat).Elem()
			typeOfElem := elem.Type()

			for i := 0; i < elem.NumField(); i++ {
				field := typeOfElem.Field(i)

				if field.Name == "Iface" {
					continue
				}

				value := elem.Field(i).Uint()

				_, collect := StatsMap[field.Name]

				if metricName := field.Tag.Get("json"); collect && metricName != "" {
					if err := StatsdClient.Count(strings.Join([]string{networkDeviceStatsMetricPrefix, metricName}, ""), int64(value), []string{"iface:" + stat.Iface}, 1); err != nil {
						Log.WithField("error", err).Error("Couldn't submit event to statsd.")
					}
				}
			}
		} else {
			Log.Infof("Filtered interface %s from network device stat collection.", stat.Iface)
		}
	}
}

func collectNetstatStats() {
	stats, err := readNetstatStats()

	if err != nil {
		Log.WithField("error", err).Error("Error reading network device stats. Can't collect netstat stats.")
	}

	jsonBytes, err := json.Marshal(stats)
	if err != nil {
		Log.WithField("error", err).Errorf("Error JSON marshalling: %+v.", stats)
	}

	Log.WithField("stat", NetstatStatPath).Infof(string(jsonBytes))

	elem := reflect.ValueOf(stats).Elem()
	typeOfElem := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := typeOfElem.Field(i)

		value := elem.Field(i).Uint()

		_, collect := StatsMap[field.Name]

		if metricName := field.Tag.Get("json"); collect && metricName != "" {
			if err := StatsdClient.Count(strings.Join([]string{netstatStatsMetricPrefix, metricName}, ""), int64(value), []string{}, 1); err != nil {
				Log.WithField("error", err).Error("Couldn't submit event to statsd.")
			}
		}
	}
}
