package framework

import (
	"github.com/appscode/go/crypto/rand"
	tcs "github.com/k8sdb/apimachinery/client/typed/kubedb/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type Framework struct {
	kubeClient   kubernetes.Interface
	extClient    tcs.KubedbV1alpha1Interface
	namespace    string
	name         string
	StorageClass string
}

func New(kubeClient kubernetes.Interface, extClient tcs.KubedbV1alpha1Interface, storageClass string) *Framework {
	return &Framework{
		kubeClient:   kubeClient,
		extClient:    extClient,
		name:         "memcached-operator",
		namespace:    rand.WithUniqSuffix("memcached"),
		StorageClass: storageClass,
	}
}

func (f *Framework) Invoke() *Invocation {
	return &Invocation{
		Framework: f,
		app:       rand.WithUniqSuffix("memcached-e2e"),
	}
}

type Invocation struct {
	*Framework
	app string
}
