package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// pkg/clientgo_example/clientset_example.go --config-path=/Users/gutao/gutaodev/gocode/operator-develop/config/kubeconfig-vm1-16

func main() {

	// 参数1：间隔秒数，默认1秒
	interval := flag.Int("interval", 5, "Interval in seconds between pod queries")
	configPath := flag.String("config-path", "", "Path to kubeconfig file")
	flag.Parse()

	// 1. 构建config
	config, err := clientcmd.BuildConfigFromFlags("", *configPath)
	if err != nil {
		panic(err)
	}

	// 2. 创建 ClientSet 对象
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()
	iteration := 0

	for range ticker.C {
		// 使用ClientSet，直接获取已经实现好的client
		pods, err := clientSet.
			CoreV1().
			Pods("kube-system").
			List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		iteration++
		fmt.Println("[print] iteration: ", iteration)
		fmt.Println("[print] pods count: ", len(pods.Items))
		if len(pods.Items) > 0 {
			fmt.Println("[print] pods[0]: ", pods.Items[0].Name, pods.Items[0].Spec.Containers[0].Image)
		}
		fmt.Println("----------------------------------------------------------------------------------------------------------")
		fmt.Println()

	}
}
