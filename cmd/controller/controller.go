package main

import (
	"flag"

	"github.com/zdnscloud/cement/log"

	"github.com/trymanytimes/alg-manager/config"
	rpcpb "github.com/trymanytimes/alg-manager/proto/etcdserverpb"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "c", "web-controller.conf", "configure file path")
	flag.Parse()

	log.InitLogger(log.Debug)
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("load config file failed: %s", err.Error())
	}
	watch := rpcpb.WatchRequest{RequestUnion: &rpcpb.WatchRequest_CreateRequest{}}
}
