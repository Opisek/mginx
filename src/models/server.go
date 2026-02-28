package models

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/goccy/go-yaml"
)

var serverAddressRegex = regexp.MustCompile(`^([^:]+)(:(\d+))?$`)

func (addr Address) String() string {
	if addr.Port == uint16(25565) {
		return addr.Hostname
	}
	return fmt.Sprintf("%v:%v", addr.Hostname, addr.Port)
}

func (addr Address) MarshalYAML() ([]byte, error) {
	if addr.Port == uint16(25565) {
		return []byte(addr.Hostname), nil
	}
	return fmt.Appendf(nil, "%v:%v", addr.Hostname, addr.Port), nil
}

func (addr *Address) UnmarshalYAML(b []byte) error {
	var str string
	if err := yaml.Unmarshal(b, &str); err != nil {
		return err
	}

	matches := serverAddressRegex.FindAllStringSubmatch(str, -1)

	if matches == nil || len(matches) != 1 || len(matches[0]) != 4 {
		return fmt.Errorf("malformed hostname or IP address: %v", str)
	}

	if len(matches[0][3]) == 0 {
		addr.Port = 25565
	} else {
		value, err := strconv.ParseUint(matches[0][3], 10, 16)

		if err != nil {
			return errors.Join(errors.New("could not parse port"), err)
		}

		addr.Port = uint16(value)
	}

	addr.Hostname = matches[0][1]

	return nil
}

type Address struct {
	Hostname string
	Port     uint16
}

const (
	ServerStateUnmanaged = iota
	ServerStateUnknown
	ServerStateDown
	ServerStateStarting
	ServerStateUp
	ServerStateStopping
)

type UpstreamServer struct {
	From            []Address `yaml:"from"`
	To              Address   `yaml:"to"`
	Redirect        bool      `yaml:"redirect"`
	connectionState int       `yaml:""`
}
