/*
Copyright 2018 The Kubernetes Authors.

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

package backends

import (
	"k8s.io/api/core/v1"
)

// Backend builds PersistentVolumeSource
type Backend interface {
	// Name of the share backend
	Name() string

	// Called during share provision, the result is used in the final PersistentVolume object.
	BuildSource(*BuildSourceArgs) (*v1.PersistentVolumeSource, error)
}

// BuildSourceArgs contains arguments
type BuildSourceArgs struct {
	Location string
}
