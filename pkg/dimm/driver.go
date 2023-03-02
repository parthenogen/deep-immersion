package dimm

type driver struct {
	config config

	queries   chan Query
	responses chan Response
	errors    chan error

	stop chan struct{}
}

func NewDriver(c config) (d *driver) {
	d = &driver{
		config:    c,
		queries:   make(chan Query),
		responses: make(chan Response),
		errors:    make(chan error),
		stop:      make(chan struct{}),
	}

	return
}

func (d *driver) Run() {
	var (
		s source
		c dnsClient
		i inspector
	)

	for _, s = range d.config.Sources() {
		go d.driveSource(s)
	}

	for _, c = range d.config.DNSClients() {
		go d.driveDNSClient(c)
	}

	for _, i = range d.config.Inspectors() {
		go d.driveInspector(i)
	}
}

func (d *driver) Stop() {
	d.stop <- struct{}{}

	return
}

func (d *driver) driveSource(s source) {
	for {
		select {
		case <-d.stop:
			return

		default:
			d.queries <- &query{
				fqdn: s.GenerateFQDN(),
			}
		}
	}
}

func (d *driver) driveDNSClient(client dnsClient) {
	var (
		response Response
		e        error
	)

	for {
		select {
		case <-d.stop:
			return

		case <-d.config.Conductor().Beats():
			response, e = client.Send(<-d.queries)

			if e != nil {
				d.errors <- e

			} else {
				d.responses <- response
			}
		}
	}
}

func (d *driver) driveInspector(i inspector) {
	for {
		select {
		case <-d.stop:
			return

		default:
			i.Inspect(<-d.responses)
		}
	}
}
