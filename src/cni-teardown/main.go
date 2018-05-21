package main

import (
	"cni-teardown/config"
	"flag"
	"os"
	"strings"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/silk/lib/adapter"
)

func main() {
	logger := lager.NewLogger("cni-teardown")
	sink := lager.NewWriterSink(os.Stdout, lager.INFO)
	logger.RegisterSink(sink)

	logger.Info("starting")
	netlinkAdapter := &adapter.NetlinkAdapter{}

	links, err := netlinkAdapter.LinkList()
	if err != nil {
		logger.Error("failed-to-list-network-devices", err) // not tested
	}

	for _, link := range links {
		if link.Type() == "ifb" && strings.HasPrefix(link.Attrs().Name, "i") {
			err = netlinkAdapter.LinkDel(link)
			if err != nil {
				logger.Error("failed-to-remove-ifb", err)
			}
		}
	}

	configFilePath := flag.String("config", "", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configFilePath)
	if err != nil {
		logger.Error("read-config-file", err)
		os.Exit(1)
	}

	for _, path := range cfg.PathsToDelete {
		err := os.RemoveAll(path)
		if err != nil {
			logger.Info("failed-to-remove-path", lager.Data{"path": path, "err": err})
		}
	}

	logger.Info("complete")
}
