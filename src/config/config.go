package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/goccy/go-yaml"
)

var serverAddressRegex = regexp.MustCompile(`^([^:]+)(:(\d+))?$`)

type serverAddress struct {
	Hostname string
	Port     uint16
}

func (addr serverAddress) String() string {
	if addr.Port == uint16(25565) {
		return addr.Hostname
	}
	return fmt.Sprintf("%v:%v", addr.Hostname, addr.Port)
}

func (addr serverAddress) MarshalYAML() ([]byte, error) {
	if addr.Port == uint16(25565) {
		return []byte(addr.Hostname), nil
	}
	return fmt.Appendf(nil, "%v:%v", addr.Hostname, addr.Port), nil
}

func (addr *serverAddress) UnmarshalYAML(b []byte) error {
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

type ServerConfig struct {
	From []serverAddress `yaml:"from"`
	To   serverAddress   `yaml:"to"`
}

type Configuration struct {
	Servers    map[string]ServerConfig `yaml:"servers"`
	fromToServ map[string]string       `yaml:""`
}

func (conf *Configuration) GetUpstream(hostname string, port uint16) *ServerConfig {
	fromAddr := serverAddress{
		Hostname: hostname,
		Port:     port,
	}

	upstreamName, ok := conf.fromToServ[fromAddr.String()]
	if !ok {
		return nil
	}

	upstream, ok := conf.Servers[upstreamName]
	if !ok {
		return nil
	}

	return &upstream
}

func ReadConfig() *Configuration {
	var conf Configuration

	data, err := os.ReadFile("../config.yml")

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &conf)

	if err != nil {
		panic(err)
	}

	conf.fromToServ = make(map[string]string)

	for serverName, server := range conf.Servers {
		for _, from := range server.From {
			_, ok := conf.fromToServ[from.String()]
			if ok {
				panic(fmt.Errorf("duplicate from address: %v", from.String()))
			}
			conf.fromToServ[from.String()] = serverName
		}
	}

	return &conf
}
