package controller

import (
	"fmt"
	"reflect"

	"github.com/appscode/go/log"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/kubedb/apimachinery/client/typed/kubedb/v1alpha1/util"
	"github.com/kubedb/apimachinery/pkg/eventer"
	"github.com/kubedb/memcached/pkg/validator"
	"github.com/the-redback/go-oneliners"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Controller) create(memcached *api.Memcached) error {
	if memcached.Status.CreationTime == nil {
		mc, err := util.PatchMemcached(c.ExtClient, memcached, func(in *api.Memcached) *api.Memcached {
			t := metav1.Now()
			in.Status.CreationTime = &t
			in.Status.Phase = api.DatabasePhaseCreating
			return in
		})

		if err != nil {
			c.recorder.Eventf(
				memcached.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToUpdate,
				err.Error(),
			)
			return err
		}
		memcached.Status = mc.Status
	}

	// Assign Default Monitoring Port
	if err := c.setMonitoringPort(memcached); err != nil {
		return err
	}

	if err := validator.ValidateMemcached(c.Client, memcached); err != nil {
		c.recorder.Event(
			memcached.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonInvalid,
			err.Error(),
		)
		return err
	}

	// Event for successful validation
	c.recorder.Event(
		memcached.ObjectReference(),
		core.EventTypeNormal,
		eventer.EventReasonSuccessfulValidate,
		"Successfully validate Memcached",
	)

	// Check DormantDatabase
	if matched, err := c.matchDormantDatabase(memcached); err != nil || matched {
		return err
	}

	// Event for notification that kubernetes objects are creating
	//c.recorder.Event(
	//	memcached.ObjectReference(),
	//	core.EventTypeNormal,
	//	eventer.EventReasonCreating,
	//	"Updating Kubernetes objects",
	//)

	// ensure database Service
	ok1, er1 := c.ensureService(memcached)
	if er1 != nil {
		return er1
	}

	// ensure database Deployment
	ok2, er2 := c.ensureDeployment(memcached)
	if er2 != nil {
		return er2
	}

	if ok1 && ok2 {
		c.recorder.Event(
			memcached.ObjectReference(),
			core.EventTypeNormal,
			eventer.EventReasonSuccessfulCreate,
			"Successfully updated Memcached",
		)
	}

	if ok, err := c.manageMonitor(memcached); err != nil {
		c.recorder.Eventf(
			memcached.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToCreate,
			"Failed to manage monitoring system. Reason: %v",
			err,
		)
		log.Errorln(err)
		return nil
	} else if memcached.Spec.Monitor != nil && ok {
		c.recorder.Event(
			memcached.ObjectReference(),
			core.EventTypeNormal,
			eventer.EventReasonSuccessfulCreate,
			"Successfully updated monitoring system.",
		)
	} else if ok {
		c.recorder.Event(
			memcached.ObjectReference(),
			core.EventTypeNormal,
			eventer.EventReasonSuccessfulCreate,
			"Successfully deleted monitoring system.",
		)
	}
	return nil
}

func (c *Controller) setMonitoringPort(memcached *api.Memcached) error {
	if memcached.Spec.Monitor != nil &&
		memcached.Spec.Monitor.Prometheus != nil {
		if memcached.Spec.Monitor.Prometheus.Port == 0 {
			mc, err := util.PatchMemcached(c.ExtClient, memcached, func(in *api.Memcached) *api.Memcached {
				in.Spec.Monitor.Prometheus.Port = api.PrometheusExporterPortNumber
				return in
			})

			if err != nil {
				c.recorder.Eventf(
					memcached.ObjectReference(),
					core.EventTypeWarning,
					eventer.EventReasonFailedToUpdate,
					err.Error(),
				)
				return err
			}
			memcached.Spec = mc.Spec
		}
	}
	return nil
}

func (c *Controller) matchDormantDatabase(memcached *api.Memcached) (bool, error) {
	// Check if DormantDatabase exists or not
	dormantDb, err := c.ExtClient.DormantDatabases(memcached.Namespace).Get(memcached.Name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			c.recorder.Eventf(
				memcached.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToGet,
				`Fail to get DormantDatabase: "%v". Reason: %v`,
				memcached.Name,
				err,
			)
			return false, err
		}
		return false, nil
	}

	var sendEvent = func(message string, args ...interface{}) (bool, error) {
		c.recorder.Eventf(
			memcached.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToCreate,
			message,
			args,
		)
		return false, fmt.Errorf(message, args)
	}

	// Check DatabaseKind
	if dormantDb.Labels[api.LabelDatabaseKind] != api.ResourceKindMemcached {
		return sendEvent(fmt.Sprintf(`Invalid Memcached: "%v". Exists DormantDatabase "%v" of different Kind`, memcached.Name, dormantDb.Name))
	}

	// Check Origin Spec
	drmnOriginSpec := dormantDb.Spec.Origin.Spec.Memcached
	originalSpec := memcached.Spec

	if !reflect.DeepEqual(drmnOriginSpec, &originalSpec) {
		return sendEvent("Memcached spec mismatches with OriginSpec in DormantDatabases")
	}

	if err := c.ExtClient.Memcacheds(memcached.Namespace).Delete(memcached.Name, &metav1.DeleteOptions{}); err != nil {
		return sendEvent(`failed to resume Memcached "%v" from DormantDatabase "%v". Error: %v`, memcached.Name, memcached.Name, err)
	}

	_, err = util.PatchDormantDatabase(c.ExtClient, dormantDb, func(in *api.DormantDatabase) *api.DormantDatabase {
		in.Spec.Resume = true
		return in
	})
	if err != nil {
		c.recorder.Eventf(
			memcached.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToUpdate,
			err.Error(),
		)
		return sendEvent(err.Error())
	}
	return true, nil
}

func (c *Controller) pause(memcached *api.Memcached) error {
	oneliners.PrettyJson(memcached, "memcached")
	c.recorder.Event(
		memcached.ObjectReference(),
		core.EventTypeNormal,
		eventer.EventReasonPausing,
		"Pausing Memcached",
	)

	// Assign default monitoring port
	if memcached.Spec.Monitor != nil &&
		memcached.Spec.Monitor.Prometheus != nil {
		if memcached.Spec.Monitor.Prometheus.Port == 0 {
			memcached.Spec.Monitor.Prometheus.Port = api.PrometheusExporterPortNumber
		}
	}

	if memcached.Spec.DoNotPause {
		c.recorder.Eventf(
			memcached.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToPause,
			`Memcached "%v" is locked.`,
			memcached.Name,
		)

		if err := c.reCreateMemcached(memcached); err != nil {
			c.recorder.Eventf(
				memcached.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToCreate,
				`Failed to recreate Memcached: "%v". Reason: %v`,
				memcached.Name,
				err,
			)
			return err
		}
		return nil
	}

	if _, err := c.createDormantDatabase(memcached); err != nil {
		c.recorder.Eventf(
			memcached.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToCreate,
			`Failed to create DormantDatabase: "%v". Reason: %v`,
			memcached.Name,
			err,
		)
		return err
	}
	c.recorder.Eventf(
		memcached.ObjectReference(),
		core.EventTypeNormal,
		eventer.EventReasonSuccessfulCreate,
		`Successfully created DormantDatabase: "%v"`,
		memcached.Name,
	)

	if memcached.Spec.Monitor != nil {
		if ok, err := c.deleteMonitor(memcached); err != nil {
			c.recorder.Eventf(
				memcached.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				"Failed to delete monitoring system. Reason: %v",
				err,
			)
			log.Errorln(err)
			return nil
		} else if ok {
			c.recorder.Event(
				memcached.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulMonitorDelete,
				"Successfully deleted monitoring system.",
			)
		}
	}
	return nil
}

func (c *Controller) reCreateMemcached(memcached *api.Memcached) error {
	_memcached := &api.Memcached{
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
