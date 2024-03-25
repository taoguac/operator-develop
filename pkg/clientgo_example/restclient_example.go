package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// 使用restclient访问资源

func Restclient_example() {
	// 1. 构建config
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/gutao/gutaodev/gocode/operator-develop/pkg/clientgo_example/minikube_kubeconfig")
	if err != nil {
		panic(err)
	}

	// 2. 配置config
	config.APIPath = "api"
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	// 3. 创建restclient
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	result := &v1.PodList{}
	err = restClient.Get().
		Namespace("kube-system").
		Resource("pods").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.Background()).
		Into(result)
	if err != nil {
		panic(err)
	}

	for _, pod := range result.Items {
		fmt.Printf("namespace: %s, name: %s\n", pod.Namespace, pod.Name)
	}

}
