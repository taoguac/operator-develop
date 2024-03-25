package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Watch_example() {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/gutao/gutaodev/gocode/operator-develop/pkg/clientgo_example/minikube_kubeconfig")
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 调用watch（监听）方法
	w, err := clientset.AppsV1().Deployments("ns1").Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for {
		select {
		case e, _ := <-w.ResultChan():
			fmt.Println(e.Type, e.Object)
		}
	}
}
