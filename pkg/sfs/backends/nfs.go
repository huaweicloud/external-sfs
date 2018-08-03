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
	"fmt"
	"strings"

	"k8s.io/api/core/v1"
)

// NFS implements ShareBackend interface for k8s NFS
type NFS struct {
	Backend
}

// Name of the backend
func (NFS) Name() string { return "NFS" }

// BuildSource builds PersistentVolumeSource for k8s NFS
func (NFS) BuildSource(args *BuildSourceArgs) (*v1.PersistentVolumeSource, error) {
	delimPos := strings.LastIndexByte(args.Location, ':')
	if delimPos <= 0 {
		return &v1.PersistentVolumeSource{},
			fmt.Errorf("failed to parse address and location from export location '%s'", args.Location)
	}

	server := args.Location[:delimPos]
	path := args.Location[delimPos+1:]

	return &v1.PersistentVolumeSource{
		NFS: &v1.NFSVolumeSource{
			Server:   server,
			Path:     path,
			ReadOnly: false,
		},
	}, nil
}
