package main

import (
	"flag"
	"fmt"

	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"

	"github.com/trymanytimes/alg-manager/config"
	"github.com/trymanytimes/alg-manager/pkg/etcdwatch"
	"github.com/trymanytimes/alg-manager/pkg/firewalld"
	"github.com/trymanytimes/alg-manager/pkg/proto/etcdserverpb"
	rpcpb "github.com/trymanytimes/alg-manager/pkg/proto/etcdserverpb"
	"github.com/trymanytimes/alg-manager/pkg/updatedip"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "c", "../etc/controller.conf", "configure file path")
	flag.Parse()

	log.InitLogger(log.Debug)
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("load config file failed: %s", err.Error())
	}

	conn, err := grpc.Dial(conf.ETCDAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dail etcd grpc failed: %s", err.Error())
	}
	defer conn.Close()

	kvClient := etcdserverpb.NewKVClient(conn)

	watchClient, err := etcdwatch.NewETCDWatchClient(conn, conf)
	if err != nil {
		log.Errorf("NewETCDWatchClient:%s", err.Error())
		return
	}
	if err := watchClient.Watch("/config", "", updatedip.NewUpdatedIPHandler(kvClient)); err != nil {
		log.Errorf("Watch /config error:%s", err.Error())
		return
	}
	if err := watchClient.Watch("/config/firewalld", "", firewalld.NewFirewalldHandler(kvClient)); err != nil {
		log.Errorf("Watch /config/firewalld error:%s", err.Error())
		return
	}
	var rsp rpcpb.WatchResponse
	for {
		err = watchClient.Watch_watchclient.RecvMsg(&rsp)
		if err != nil {
			log.Fatalf("RecvMsg error,%s", err.Error())
		}
		for _, e := range rsp.Events {
			updatedip.NewUpdatedIPHandler(kvClient).DealEvent(e)
			firewalld.NewFirewalldHandler(kvClient).DealEvent(e)
		}
		fmt.Println(rsp.Events)
	}
}
