package main

import (
	"fmt"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func DiscoveryClient_example() {
	// 1. 构建config
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/gutao/gutaodev/gocode/operator-develop/pkg/clientgo_example/minikube_kubeconfig")
	if err != nil {
		panic(err)
	}

	// 2. 创建discoveryClient
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}

	_, apiResourcesList, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}

	for _, ls := range apiResourcesList {
		for _, apiResources := range ls.APIResources {
			fmt.Printf("groupVersion: %s, group: %s, version: %s, name: %s\n", ls.GroupVersion, apiResources.Group, apiResources.Version, apiResources.Name)
		}
	}
}
