package models

type Compose struct {
	Version  string
	Networks map[string]Network
	Volumes  map[string]Volume
	Services map[string]Service
}

type Network struct {
	Driver   string
	External string
}

type Volume struct {
	Driver   string
	External string
}

type Service struct {
	Image                                    string
	Networks, Ports, Volumes, Command, Links []string
	Environment                              []string
}

func (c *Compose) GetServiceNames() []string {
	keys := make([]string, 0, len(c.Services))
	for k := range c.Services {
		keys = append(keys, k)
	}
	return keys
}
