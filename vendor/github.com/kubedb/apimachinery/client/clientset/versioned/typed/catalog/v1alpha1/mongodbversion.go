/*
Copyright 2019 The KubeDB Authors.

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

package v1alpha1

import (
	v1alpha1 "github.com/kubedb/apimachinery/apis/catalog/v1alpha1"
	scheme "github.com/kubedb/apimachinery/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MongoDBVersionsGetter has a method to return a MongoDBVersionInterface.
// A group's client should implement this interface.
type MongoDBVersionsGetter interface {
	MongoDBVersions() MongoDBVersionInterface
}

// MongoDBVersionInterface has methods to work with MongoDBVersion resources.
type MongoDBVersionInterface interface {
	Create(*v1alpha1.MongoDBVersion) (*v1alpha1.MongoDBVersion, error)
	Update(*v1alpha1.MongoDBVersion) (*v1alpha1.MongoDBVersion, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.MongoDBVersion, error)
	List(opts v1.ListOptions) (*v1alpha1.MongoDBVersionList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MongoDBVersion, err error)
	MongoDBVersionExpansion
}

// mongoDBVersions implements MongoDBVersionInterface
type mongoDBVersions struct {
	client rest.Interface
}

// newMongoDBVersions returns a MongoDBVersions
func newMongoDBVersions(c *CatalogV1alpha1Client) *mongoDBVersions {
	return &mongoDBVersions{
		client: c.RESTClient(),
	}
}

// Get takes name of the mongoDBVersion, and returns the corresponding mongoDBVersion object, and an error if there is any.
func (c *mongoDBVersions) Get(name string, options v1.GetOptions) (result *v1alpha1.MongoDBVersion, err error) {
	result = &v1alpha1.MongoDBVersion{}
	err = c.client.Get().
		Resource("mongodbversions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MongoDBVersions that match those selectors.
func (c *mongoDBVersions) List(opts v1.ListOptions) (result *v1alpha1.MongoDBVersionList, err error) {
	result = &v1alpha1.MongoDBVersionList{}
	err = c.client.Get().
		Resource("mongodbversions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mongoDBVersions.
func (c *mongoDBVersions) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("mongodbversions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a mongoDBVersion and creates it.  Returns the server's representation of the mongoDBVersion, and an error, if there is any.
func (c *mongoDBVersions) Create(mongoDBVersion *v1alpha1.MongoDBVersion) (result *v1alpha1.MongoDBVersion, err error) {
	result = &v1alpha1.MongoDBVersion{}
	err = c.client.Post().
		Resource("mongodbversions").
		Body(mongoDBVersion).
		Do().
		Into(result)
	return
}

// Update takes the representation of a mongoDBVersion and updates it. Returns the server's representation of the mongoDBVersion, and an error, if there is any.
func (c *mongoDBVersions) Update(mongoDBVersion *v1alpha1.MongoDBVersion) (result *v1alpha1.MongoDBVersion, err error) {
	result = &v1alpha1.MongoDBVersion{}
	err = c.client.Put().
		Resource("mongodbversions").
		Name(mongoDBVersion.Name).
		Body(mongoDBVersion).
		Do().
		Into(result)
	return
}

// Delete takes name of the mongoDBVersion and deletes it. Returns an error if one occurs.
func (c *mongoDBVersions) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("mongodbversions").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mongoDBVersions) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("mongodbversions").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched mongoDBVersion.
func (c *mongoDBVersions) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MongoDBVersion, err error) {
	result = &v1alpha1.MongoDBVersion{}
	err = c.client.Patch(pt).
		Resource("mongodbversions").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
