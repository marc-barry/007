package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
)

const (
	PortFlag = "p"
)

var (
	Log       = logrus.New()
	stopOnce  sync.Once
	IfaceList = NewInterfaceList()

	port = flag.Int(PortFlag, 8001, "HTTP server listening port.")
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

	// Read stats to get initial data for network interface stats.
	readNetworkDeviceStats()

	if err = <-StartHTTPServer(*port); err != nil {
		Log.WithField("error", err).Fatal("Error starting HTTP server.")
	}
}

func shutdown(code int) {
	Log.WithField("code", code).Infof("Stopping.")

	os.Exit(0)
}
