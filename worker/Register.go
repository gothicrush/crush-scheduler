package worker

import (
	"context"
	"github.com/gothicrush/crush-scheduler/common"
	"go.etcd.io/etcd/clientv3"
	"net"
	"time"
)

// 注册节点到etcd
type Register struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	localIP string //本地ip
}

var (
	// 单例
	G_register *Register
)

// 初始化服务注册
func InitRegister() error {

	// etcd 连接配置
	config := clientv3.Config{
		Endpoints:   G_config.EtcdEndPoints,
		DialTimeout: time.Duration(G_config.EtcdTimeout) * time.Millisecond,
	}

	// 建立连接
	client, err := clientv3.New(config)

	if err != nil {
		return err
	}

	// 得到kv客户端
	kv := clientv3.NewKV(client)

	// 得到lease客户端
	lease := clientv3.NewLease(client)

	G_register = &Register{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	// 获取本机网卡的IP地址，作为唯一标识
	ip, err := getLocalIP()

	if err != nil {
		return err
	}

	G_register.localIP = ip

	// 服务注册
	go G_register.keepOnline()

	return nil
}

// 注册到 /cron/workers/IP，并自动续租
func (register *Register) keepOnline() {

	// 注册路径
	regKey := common.JOB_WORKER_DIR + register.localIP

	var cancelCtx context.Context
	var cancelFunc context.CancelFunc = nil

	var leaseResp *clientv3.LeaseGrantResponse
	var err error
	var keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse

	for {

		// 创建租约
		leaseResp, err = register.lease.Grant(context.TODO(), 10)

		if err != nil {
			continue
		}

		// 自动续租
		keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseResp.ID)

		if err != nil {
			continue
		}

		cancelCtx, cancelFunc = context.WithCancel(context.TODO())

		// 注册到etcd
		_, err = register.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseResp.ID))

		if err != nil {
			continue
		}

		// 处理续租应答
		for {
			select {
			case keepResp := <-keepAliveChan:
				if keepResp == nil { // 续租失败
					goto RETRY
				}
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			cancelFunc()
		}
	}

}

func getLocalIP() (string, error) {

	var ipv4 string
	var err error

	// 获取所有网卡
	arr, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	// 取第一个非Localhost的IP
	for _, addr := range arr {
		ipNet, isIPNet := addr.(*net.IPNet)

		if isIPNet && ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				break
			}
		}
	}

	return ipv4, nil
}
