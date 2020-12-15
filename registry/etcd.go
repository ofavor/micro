package registry

type etcdRegistry struct {
}

func newETCDRegistry(opts ...Option) Registry {
	return &etcdRegistry{}
}

func (r *etcdRegistry) Register(svc *Service) error {
	return nil
}
func (r *etcdRegistry) Deregister(svc *Service) error {
	return nil
}
