package main

import (
	"time"

	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/tools/clientcmd"
)

func CacheDiscoveryClient_example() {
	// 1. 构建config
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/gutao/gutaodev/gocode/operator-develop/pkg/clientgo_example/minikube_kubeconfig")
	if err != nil {
		panic(err)
	}

	// 2. 创建cachedDiscoveryClient
	cacheDiscoveryClient, err := disk.NewCachedDiscoveryClientForConfig(config, "./cache/discovery", "./cache/http", time.Hour)
	if err != nil {
		panic(err)
	}

	// gvr缓存到当前目录
	cacheDiscoveryClient.ServerGroupsAndResources()
}
