package util

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/trymanytimes/alg-manager/pkg/proto/etcdserverpb"
	"google.golang.org/grpc"
)

func GetIPPorts(kvClient etcdserverpb.KVClient) (map[string]int, error) {
	rsp, err := kvClient.Range(context.Background(), &etcdserverpb.RangeRequest{
		Key:      []byte("/project_nm/version1_0_0/ate/website/domain_id/"),
		RangeEnd: []byte("/project_nm/version1_0_0/ate/website/domain_id0"),
	}, grpc.EmptyCallOption{})
	if err != nil {
		return nil, fmt.Errorf("Range error:%s", err.Error())
	}
	allIPPorts := make(map[string]int, 20)
	for _, v := range rsp.Kvs {
		if strings.Index(string(v.Key), "dst_ip_port") > 1 {
			allIPPorts[string(v.Key)[strings.LastIndex(string(v.Key), "/")+1:]] += 1
		}
	}
	return allIPPorts, nil
}

func GetNodeIPs(kvClient etcdserverpb.KVClient) ([]string, error) {
	rsp, err := kvClient.Range(context.Background(), &etcdserverpb.RangeRequest{
		Key:      []byte("/project_nm/version1_0_0/nodes/"),
		RangeEnd: []byte("/project_nm/version1_0_0/nodes0"),
	}, grpc.EmptyCallOption{})
	if err != nil {
		return nil, fmt.Errorf("Range error:%s", err.Error())
	}
	var IPv6Addrs []string
	for _, v := range rsp.Kvs {
		IPRsp, err := kvClient.Range(context.Background(), &etcdserverpb.RangeRequest{
			Key:      []byte("/project_nm/version1_0_0/host/" + string(v.Key)[strings.LastIndex(string(v.Key), "/")+1:] + "/hostinfo/ipv6_mgr_addr"),
			RangeEnd: []byte(""),
		}, grpc.EmptyCallOption{})
		if err != nil {
			return nil, fmt.Errorf("Range error:%s", err.Error())
		}
		if len(IPRsp.Kvs) == 1 {
			IPv6Addrs = append(IPv6Addrs, string(IPRsp.Kvs[0].Value))
		}

	}
	return IPv6Addrs, nil
}

func AddOrRemoveFirewallRule(ip, port string, oper int) error {
	var command string
	if oper == 1 { //add
		command = "firewall-cmd --permanent --add-rich-rule=\"rule family=\"ipv6\" port protocol=\"tcp\" port=\"" + port + "\" accept\""
	} else if oper == 2 { //remove
		command = "firewall-cmd --permanent --remove-rich-rule=\"rule family=\"ipv6\" port protocol=\"tcp\" port=\"" + port + "\" accept\""
	}
	cmd := exec.Command("/bin/bash", "-c", command)
	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("exec %s error: %s", command, err.Error())
	}
	command = "firewall-cmd --reload"
	cmd = exec.Command("/bin/bash", "-c", command)
	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("exec %s error: %s", command, err.Error())
	}
	return nil
}
