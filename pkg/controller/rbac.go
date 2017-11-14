package controller

import (
	kutilcore "github.com/appscode/kutil/core/v1"
	kutilrbac "github.com/appscode/kutil/rbac/v1beta1"
	"github.com/k8sdb/apimachinery/apis/kubedb"
	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Controller) deleteRole(memcached *tapi.Memcached) error {
	// Delete existing Roles
	if err := c.Client.RbacV1beta1().Roles(memcached.Namespace).Delete(memcached.OffshootName(), nil); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (c *Controller) createRole(memcached *tapi.Memcached) error {
	// Create new Roles
	_, err := kutilrbac.CreateOrPatchRole(
		c.Client,
		metav1.ObjectMeta{
			Name:      memcached.OffshootName(),
			Namespace: memcached.Namespace,
		},
		func(in *rbac.Role) *rbac.Role {
			in.Rules = []rbac.PolicyRule{
				{
					APIGroups:     []string{kubedb.GroupName},
					Resources:     []string{tapi.ResourceTypeMemcached},
					ResourceNames: []string{memcached.Name},
					Verbs:         []string{"get"},
				},
			}
			return in
		},
	)
	return err
}

func (c *Controller) deleteServiceAccount(memcached *tapi.Memcached) error {
	// Delete existing ServiceAccount
	if err := c.Client.CoreV1().ServiceAccounts(memcached.Namespace).Delete(memcached.OffshootName(), nil); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (c *Controller) createServiceAccount(memcached *tapi.Memcached) error {
	// Create new ServiceAccount
	_, err := kutilcore.CreateOrPatchServiceAccount(
		c.Client,
		metav1.ObjectMeta{
			Name:      memcached.OffshootName(),
			Namespace: memcached.Namespace,
		},
		func(in *core.ServiceAccount) *core.ServiceAccount {
			return in
		},
	)
	return err
}

func (c *Controller) deleteRoleBinding(memcached *tapi.Memcached) error {
	// Delete existing RoleBindings
	if err := c.Client.RbacV1beta1().RoleBindings(memcached.Namespace).Delete(memcached.OffshootName(), nil); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (c *Controller) createRoleBinding(memcached *tapi.Memcached) error {
	// Ensure new RoleBindings
	_, err := kutilrbac.CreateOrPatchRoleBinding(
		c.Client,
		metav1.ObjectMeta{
			Name:      memcached.OffshootName(),
			Namespace: memcached.Namespace,
		},
		func(in *rbac.RoleBinding) *rbac.RoleBinding {
			in.RoleRef = rbac.RoleRef{
				APIGroup: rbac.GroupName,
				Kind:     "Role",
				Name:     memcached.OffshootName(),
			}
			in.Subjects = []rbac.Subject{
				{
					Kind:      rbac.ServiceAccountKind,
					Name:      memcached.OffshootName(),
					Namespace: memcached.Namespace,
				},
			}
			return in
		},
	)
	return err
}

func (c *Controller) createRBACStuff(memcached *tapi.Memcached) error {
	// Delete Existing Role
	if err := c.deleteRole(memcached); err != nil {
		return err
	}
	// Create New Role
	if err := c.createRole(memcached); err != nil {
		return err
	}

	// Create New ServiceAccount
	if err := c.createServiceAccount(memcached); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
	}

	// Create New RoleBinding
	if err := c.createRoleBinding(memcached); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

func (c *Controller) deleteRBACStuff(memcached *tapi.Memcached) error {
	// Delete Existing Role
	if err := c.deleteRole(memcached); err != nil {
		return err
	}

	// Delete ServiceAccount
	if err := c.deleteServiceAccount(memcached); err != nil {
		return err
	}

	// Delete New RoleBinding
	if err := c.deleteRoleBinding(memcached); err != nil {
		return err
	}

	return nil
}
