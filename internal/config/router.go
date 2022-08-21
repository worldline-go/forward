package config

type TCPListen struct {
	Name string `cfg:"name"`
	Host string `cfg:"host"`
}

type TCPConnect struct {
	Name string `cfg:"name"`
	URL  string `cfg:"url"`
}

type Forward struct {
	Name string `cfg:"name"`
	From string `cfg:"from"`
	To   string `cfg:"to"`
}

// Router is the configuration for the router.
type Router struct {
	TCPListen  []TCPListen  `cfg:"tcpListen"`
	TCPConnect []TCPConnect `cfg:"tcpConnect"`

	Forward []Forward `cfg:"forward"`
}
