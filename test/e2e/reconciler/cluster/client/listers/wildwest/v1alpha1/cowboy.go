/*
Copyright 2021 The KCP Authors.

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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	v1alpha1 "github.com/kcp-dev/kcp/test/e2e/reconciler/cluster/apis/wildwest/v1alpha1"
)

// CowboyLister helps list Cowboys.
// All objects returned here must be treated as read-only.
type CowboyLister interface {
	// List lists all Cowboys in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Cowboy, err error)
	// Cowboys returns an object that can list and get Cowboys.
	Cowboys(namespace string) CowboyNamespaceLister
	CowboyListerExpansion
}

// cowboyLister implements the CowboyLister interface.
type cowboyLister struct {
	indexer cache.Indexer
}

// NewCowboyLister returns a new CowboyLister.
func NewCowboyLister(indexer cache.Indexer) CowboyLister {
	return &cowboyLister{indexer: indexer}
}

// List lists all Cowboys in the indexer.
func (s *cowboyLister) List(selector labels.Selector) (ret []*v1alpha1.Cowboy, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Cowboy))
	})
	return ret, err
}

// Cowboys returns an object that can list and get Cowboys.
func (s *cowboyLister) Cowboys(namespace string) CowboyNamespaceLister {
	return cowboyNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CowboyNamespaceLister helps list and get Cowboys.
// All objects returned here must be treated as read-only.
type CowboyNamespaceLister interface {
	// List lists all Cowboys in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Cowboy, err error)
	// Get retrieves the Cowboy from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Cowboy, error)
	CowboyNamespaceListerExpansion
}

// cowboyNamespaceLister implements the CowboyNamespaceLister
// interface.
type cowboyNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Cowboys in the indexer for a given namespace.
func (s cowboyNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Cowboy, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Cowboy))
	})
	return ret, err
}

// Get retrieves the Cowboy from the indexer for a given namespace and name.
func (s cowboyNamespaceLister) Get(name string) (*v1alpha1.Cowboy, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("cowboy"), name)
	}
	return obj.(*v1alpha1.Cowboy), nil
}