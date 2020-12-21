package selector

import (
	"math"
	"strconv"
	"strings"

	"github.com/ofavor/micro-lite/internal/log"
	"github.com/ofavor/micro-lite/registry"
)

// Options of selector
type Options struct {
	SelectOpts SelectOptions
}

// Option function to set selector options
type Option func(opts *Options)

// Filter service filter function
type Filter func([]*registry.Service) []*registry.Service

// SelectOptions options for select
type SelectOptions struct {
	Filters []Filter
}

// SelectOption function to set select options
type SelectOption func(opts *SelectOptions)

func defaultOptions() Options {
	return Options{}
}

func isStringIn(target string, strs []string) bool {
	for _, s := range strs {
		if target == s {
			return true
		}
	}
	return false
}

func verToInt(ver string) int {
	ss := strings.Split(ver, ".")
	v := 0
	for i := len(ss); i > 0; i-- {
		j, _ := strconv.Atoi(ss[i-1])
		k := math.Pow(10, float64(3-i)*3)
		v += int(k) * j
	}
	return v
}

func isVersionIn(target string, vers []string) bool {
	if len(vers) == 0 {
		return true
	}
	tv := verToInt(target)
	v1 := 0
	v2 := 0
	if len(vers) > 0 {
		v1 = verToInt(vers[0])
	}
	if len(vers) > 1 {
		v2 = verToInt(vers[1])
	}
	if (v1 > 0 && tv < v1) || (v2 > 0 && tv > v2) {
		return false
	}
	return true
}

// WithAddressFilter node address filter
func WithAddressFilter(addrs []string) SelectOption {
	return func(opts *SelectOptions) {
		opts.Filters = append(opts.Filters, func(services []*registry.Service) []*registry.Service {
			ret := []*registry.Service{}
			for _, s := range services {
				ts := &registry.Service{
					Name:      s.Name,
					Version:   s.Version,
					Metadata:  s.Metadata,
					Endpoints: s.Endpoints,
					Nodes:     []*registry.Node{},
				}
				for _, n := range s.Nodes {
					if isStringIn(n.Address, addrs) {
						log.Debugf("Node address match: %s => %v", n.Address, addrs)
						ts.Nodes = append(ts.Nodes, n)
					} else {
						log.Debugf("Node address mismatch: %s => %v", n.Address, addrs)
					}
				}
				if len(ts.Nodes) > 0 {
					ret = append(ret, ts)
				}
			}
			return ret
		})
	}
}

// WithIDFilter service id filter
func WithIDFilter(ids []string) SelectOption {
	return func(opts *SelectOptions) {
		opts.Filters = append(opts.Filters, func(services []*registry.Service) []*registry.Service {
			ret := []*registry.Service{}
			for _, s := range services {
				ts := &registry.Service{
					Name:      s.Name,
					Version:   s.Version,
					Metadata:  s.Metadata,
					Endpoints: s.Endpoints,
					Nodes:     []*registry.Node{},
				}
				for _, n := range s.Nodes {
					if isStringIn(n.ID, ids) {
						log.Debugf("Node ID match: %s => %v", n.ID, ids)
						ts.Nodes = append(ts.Nodes, n)
					} else {
						log.Debugf("Node ID mismatch: %s => %v", n.ID, ids)
					}
				}
				if len(ts.Nodes) > 0 {
					ret = append(ret, ts)
				}
			}
			return ret
		})
	}
}

// WithVersionFilter service version filter
func WithVersionFilter(vers []string) SelectOption {
	return func(opts *SelectOptions) {
		opts.Filters = append(opts.Filters, func(services []*registry.Service) []*registry.Service {
			ret := []*registry.Service{}
			for _, s := range services {
				if isVersionIn(s.Version, vers) {
					log.Debugf("Service version match: %s => %v", s.Version, vers)
					ret = append(ret, s)
				} else {
					log.Debugf("Service version mismatch: %s => %v", s.Version, vers)
				}
			}
			return ret
		})
	}
}
