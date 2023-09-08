package conf

type Conf struct {
	addr string
}

func NewConf() *Conf {
	return &Conf{
		addr: ":8080",
	}
}

func (c Conf) Addr() string {
	return c.addr
}
