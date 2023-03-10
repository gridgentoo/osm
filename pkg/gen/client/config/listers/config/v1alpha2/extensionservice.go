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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha2

import (
	v1alpha2 "github.com/openservicemesh/osm/pkg/apis/config/v1alpha2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ExtensionServiceLister helps list ExtensionServices.
// All objects returned here must be treated as read-only.
type ExtensionServiceLister interface {
	// List lists all ExtensionServices in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha2.ExtensionService, err error)
	// ExtensionServices returns an object that can list and get ExtensionServices.
	ExtensionServices(namespace string) ExtensionServiceNamespaceLister
	ExtensionServiceListerExpansion
}

// extensionServiceLister implements the ExtensionServiceLister interface.
type extensionServiceLister struct {
	indexer cache.Indexer
}

// NewExtensionServiceLister returns a new ExtensionServiceLister.
func NewExtensionServiceLister(indexer cache.Indexer) ExtensionServiceLister {
	return &extensionServiceLister{indexer: indexer}
}

// List lists all ExtensionServices in the indexer.
func (s *extensionServiceLister) List(selector labels.Selector) (ret []*v1alpha2.ExtensionService, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.ExtensionService))
	})
	return ret, err
}

// ExtensionServices returns an object that can list and get ExtensionServices.
func (s *extensionServiceLister) ExtensionServices(namespace string) ExtensionServiceNamespaceLister {
	return extensionServiceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ExtensionServiceNamespaceLister helps list and get ExtensionServices.
// All objects returned here must be treated as read-only.
type ExtensionServiceNamespaceLister interface {
	// List lists all ExtensionServices in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha2.ExtensionService, err error)
	// Get retrieves the ExtensionService from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha2.ExtensionService, error)
	ExtensionServiceNamespaceListerExpansion
}

// extensionServiceNamespaceLister implements the ExtensionServiceNamespaceLister
// interface.
type extensionServiceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ExtensionServices in the indexer for a given namespace.
func (s extensionServiceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha2.ExtensionService, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.ExtensionService))
	})
	return ret, err
}

// Get retrieves the ExtensionService from the indexer for a given namespace and name.
func (s extensionServiceNamespaceLister) Get(name string) (*v1alpha2.ExtensionService, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha2.Resource("extensionservice"), name)
	}
	return obj.(*v1alpha2.ExtensionService), nil
}
