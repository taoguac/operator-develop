package ingressexample

import (
	"log"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func RunInMain() {
	// 构建和运行ingress的controller
	// 1. 读取config
	config, err := clientcmd.BuildConfigFromFlags("",
		"/Users/gutao/gutaodev/gocode/operator-develop/config/minikube_kubeconfig")
	if err != nil {
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln("get kubenetes config failed")
		}
		config = inClusterConfig
	}
	// 2. 构建client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("get kubernetes clientset failed")
	}

	// 3. 构建informer
	// 4. 向informer注册event handler
	factory := informers.NewSharedInformerFactory(clientset, 0)
	serviceInformer := factory.Core().V1().Services()
	ingressInformer := factory.Networking().V1().Ingresses()
	controller := NewController(clientset, serviceInformer, ingressInformer)

	// 5. 启动informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)
	controller.Run(stopCh)

}
