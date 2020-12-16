package selector

// Options of selector
type Options struct {
}

// Option function to set selector options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{}
}
