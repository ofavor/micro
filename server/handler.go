package server

import (
	"reflect"
	"strings"

	"github.com/ofavor/micro-lite/internal/log"
	"github.com/ofavor/micro-lite/registry"
)

type handler struct {
	name      string
	target    interface{}
	endpoints []*registry.Endpoint
}

func newHandler(name string, target interface{}) Handler {
	typ := reflect.TypeOf(target)

	var endpoints []*registry.Endpoint

	for m := 0; m < typ.NumMethod(); m++ {
		if e := extractEndpoint(typ.Method(m)); e != nil {
			e.Name = name + "." + e.Name
			log.Debug("Handler endpoint:", e)
			endpoints = append(endpoints, e)
		}
	}
	return &handler{
		name:      name,
		target:    target,
		endpoints: endpoints,
	}
}

func extractValue(v reflect.Type, d int) *registry.Value {
	if d == 3 {
		return nil
	}
	if v == nil {
		return nil
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	arg := &registry.Value{
		Name: v.Name(),
		Type: v.Name(),
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			val := extractValue(f.Type, d+1)
			if val == nil {
				continue
			}

			// if we can find a json tag use it
			if tags := f.Tag.Get("json"); len(tags) > 0 {
				parts := strings.Split(tags, ",")
				if parts[0] == "-" || parts[0] == "omitempty" {
					continue
				}
				val.Name = parts[0]
			}

			// if there's no name default it
			if len(val.Name) == 0 {
				val.Name = v.Field(i).Name
			}

			arg.Values = append(arg.Values, val)
		}
	case reflect.Slice:
		p := v.Elem()
		if p.Kind() == reflect.Ptr {
			p = p.Elem()
		}
		arg.Type = "[]" + p.Name()
	}

	return arg
}

func extractEndpoint(method reflect.Method) *registry.Endpoint {
	if method.PkgPath != "" {
		return nil
	}

	var rspType, reqType reflect.Type
	mt := method.Type

	switch mt.NumIn() {
	case 4:
		reqType = mt.In(2)
		rspType = mt.In(3)
	default:
		return nil
	}

	request := extractValue(reqType, 0)
	response := extractValue(rspType, 0)

	ep := &registry.Endpoint{
		Name:     method.Name,
		Request:  request,
		Response: response,
	}
	return ep
}

func (h *handler) Name() string {
	return h.name
}

func (h *handler) Target() interface{} {
	return h.target
}

func (h *handler) Endpoints() []*registry.Endpoint {
	return h.endpoints
}
