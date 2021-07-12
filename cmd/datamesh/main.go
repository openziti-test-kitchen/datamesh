package main

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti-incubator/datamesh/cmd/datamesh/cli"
	"github.com/openziti/foundation/transport"
	"github.com/openziti/foundation/transport/tcp"
	"github.com/openziti/foundation/transport/tls"
	"github.com/sirupsen/logrus"
)

func init() {
	pfxlog.GlobalInit(logrus.InfoLevel, pfxlog.DefaultOptions().SetTrimPrefix("github.com/openziti/"))
	transport.AddAddressParser(&tcp.AddressParser{})
	transport.AddAddressParser(&tls.AddressParser{})
}

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		logrus.Fatalf("error (%v)", err)
	}
}
