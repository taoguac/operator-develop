package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func DynamicClient_example() {
	// 1. 构建config
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/gutao/gutaodev/gocode/operator-develop/pkg/clientgo_example/minikube_kubeconfig")
	if err != nil {
		panic(err)
	}

	// 2. 创建 DynamicClient
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 3. 配置GVR
	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	// 4. 获取unStructured非结构化数据
	unStructData, err := dynamicClient.Resource(gvr).Namespace("kube-system").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// 5. 将unStructured转换为结构化数据
	pods := &v1.PodList{}
	_ = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructData.UnstructuredContent(), pods)

	for _, pod := range pods.Items {
		fmt.Printf("namespace: %s, name: %s\n", pod.Namespace, pod.Name)
	}

}
