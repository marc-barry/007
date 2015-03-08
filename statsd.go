package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	networkDeviceStatsMetricPrefix = "net.dev."
	netstatStatsMetricPrefix       = "net.netstat."
)

type LoggedStat struct {
	Name string      `json:"name"`
	Time time.Time   `json:"time"`
	Stat interface{} `json:"stat"`
}

func collectNetworkDeviceStats() {
	stats, err := readNetworkDeviceStats()

	if err != nil {
		Log.WithField("error", err).Error("Error reading network device stats. Can't collect network device stats.")
	}

	for _, stat := range stats {
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

func logNetworkDeviceStats() {
	stats, err := readNetworkDeviceStats()

	if err != nil {
		Log.WithField("error", err).Error("Error reading network device stats. Can't log network device stats.")
	}

	loggedStat := LoggedStat{NetworkStatPath, time.Now(), stats}

	jsonBytes, err := json.Marshal(loggedStat)
	if err != nil {
		Log.WithField("error", err).Errorf("Error JSON marshalling: %+v.", stats)
	}

	fmt.Println(string(jsonBytes))
}

func collectNetstatStats() {
	stats, err := readNetstatStats()

	if err != nil {
		Log.WithField("error", err).Error("Error reading network device stats. Can't collect netstat stats.")
	}

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

func logNetstatStats() {
	stats, err := readNetstatStats()

	if err != nil {
		Log.WithField("error", err).Error("Error reading network device stats. Can't log netstat stats.")
	}

	loggedStat := LoggedStat{NetstatStatPath, time.Now(), stats}

	jsonBytes, err := json.Marshal(loggedStat)
	if err != nil {
		Log.WithField("error", err).Errorf("Error JSON marshalling: %+v.", stats)
	}

	fmt.Println(string(jsonBytes))
}
