package server

import "bitbucket.org/aukbit/pluto/discovery"

// register Server within the service discovery system
func (ds *defaultServer) register() error {
	_, err := discovery.IsAvailable()
	if err != nil {
		ds.logger.Error("service discovery not available")
		return nil
	}
	s := &discovery.Service{
		ID:   ds.cfg.Name,
		Name: ds.cfg.Name,
		Port: ds.cfg.Port(),
		Tags: []string{ds.cfg.Version},
	}
	err = discovery.RegisterService(s)
	if err != nil {
		return err
	}
	// c := &discovery.Check{
	// 	Name:  fmt.Sprintf("Service '%s' check", ds.cfg.Name),
	// 	Notes: fmt.Sprintf("Ensure the server is listening on port %s", ds.cfg.Addr),
	// 	DeregisterCriticalServiceAfter: "10m",
	// 	HTTP:      fmt.Sprintf("http://localhost:%d/_health/server", ds.cfg.Port()),
	// 	Interval:  "10s",
	// 	Timeout:   "1s",
	// 	ServiceID: ds.cfg.Name,
	// }
	// c := &discovery.Check{
	// 	Name:  fmt.Sprintf("Service '%s' check", ds.cfg.Name),
	// 	Notes: fmt.Sprintf("Ensure the server is listening on port %s", ds.cfg.Addr),
	// 	DeregisterCriticalServiceAfter: "10m",
	// 	TCP:       ds.cfg.Addr,
	// 	Interval:  "10s",
	// 	Timeout:   "1s",
	// 	ServiceID: ds.cfg.ID,
	// }
	// err = discovery.RegisterCheck(c)
	if err != nil {
		return err
	}
	ds.isDiscovered = true
	return nil
}

// unregister Server from the service discovery system
func (ds *defaultServer) unregister() {
	defer ds.wg.Done()
	if ds.isDiscovered {
		err := discovery.DeregisterService(ds.cfg.Name)
		if err != nil {
			ds.logger.Error(err.Error())
		}
	}
}
