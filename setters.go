package datamesh

import (
	"github.com/openziti-incubator/cf"
	"github.com/openziti/dilithium"
	"github.com/openziti/foundation/transport"
	"github.com/pkg/errors"
	"reflect"
)

func TransportAddressSetter(v interface{}, f reflect.Value) error {
	if vt, ok := v.(string); ok {
		addr, err := transport.ParseAddress(vt)
		if err != nil {
			return errors.Wrapf(err, "error parsing '%s'", vt)
		}
		if f.Kind() == reflect.Ptr {
			f.Elem().Set(reflect.ValueOf(addr))
		} else {
			f.Set(reflect.ValueOf(addr))
		}
		return nil
	}
	return errors.Errorf("got '%s', expected '%s'", reflect.TypeOf(v), f.Type())
}

func WestworldProfileFlexibleSetter(v interface{}, opt *cf.Options) (interface{}, error) {
	wp := dilithium.NewBaselineWestworldProfile()
	if err := cf.Bind(wp, v.(map[string]interface{}), opt); err != nil {
		return nil, err
	}
	return wp, nil
}

func ProxyListenerFactorySetter(v interface{}, opt *cf.Options) (interface{}, error) {
	pxlf := &ProxyListenerFactory{}
	if err := cf.Bind(pxlf, v.(map[string]interface{}), opt); err != nil {
		return nil, err
	}
	return pxlf, nil
}

func ProxyTerminatorFactorySetter(v interface{}, opt *cf.Options) (interface{}, error) {
	pxtf := &ProxyTerminatorFactory{}
	if err := cf.Bind(pxtf, v.(map[string]interface{}), opt); err != nil {
		return nil, err
	}
	return pxtf, nil
}