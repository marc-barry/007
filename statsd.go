package main

import (
	"encoding/json"
	"reflect"
	"strings"
)

func collectNetworkDeviceStats() {
	stats, err := readNetworkDeviceStats()

	metricPrefix := "net.dev."

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

				if metricName := field.Tag.Get("json"); metricName != "" {
					if err := StatsdClient.Count(strings.Join([]string{metricPrefix, metricName}, ""), int64(value), []string{"iface:" + stat.Iface}, 1); err != nil {
						Log.WithField("error", err).Error("Couldn't submit event to statsd.")
					}
				}
			}
		} else {
			Log.Infof("Filtered interface %s from network device stat collection.", stat.Iface)
		}
	}
}
