package main

import (
	"github.com/openziti-incubator/cf"
	"github.com/openziti-incubator/datamesh"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node <config>",
	Short: "Start datamesh node",
	Args:  cobra.ExactArgs(1),
	Run:   node,
}

func node(_ *cobra.Command, args []string) {
	cfg := &datamesh.Config{}
	if err := cf.BindYaml(cfg, args[0], cf.DefaultOptions()); err != nil {
		logrus.Error(err)
		return
	}
}
