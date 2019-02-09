package grpcserver

type Option func(c Config) Config

type Config struct {
	Servers []Server
	Port    int64
}

func WithServers(servers ...Server) Option {
	return func(c Config) Config {
		c.Servers = append(c.Servers, servers...)
		return c
	}
}

func WithPort(port int64) Option {
	return func(c Config) Config {
		c.Port = port
		return c
	}
}
