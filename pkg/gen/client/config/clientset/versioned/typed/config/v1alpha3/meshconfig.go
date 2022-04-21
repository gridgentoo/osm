/*
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

package v1alpha3

import (
	"context"
	"time"

	v1alpha3 "github.com/openservicemesh/osm/pkg/apis/config/v1alpha3"
	scheme "github.com/openservicemesh/osm/pkg/gen/client/config/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MeshConfigsGetter has a method to return a MeshConfigInterface.
// A group's client should implement this interface.
type MeshConfigsGetter interface {
	MeshConfigs(namespace string) MeshConfigInterface
}

// MeshConfigInterface has methods to work with MeshConfig resources.
type MeshConfigInterface interface {
	Create(ctx context.Context, meshConfig *v1alpha3.MeshConfig, opts v1.CreateOptions) (*v1alpha3.MeshConfig, error)
	Update(ctx context.Context, meshConfig *v1alpha3.MeshConfig, opts v1.UpdateOptions) (*v1alpha3.MeshConfig, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha3.MeshConfig, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha3.MeshConfigList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha3.MeshConfig, err error)
	MeshConfigExpansion
}

// meshConfigs implements MeshConfigInterface
type meshConfigs struct {
	client rest.Interface
	ns     string
}

// newMeshConfigs returns a MeshConfigs
func newMeshConfigs(c *ConfigV1alpha3Client, namespace string) *meshConfigs {
	return &meshConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the meshConfig, and returns the corresponding meshConfig object, and an error if there is any.
func (c *meshConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha3.MeshConfig, err error) {
	result = &v1alpha3.MeshConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("meshconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MeshConfigs that match those selectors.
func (c *meshConfigs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha3.MeshConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha3.MeshConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("meshconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested meshConfigs.
func (c *meshConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("meshconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a meshConfig and creates it.  Returns the server's representation of the meshConfig, and an error, if there is any.
func (c *meshConfigs) Create(ctx context.Context, meshConfig *v1alpha3.MeshConfig, opts v1.CreateOptions) (result *v1alpha3.MeshConfig, err error) {
	result = &v1alpha3.MeshConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("meshconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(meshConfig).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a meshConfig and updates it. Returns the server's representation of the meshConfig, and an error, if there is any.
func (c *meshConfigs) Update(ctx context.Context, meshConfig *v1alpha3.MeshConfig, opts v1.UpdateOptions) (result *v1alpha3.MeshConfig, err error) {
	result = &v1alpha3.MeshConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("meshconfigs").
		Name(meshConfig.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(meshConfig).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the meshConfig and deletes it. Returns an error if one occurs.
func (c *meshConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("meshconfigs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *meshConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("meshconfigs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched meshConfig.
func (c *meshConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha3.MeshConfig, err error) {
	result = &v1alpha3.MeshConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("meshconfigs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
