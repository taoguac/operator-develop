package ingressexample

import (
	"context"
	"reflect"
	"time"

	k8sv1 "k8s.io/api/core/v1"
	k8snetv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	v1informer "k8s.io/client-go/informers/core/v1"
	networkinginformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/listers/core/v1"
	networkingv1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	workNum                     = 5
	maxRetry                    = 10
	ServiceControllerAnnotation = "ingress.controller/example"
)

type controller struct {
	client kubernetes.Interface
	// Lister是Informer的cache
	ingressLister networkingv1.IngressLister
	serviceLister corev1.ServiceLister
	// 是业务相关的workQueue?
	queue workqueue.RateLimitingInterface
}

func NewController(client kubernetes.Interface, serviceInformer v1informer.ServiceInformer, ingInformer networkinginformer.IngressInformer) *controller {
	c := controller{
		client:        client,
		ingressLister: ingInformer.Lister(),
		serviceLister: serviceInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingress_example_manager"),
	}
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})
	ingInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})
	return &c
}

func (c *controller) enqueue(obj interface{}) {
	// queue中默认存放的是service对象
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	c.queue.Add(key)
}

func (c *controller) addService(obj interface{}) {
	c.enqueue(obj)
}

func (c *controller) updateService(oldObj, newObj interface{}) {
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}
	c.enqueue(newObj)
}

// 删除ingress后，如果对应的service还存在，需要重新再拉起来
func (c *controller) deleteIngress(obj interface{}) {
	ingress := obj.(*k8snetv1.Ingress)
	ownerReference := v1.GetControllerOf(ingress)
	if ownerReference == nil {
		return
	}
	if ownerReference.Kind != "Service" {
		return
	}
	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

func (c *controller) Run(stopCh chan struct{}) {
	for i := 0; i < workNum; i++ {
		go wait.Until(c.worker, time.Minute, stopCh)
	}
	<-stopCh
}

func (c *controller) worker() {
	for c.processNextItem() {
		klog.Info("[ingress_controller]")
	}
}

func (c *controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	defer c.queue.Done(item)
	if shutdown {
		return false
	}

	key := item.(string)
	err := c.syncService(key)
	if err != nil {
		c.handlerError(key, err)
	}
	return true
}

func (c *controller) syncService(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// service删除事件
	service, err := c.serviceLister.Services(namespace).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	// service修改事件，ingress删除事件
	_, ok := service.Annotations[ServiceControllerAnnotation]
	ingress, err := c.ingressLister.Ingresses(namespace).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if ok && errors.IsNotFound(err) {
		ig := c.constructIngressOfService(service)
		_, err := c.client.NetworkingV1().Ingresses(namespace).Create(context.Background(), ig, v1.CreateOptions{})
		if err != nil {
			return err
		}
	} else if !ok && ingress != nil {
		err := c.client.NetworkingV1().Ingresses(namespace).Delete(context.Background(), name, v1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *controller) constructIngressOfService(service *k8sv1.Service) *k8snetv1.Ingress {
	ig := &k8snetv1.Ingress{}

	ig.OwnerReferences = []v1.OwnerReference{
		*v1.NewControllerRef(service, k8sv1.SchemeGroupVersion.WithKind("Service")),
	}

	ig.Name = service.Name
	ig.Namespace = service.Namespace
	pathType := k8snetv1.PathTypePrefix
	ingressClassName := "nginx"
	ig.Spec.IngressClassName = &ingressClassName
	ig.Spec = k8snetv1.IngressSpec{
		Rules: []k8snetv1.IngressRule{
			{
				Host: "example.com",
				IngressRuleValue: k8snetv1.IngressRuleValue{
					HTTP: &k8snetv1.HTTPIngressRuleValue{
						Paths: []k8snetv1.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: k8snetv1.IngressBackend{
									Service: &k8snetv1.IngressServiceBackend{
										Name: service.Name,
										Port: k8snetv1.ServiceBackendPort{
											Number: 3333,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return ig
}

func (c *controller) handlerError(key string, err error) {
	if c.queue.NumRequeues(key) < maxRetry {
		c.queue.AddRateLimited(key)
		return
	}
	runtime.HandleError(err)
	c.queue.Forget(key)
}
