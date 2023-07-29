package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/realmicro/realmicro"
	"github.com/realmicro/realmicro/client"
	mconfig "github.com/realmicro/realmicro/config"
	"github.com/realmicro/realmicro/config/encoder/yaml"
	"github.com/realmicro/realmicro/config/reader"
	"github.com/realmicro/realmicro/config/reader/json"
	"github.com/realmicro/realmicro/config/source/file"
	"github.com/realmicro/realmicro/registry"
	"github.com/realmicro/realmicro/registry/etcd"
	"github.com/realmicro/realmicro/wrapper/select/dc"
	log "github.com/sirupsen/logrus"
	"qingyun/common/wrapper"
	"qingyun/services/fishing_gateway/common"
	"qingyun/services/fishing_gateway/internal/config"
	"qingyun/services/fishing_gateway/internal/logic"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	flag.Parse()
	fmt.Println(*configFile)
	c, _ := mconfig.NewConfig(
		mconfig.WithReader(
			json.NewReader(
				reader.WithEncoder(yaml.NewEncoder()),
			),
		),
	)
	var err error
	// load the config from a file source
	if err = c.Load(file.NewSource(
		file.WithPath(*configFile),
	)); err != nil {
		fmt.Println(err)
		return
	}
	var cfg config.Config
	if err = c.Scan(&cfg); err != nil {
		fmt.Println(err)
		return
	}

	service := realmicro.NewService(
		realmicro.Name(cfg.ServiceName),
		realmicro.Version(cfg.Version),
		realmicro.Metadata(map[string]string{
			"env":     cfg.Env,
			"project": cfg.Project,
		}),
		realmicro.Registry(etcd.NewRegistry(registry.Addrs(cfg.Hosts.Etcd.Address...))),
		realmicro.WrapHandler(wrapper.LogHandler()),
		realmicro.WrapClient(dc.NewDCWrapper, wrapper.LogCall),
	)
	client.DefaultClient = client.NewClient(
		client.Registry(etcd.NewRegistry(registry.Addrs(cfg.Hosts.Etcd.Address...))),
	)

	service.Init()
	r := gin.Default()
	r.Use(common.Logger())
	r.GET("/", logic.Home)
	r.POST("/rpc", logic.Rpc)
	log.Info("start.....")
	go r.Run()
	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
