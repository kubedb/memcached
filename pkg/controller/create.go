package controller

import (
	"fmt"
	"time"

	"github.com/appscode/go/types"
	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/k8sdb/apimachinery/pkg/docker"
	apps "k8s.io/api/apps/v1beta1"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// Duration in Minute
	// Check whether pod under StatefulSet is running or not
	// Continue checking for this duration until failure
	durationCheckStatefulSet = time.Minute * 30
)

func (c *Controller) findService(memcached *tapi.Memcached) (bool, error) {
	name := memcached.OffshootName()
	service, err := c.Client.CoreV1().Services(memcached.Namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	if service.Spec.Selector[tapi.LabelDatabaseName] != name {
		return false, fmt.Errorf(`Intended service "%v" already exists`, name)
	}

	return true, nil
}

func (c *Controller) createService(memcached *tapi.Memcached) error {
	svc := &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   memcached.OffshootName(),
			Labels: memcached.OffshootLabels(),
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{
					Name:       "db",
					Port:       11211,
					TargetPort: intstr.FromString("db"),
				},
			},
			Selector: memcached.OffshootLabels(),
		},
	}
	if memcached.Spec.Monitor != nil &&
		memcached.Spec.Monitor.Agent == tapi.AgentCoreosPrometheus &&
		memcached.Spec.Monitor.Prometheus != nil {
		svc.Spec.Ports = append(svc.Spec.Ports, core.ServicePort{
			Name:       tapi.PrometheusExporterPortName,
			Port:       tapi.PrometheusExporterPortNumber,
			TargetPort: intstr.FromString(tapi.PrometheusExporterPortName),
		})
	}

	if _, err := c.Client.CoreV1().Services(memcached.Namespace).Create(svc); err != nil {
		return err
	}

	return nil
}

func (c *Controller) findStatefulSet(memcached *tapi.Memcached) (bool, error) {
	// SatatefulSet for Memcached database
	statefulSet, err := c.Client.AppsV1beta1().StatefulSets(memcached.Namespace).Get(memcached.OffshootName(), metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	if statefulSet.Labels[tapi.LabelDatabaseKind] != tapi.ResourceKindMemcached {
		return false, fmt.Errorf(`Intended statefulSet "%v" already exists`, memcached.OffshootName())
	}

	return true, nil
}

func (c *Controller) createStatefulSet(memcached *tapi.Memcached) (*apps.StatefulSet, error) {
	// SatatefulSet for Memcached database
	statefulSet := &apps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        memcached.OffshootName(),
			Namespace:   memcached.Namespace,
			Labels:      memcached.StatefulSetLabels(),
			Annotations: memcached.StatefulSetAnnotations(),
		},
		Spec: apps.StatefulSetSpec{
			Replicas:    types.Int32P(1),
			ServiceName: c.opt.GoverningService,
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: memcached.OffshootLabels(),
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            tapi.ResourceNameMemcached,
							Image:           fmt.Sprintf("%s:%s", docker.ImageMemcached, memcached.Spec.Version),
							ImagePullPolicy: core.PullIfNotPresent,
							Ports: []core.ContainerPort{
								{
									Name:          "db",
									ContainerPort: 11211,
								},
							},
							Resources: memcached.Spec.Resources,
						},
					},
					NodeSelector:  memcached.Spec.NodeSelector,
					Affinity:      memcached.Spec.Affinity,
					SchedulerName: memcached.Spec.SchedulerName,
					Tolerations:   memcached.Spec.Tolerations,
				},
			},
		},
	}

	if memcached.Spec.Monitor != nil &&
		memcached.Spec.Monitor.Agent == tapi.AgentCoreosPrometheus &&
		memcached.Spec.Monitor.Prometheus != nil {
		exporter := core.Container{
			Name: "exporter",
			Args: []string{
				"export",
				fmt.Sprintf("--address=:%d", tapi.PrometheusExporterPortNumber),
				"--v=3",
			},
			Image:           docker.ImageOperator + ":" + c.opt.ExporterTag,
			ImagePullPolicy: core.PullIfNotPresent,
			Ports: []core.ContainerPort{
				{
					Name:          tapi.PrometheusExporterPortName,
					Protocol:      core.ProtocolTCP,
					ContainerPort: int32(tapi.PrometheusExporterPortNumber),
				},
			},
		}
		statefulSet.Spec.Template.Spec.Containers = append(statefulSet.Spec.Template.Spec.Containers, exporter)
	}

	if c.opt.EnableRbac {
		// Ensure ClusterRoles for database statefulsets
		if err := c.createRBACStuff(memcached); err != nil {
			return nil, err
		}

		statefulSet.Spec.Template.Spec.ServiceAccountName = memcached.Name
	}

	if _, err := c.Client.AppsV1beta1().StatefulSets(statefulSet.Namespace).Create(statefulSet); err != nil {
		return nil, err
	}

	return statefulSet, nil
}

func (c *Controller) createDormantDatabase(memcached *tapi.Memcached) (*tapi.DormantDatabase, error) {
	dormantDb := &tapi.DormantDatabase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcached.Name,
			Namespace: memcached.Namespace,
			Labels: map[string]string{
				tapi.LabelDatabaseKind: tapi.ResourceKindMemcached,
			},
		},
		Spec: tapi.DormantDatabaseSpec{
			Origin: tapi.Origin{
				ObjectMeta: metav1.ObjectMeta{
					Name:        memcached.Name,
					Namespace:   memcached.Namespace,
					Labels:      memcached.Labels,
					Annotations: memcached.Annotations,
				},
				Spec: tapi.OriginSpec{
					Memcached: &memcached.Spec,
				},
			},
		},
	}

	return c.ExtClient.DormantDatabases(dormantDb.Namespace).Create(dormantDb)
}

func (c *Controller) reCreateMemcached(memcached *tapi.Memcached) error {
	_memcached := &tapi.Memcached{
		ObjectMeta: metav1.ObjectMeta{
			Name:        memcached.Name,
			Namespace:   memcached.Namespace,
			Labels:      memcached.Labels,
			Annotations: memcached.Annotations,
		},
		Spec:   memcached.Spec,
		Status: memcached.Status,
	}

	if _, err := c.ExtClient.Memcacheds(_memcached.Namespace).Create(_memcached); err != nil {
		return err
	}

	return nil
}
