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
	cfO := cf.DefaultOptions()
	cfO = cfO.AddSetter(reflect.TypeOf((*transport.Address)(nil)).Elem(), datamesh.TransportAddressSetter)

	cfI := &Config{}
	if err := cf.BindYaml(cfI, args[0], cfO); err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info(cf.Dump(cfI, cfO))

	d := datamesh.NewDatamesh(cfI.Datamesh)
	d.Start()

	for _, peer := range cfI.Peers {
		linkCh, err := d.Dial("default", peer)
		if err == nil {
			logrus.Infof("connected link [%s]", linkCh.Label())
		} else {
			logrus.Errorf("error connecting link (%v)", err)
		}
	}

	for {
		time.Sleep(30 * time.Second)
	}
}
