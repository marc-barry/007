package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	HTTPEnableFlag                = "http-enable"
	HTTPPortFlag                  = "http-port"
	CollectInterfaceRateStatsRate = "collect-iface-rate-stat-rate"
)

var (
	Log = logrus.New()

	httpEnable                    = flag.Bool(HTTPEnableFlag, false, "Enable HTTP server.")
	httpPort                      = flag.Int(HTTPPortFlag, 8001, "HTTP server listening port.")
	collectInterfaceRateStatsRate = flag.Int64(CollectInterfaceRateStatsRate, 2, "Rate (in seconds) for which the interface rate stats are collected.")

	stopOnce sync.Once
	stopWg   sync.WaitGroup

	IfaceList = NewInterfaceList()
)

func withLogging(f func()) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("Recovered from panic(%+v)", r)

			Log.WithField("error", err).Panicf("Stopped with panic: %s", err.Error())
		}
	}()

	f()
}

func main() {
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			Log.WithField("signal", sig).Infof("Signalled. Shutting down.")

			stopOnce.Do(func() { shutdown(0) })
		}
	}()

	// Get a list of all interfaces.
	ifaces, err := net.Interfaces()
	if err != nil {
		Log.WithField("error", err).Fatalf("Error getting the list of interfaces.")
	}

	Log.WithField("count", len(ifaces)).Info("Found interfaces.")

	for _, iface := range ifaces {
		Log.WithFields(logrus.Fields{
			"index": iface.Index,
			"name":  iface.Name,
		}).Info("Found interface.")

		IfaceList.Append(iface)
	}

	startCollectors()

	if *httpEnable {
		if err := <-StartHTTPServer(*httpPort); err != nil {
			Log.WithField("error", err).Fatal("Error starting HTTP server.")
		}

		return
	}

	// If HTTP is not enabled we need to block with a wait on a WaitGroup.
	stopWg.Add(1)
	stopWg.Wait()
}

func shutdown(code int) {
	Log.WithField("code", code).Infof("Stopping.")

	// If HTTP is enabled we must exit in order to cause the HTTP server to shutdown.
	if *httpEnable {
		os.Exit(0)
	}

	stopWg.Done()
}

func startCollectors() {
	if err := calculateInterfaceRateStats(); err != nil {
		Log.WithField("error", err).Error("Error calculating interface rate stats.")
	} else {
		go withLogging(func() {
			for {
				select {
				case <-time.Tick(time.Duration(*collectInterfaceRateStatsRate) * time.Second):
					if err := calculateInterfaceRateStats(); err != nil {
						Log.WithField("error", err).Error("Error calculating interface rate stats.")
					}
				}
			}
		})
	}
}
