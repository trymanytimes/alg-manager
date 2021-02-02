package main

import (
	"flag"
	"sync"

	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"

	"github.com/trymanytimes/alg-manager/config"
	"github.com/trymanytimes/alg-manager/pkg/etcdwatch"
	"github.com/trymanytimes/alg-manager/pkg/nodes"
	"github.com/trymanytimes/alg-manager/pkg/proto/etcdserverpb"
	"github.com/trymanytimes/alg-manager/pkg/updatedip"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "c", "../etc/controller.conf", "configure file path")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
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
	updatedIPHandler, err := updatedip.NewUpdatedIPHandler(kvClient, conf)
	if err != nil {
		log.Errorf(":%s", err.Error())
	}
	nodesHandler := nodes.NewNodesHandler(kvClient)
	watchUpdatedIPPortClient, err := etcdwatch.NewETCDWatchClient(conn, conf, updatedIPHandler)
	if err != nil {
		log.Errorf("NewETCDWatchClient:%s", err.Error())
		return
	}
	watchNodesClient, err := etcdwatch.NewETCDWatchClient(conn, conf, nodesHandler)
	if err != nil {
		log.Errorf("NewETCDWatchClient:%s", err.Error())
		return
	}
	if err != nil {
		log.Errorf("NewETCDWatchClient:%s", err.Error())
		return
	}
	// watch /project_nm/version1_0_0/cluster/cluster_balance_info/socs.conf for the vip changed.
	if err := watchUpdatedIPPortClient.Watch("/project_nm/version1_0_0/ate/website/domain_id/", "/project_nm/version1_0_0/ate/website/domain_id0"); err != nil {
		log.Errorf("Watch /project_nm/version1_0_0/cluster/cluster_balance_info/socs.conf error:%s", err.Error())
		return
	}
	if err := watchNodesClient.Watch("/project_nm/version1_0_0/nodes/", "/project_nm/version1_0_0/nodes0"); err != nil {
		log.Errorf("Watch /config error:%s", err.Error())
		return
	}
	wg.Wait()
}
