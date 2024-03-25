package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func SharedInformer_Example() {
	// 1. 构建SharedInformer需要的Client
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/gutao/gutaodev/gocode/operator-develop/pkg/clientgo_example/minikube_kubeconfig")
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 2. 构建SharedInformerFactory
	sharedInformerFactory := informers.NewSharedInformerFactory(clientSet, 0)
	// 3. 生成pod的Informer
	podInformer := sharedInformerFactory.Core().V1().Pods().Informer()

	// 4. 启动Informer
	sharedInformerFactory.Start(nil)
	// 5. 等待Informer数据同步
	sharedInformerFactory.WaitForCacheSync(nil)

	// 使用indexer获取数据
	indexer := podInformer.GetIndexer()
	pods := indexer.List()

	for _, pod := range pods {
		fmt.Printf("pod: %s, %s\n", pod.(*v1.Pod).Namespace, pod.(*v1.Pod).Name)
	}
}
