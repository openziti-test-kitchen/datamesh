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
	cfo := cf.DefaultOptions()
	cfo = cfo.AddSetter(reflect.TypeOf((*transport.Address)(nil)).Elem(), datamesh.TransportAddressSetter)
	cfo = cfo.AddInstantiator(reflect.TypeOf(datamesh.DialerConfig{}), func() interface{} { return datamesh.DialerConfigDefaults() })
	cfo = cfo.AddInstantiator(reflect.TypeOf(datamesh.ListenerConfig{}), func() interface{} { return datamesh.ListenerConfigDefaults() })

	cfg := &Config{}
	if err := cf.BindYaml(cfg, args[0], cfo); err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info(cf.Dump(cfg, cfo))

	d := datamesh.NewDatamesh(cfg.Datamesh)
	d.Handlers.AddLinkUpHandler(func(l datamesh.Link) {
		logrus.Infof("would create new circuit with [link/%s]", l.Address())
	})
	d.Start()

	for _, peer := range cfg.Peers {
		l, err := d.DialLink("default", peer)
		if err == nil {
			logrus.Infof("connected link [link/%s]", l.Address())
		} else {
			logrus.Errorf("error connecting link (%v)", err)
		}
	}

	for {
		time.Sleep(30 * time.Second)
	}
}
