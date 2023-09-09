package conf

type ServerConf struct {
	addr string
}

func NewServerConf() *ServerConf {
	return &ServerConf{
		addr: ":8080",
	}
}

func (c ServerConf) Addr() string {
	return c.addr
}
