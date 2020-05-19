/*
Copyright The KubeDB Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEtcdOpsRequests implements EtcdOpsRequestInterface
type FakeEtcdOpsRequests struct {
	Fake *FakeOpsV1alpha1
	ns   string
}

var etcdopsrequestsResource = schema.GroupVersionResource{Group: "ops.kubedb.com", Version: "v1alpha1", Resource: "etcdopsrequests"}

var etcdopsrequestsKind = schema.GroupVersionKind{Group: "ops.kubedb.com", Version: "v1alpha1", Kind: "EtcdOpsRequest"}

// Get takes name of the etcdOpsRequest, and returns the corresponding etcdOpsRequest object, and an error if there is any.
func (c *FakeEtcdOpsRequests) Get(name string, options v1.GetOptions) (result *v1alpha1.EtcdOpsRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(etcdopsrequestsResource, c.ns, name), &v1alpha1.EtcdOpsRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdOpsRequest), err
}

// List takes label and field selectors, and returns the list of EtcdOpsRequests that match those selectors.
func (c *FakeEtcdOpsRequests) List(opts v1.ListOptions) (result *v1alpha1.EtcdOpsRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(etcdopsrequestsResource, etcdopsrequestsKind, c.ns, opts), &v1alpha1.EtcdOpsRequestList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.EtcdOpsRequestList{ListMeta: obj.(*v1alpha1.EtcdOpsRequestList).ListMeta}
	for _, item := range obj.(*v1alpha1.EtcdOpsRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested etcdOpsRequests.
func (c *FakeEtcdOpsRequests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(etcdopsrequestsResource, c.ns, opts))

}

// Create takes the representation of a etcdOpsRequest and creates it.  Returns the server's representation of the etcdOpsRequest, and an error, if there is any.
func (c *FakeEtcdOpsRequests) Create(etcdOpsRequest *v1alpha1.EtcdOpsRequest) (result *v1alpha1.EtcdOpsRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(etcdopsrequestsResource, c.ns, etcdOpsRequest), &v1alpha1.EtcdOpsRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdOpsRequest), err
}

// Update takes the representation of a etcdOpsRequest and updates it. Returns the server's representation of the etcdOpsRequest, and an error, if there is any.
func (c *FakeEtcdOpsRequests) Update(etcdOpsRequest *v1alpha1.EtcdOpsRequest) (result *v1alpha1.EtcdOpsRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(etcdopsrequestsResource, c.ns, etcdOpsRequest), &v1alpha1.EtcdOpsRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdOpsRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeEtcdOpsRequests) UpdateStatus(etcdOpsRequest *v1alpha1.EtcdOpsRequest) (*v1alpha1.EtcdOpsRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(etcdopsrequestsResource, "status", c.ns, etcdOpsRequest), &v1alpha1.EtcdOpsRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdOpsRequest), err
}

// Delete takes name of the etcdOpsRequest and deletes it. Returns an error if one occurs.
func (c *FakeEtcdOpsRequests) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(etcdopsrequestsResource, c.ns, name), &v1alpha1.EtcdOpsRequest{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEtcdOpsRequests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(etcdopsrequestsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.EtcdOpsRequestList{})
	return err
}

// Patch applies the patch and returns the patched etcdOpsRequest.
func (c *FakeEtcdOpsRequests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.EtcdOpsRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(etcdopsrequestsResource, c.ns, name, pt, data, subresources...), &v1alpha1.EtcdOpsRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdOpsRequest), err
}
