package config

import (
	"fmt"
	"mginx/models"
	"os"

	"github.com/goccy/go-yaml"
)

type Configuration struct {
	Servers    map[string]*models.UpstreamServer `yaml:"servers"`
	fromToServ map[string]string                 `yaml:""`
}

func (conf *Configuration) GetUpstream(hostname string, port uint16) *models.UpstreamServer {
	fromAddr := models.Address{
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

	return upstream
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
		server.InternalName = serverName
		if server.Watchdog.IsManaged() {
			server.SetUnknown()
		}

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
