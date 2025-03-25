/*
Copyright 2020 The Tekton Authors

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

package v1beta1

import (
	v1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// ClusterTaskLister helps list ClusterTasks.
// All objects returned here must be treated as read-only.
type ClusterTaskLister interface {
	// List lists all ClusterTasks in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.ClusterTask, err error)
	// Get retrieves the ClusterTask from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.ClusterTask, error)
	ClusterTaskListerExpansion
}

// clusterTaskLister implements the ClusterTaskLister interface.
type clusterTaskLister struct {
	listers.ResourceIndexer[*v1beta1.ClusterTask]
}

// NewClusterTaskLister returns a new ClusterTaskLister.
func NewClusterTaskLister(indexer cache.Indexer) ClusterTaskLister {
	return &clusterTaskLister{listers.New[*v1beta1.ClusterTask](indexer, v1beta1.Resource("clustertask"))}
}
