package conf

import "time"

type ServerConf struct {
	addr          string
	storeInterval int
	restore       bool
}

func NewServerConf(addr string, storeInterval int, restore bool) *ServerConf {
	return &ServerConf{
		addr:          addr,
		storeInterval: storeInterval,
		restore:       restore,
	}
}

func (c ServerConf) Addr() string {
	return c.addr
}

func (c ServerConf) StoreInterval() time.Duration {
	return time.Duration(c.storeInterval) * time.Second
}

func (c ServerConf) Restore() bool {
	return c.restore
}
