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

package sfs

import (
	"fmt"

	"github.com/huaweicloud/external-sfs/pkg/sfs/backends"
)

var (
	caches map[string]backends.Backend
)

// InitBackends for share
func InitBackends() {
	caches = make(map[string]backends.Backend)
	RegisterBackend(&backends.NFSBackend{})
}

// RegisterBackend for share
func RegisterBackend(b backends.Backend) {
	caches[b.Name()] = b
}

// GetBackend by name
func GetBackend(name string) (backends.Backend, error) {
	b, ok := caches[name]
	if !ok {
		return nil, fmt.Errorf("backend %s not found", name)
	}
	return b, nil
}
