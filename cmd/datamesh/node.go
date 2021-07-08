package main

import (
	"github.com/openziti-incubator/cf"
	"github.com/openziti-incubator/datamesh"
	"github.com/openziti/foundation/transport"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"reflect"
)

func init() {
	rootCmd.AddCommand(nodeCmd)
}

var nodeCmd = &cobra.Command{
	Use:   "node <config>",
	Short: "Start datamesh node",
	Args:  cobra.ExactArgs(1),
	Run:   node,
}

func node(_ *cobra.Command, args []string) {
	cfo := cf.DefaultOptions()
	cfo = cfo.AddSetter(reflect.TypeOf((*transport.Address)(nil)).Elem(), datamesh.TransportAddressSetter)

	cfg := &datamesh.Config{}
	if err := cf.BindYaml(cfg, args[0], cfo); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info(cf.Dump("config", cfg, cfo))
}
