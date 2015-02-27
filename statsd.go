package main

func collectNetworkDeviceStats() {

	stats, err := readNetworkDeviceStats()

	if err != nil {
		Log.WithField("error", err).Error("Error reading network device stats. Can't collect network device stats.")
	}

	for _, stat := range stats {
		// Rx stats
		if err := StatsdClient.Count("net.dev.rxbytes", int64(stat.RxBytes), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.rxpackets", int64(stat.RxPackets), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.rxerrs", int64(stat.RxErrs), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.rxdrop", int64(stat.RxDrop), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.rxfifo", int64(stat.RxFifo), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.rxframe", int64(stat.RxFrame), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.rxcompressed", int64(stat.RxCompressed), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.rxmulticast", int64(stat.RxMulticast), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}

		// Tx stats
		if err := StatsdClient.Count("net.dev.txbytes", int64(stat.TxBytes), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.txpackets", int64(stat.TxPackets), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.txerrs", int64(stat.TxErrs), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.txdrop", int64(stat.TxDrop), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.txfifo", int64(stat.TxFifo), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.txcolls", int64(stat.TxColls), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.txcarrier", int64(stat.TxCarrier), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
		if err := StatsdClient.Count("net.dev.txcompressed", int64(stat.TxCompressed), []string{"iface:" + stat.Iface}, 1); err != nil {
			Log.WithField("error", err).Error("Couldn't submit event to statsd.")
		}
	}
}
