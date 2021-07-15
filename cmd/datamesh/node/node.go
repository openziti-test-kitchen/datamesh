package node

import (
	"github.com/openziti-incubator/cf"
	"github.com/openziti-incubator/datamesh"
	"github.com/openziti-incubator/datamesh/cmd/datamesh/cli"
	"github.com/openziti/foundation/transport"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"reflect"
	"time"
)

func init() {
	cli.RootCmd.AddCommand(nodeCmd)
}

var nodeCmd = &cobra.Command{
	Use:   "node <config>",
	Short: "Start datamesh node",
	Args:  cobra.ExactArgs(1),
	Run:   node,
}

func node(_ *cobra.Command, args []string) {
	cfOpts := cf.DefaultOptions()
	cfOpts = cfOpts.AddSetter(reflect.TypeOf((*transport.Address)(nil)).Elem(), datamesh.TransportAddressSetter)

	cfCfg := &Config{}
	if err := cf.BindYaml(cfCfg, args[0], cfOpts); err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info(cf.Dump(cfCfg, cfOpts))

	d := datamesh.NewDatamesh(cfCfg.Datamesh)
	d.Start()

	for _, peer := range cfCfg.Peers {
		linkCh, err := d.Dial("default", peer)
		if err == nil {
			logrus.Infof("connected link [%s]", linkCh.Id().Token)
		} else {
			logrus.Errorf("error connecting link (%v)", err)
		}
	}

	for {
		time.Sleep(30 * time.Second)
	}
}
