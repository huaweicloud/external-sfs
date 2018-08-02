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
	"strconv"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// getStorageSizeInGiga from pvc
func getStorageSizeInGiga(pvc *v1.PersistentVolumeClaim) (int, error) {
	errStorageSizeNotConfigured := fmt.Errorf("requested storage capacity must be set")

	if pvc.Spec.Resources.Requests == nil {
		return 0, errStorageSizeNotConfigured
	}

	storageSize, ok := pvc.Spec.Resources.Requests[v1.ResourceStorage]
	if !ok {
		return 0, errStorageSizeNotConfigured
	}

	if storageSize.IsZero() {
		return 0, fmt.Errorf("requested storage size must not have zero value")
	}

	if storageSize.Sign() == -1 {
		return 0, fmt.Errorf("requested storage size must be greater than zero")
	}

	var buf []byte
	canonicalValue, _ := storageSize.AsScale(resource.Giga)
	storageSizeAsByteSlice, _ := canonicalValue.AsCanonicalBytes(buf)

	return strconv.Atoi(string(storageSizeAsByteSlice))
}
