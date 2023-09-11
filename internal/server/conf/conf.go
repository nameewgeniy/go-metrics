package conf

type ServerConf struct {
	addr string
}

func NewServerConf(addr string) *ServerConf {
	return &ServerConf{
		addr: addr,
	}
}

func (c ServerConf) Addr() string {
	return c.addr
}
