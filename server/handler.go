package server

type handler struct {
	name      string
	target    interface{}
	endpoints []interface{}
}

func newHandler(name string, target interface{}) Handler {
	// typ := reflect.TypeOf(target)
	// hdlr := reflect.ValueOf(target)
	// name := reflect.Indirect(hdlr).Type().Name()
	return &handler{
		name:   name,
		target: target,
	}
}

func (h *handler) Name() string {
	return h.name
}

func (h *handler) Target() interface{} {
	return h.target
}

func (h *handler) Endpoints() []interface{} {
	return h.endpoints
}
