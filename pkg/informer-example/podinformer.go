package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/labels"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig = "/Users/gutao/gutaodev/gocode/operator-develop/config/kubeconfig-vm1-16"

func BuildInformerFactory() informers.SharedInformerFactory {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}
	factory := informers.NewSharedInformerFactory(clientset, 0)
	return factory
}

func main() {
	stopCh := make(chan struct{})
	defer close(stopCh)

	// 1. 构建InformerFactory
	factory := BuildInformerFactory()

	// 2. 获取对应的PodInformer
	podFactory := factory.Core().V1().Pods()
	// 3. 获取Lister和informer, 使用Factory初始化PodInformer, 实际上是list资源到cache中
	podLister := podFactory.Lister()
	podInformer := podFactory.Informer()
	//podInformer.Informer()  // 这一步不用也行，在informer.Lister()中也会调用Informer()初始化
	// 4. 运行所有初始化过的Informer
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	// 5. 使用informer的api
	pods, err := podLister.List(labels.Everything())
	//pods, err := podLister.Pods("ns1").List(labels.Everything())
	if err != nil {
		fmt.Println("PodLister list pod failed, err:", err)
	}
	fmt.Println("get pod nums: ", len(pods))

	// 使用podinformer
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			if !strings.Contains(pod.Name, "test-app") {
				return
			}
			fmt.Println("Pod added:", pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*corev1.Pod)
			newPod := newObj.(*corev1.Pod)
			if !strings.Contains(oldPod.Name, "test-app") {
				return
			}
			if oldPod.ResourceVersion != newPod.ResourceVersion {
				fmt.Println("Pod resourceVersion 发生更新")
			}
			if !reflect.DeepEqual(oldPod.Status, newPod.Status) {
				fmt.Println("Pod status 状态发生更新")
			}
			if !reflect.DeepEqual(oldPod.Spec, newPod.Spec) {
				fmt.Println("Pod spec 状态发生更新")
			}
			fmt.Println("Pod updated:", oldPod.Name)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			if !strings.Contains(pod.Name, "test-app") {
				return
			}
			fmt.Println("Pod deleted:", pod.Name)
		},
	})

	<-stopCh
}
