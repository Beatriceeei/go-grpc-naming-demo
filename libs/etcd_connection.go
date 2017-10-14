package libs

import (
	"time"

	"github.com/coreos/etcd/clientv3"
)

// your etcd config
var (
	endPoints = []string{
		"http://localhost:2379",
		"http://localhost:2378",
		"http://localhost:2377",
	}
	dialTimeOut = 3
	username    = "username"
	password    = "password"
)

// singleTon for etcd connection
var etcdCli *clientv3.Client

// GetEtcdCli is a method for getting etcd connection
func GetEtcdCli() *clientv3.Client {

	if etcdCli == nil {
		config := clientv3.Config{
			Endpoints:   endPoints,
			DialTimeout: time.Duration(dialTimeOut) * time.Second,
			Username:    username,
			Password:    password,
		}

		var err error
		etcdCli, err = clientv3.New(config)
		if err != nil {
			panic(err)
		}
	}

	return etcdCli
}
