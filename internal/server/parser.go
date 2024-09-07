package server

import (
	"strings"
)

type Holder struct {
	Name   string
	Host   string
	Socket []string
}

func Parse(hosts, sockets []string) []Holder {
	var holders []Holder

	socketList := make(map[string][]string)

	for _, socket := range sockets {
		name, sValue := atParse(socket)
		socketList[name] = append(socketList[name], sValue)
	}

	for _, host := range hosts {
		name, addr := atParse(host)

		holders = append(holders, Holder{
			Name:   name,
			Host:   addr,
			Socket: socketList[name],
		})
	}

	return holders
}

func atParse(v string) (name string, value string) {
	vSplit := strings.SplitN(v, "@", 2)

	switch {
	case len(vSplit) < 2:
		name = "default"
		value = vSplit[0]
	default:
		name = vSplit[0]
		value = vSplit[1]
	}

	return
}
