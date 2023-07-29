package config

import (
	"sync"

	"github.com/realmicro/realmicro/config"
	"github.com/realmicro/realmicro/config/source/etcd"
	log "github.com/sirupsen/logrus"
)

var (
	ServiceName          = "FishingGateway"
	gatewayServer        = "gatewayServer"
	gatewayMethod        = "gatewayMethod"
	gatewayClientVersion = "gatewayClientVersion"
	etcdHost             = []string{"127.0.0.1:2379"}
	once                 sync.Once
	instance             *GoConfig
)

type TplReply struct {
	Score string
}

type GoConfig struct {
	conf config.Config
}

func GetGoConfig() *GoConfig {
	once.Do(func() {
		instance = &GoConfig{}
		instance.conf, _ = config.NewConfig()
		es := etcd.NewSource(
			etcd.WithAddress(etcdHost...),
			etcd.WithPrefix(ServiceName+"/"),
			etcd.StripPrefix(true),
		)
		instance.conf.Load(es)
	})

	return instance
}

func (gc *GoConfig) GetGatewayServerWhiteList() (map[string]int64, error) {
	var rsp map[string]int64
	if err := gc.conf.Get(gatewayServer, "version").Scan(&rsp); err != nil {
		log.Infof("load version whitelist controller reply scan error: %v", err)
		return nil, err
	}
	return rsp, nil
}

func (gc *GoConfig) GetClientGatewayTokenControllerWhiteList() (map[string]int64, error) {
	var rsp map[string]int64
	if err := gc.conf.Get(gatewayClientVersion, "version").Scan(&rsp); err != nil {
		log.Info("load version whitelist controller reply scan error: %v", err)
		return nil, err
	}
	return rsp, nil
}

func (gc *GoConfig) GetGatewayTokenController() (map[string]int64, error) {
	var rsp map[string]int64
	if err := gc.conf.Get(gatewayMethod, "version").Scan(&rsp); err != nil {
		log.Info("load version whitelist controller reply scan error: %v", err)
		return nil, err
	}
	return rsp, nil
}
