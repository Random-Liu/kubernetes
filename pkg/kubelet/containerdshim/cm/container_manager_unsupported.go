// +build !linux

/*
Copyright 2016 The Kubernetes Authors.

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

package cm

import (
	"fmt"

	libcontainerd "k8s.io/kubernetes/pkg/kubelet/containerdtools"
	"k8s.io/kubernetes/pkg/kubelet/dockertools"
)

type unsupportedContainerManager struct {
}

func NewContainerManager(_ string, _ dockertools.DockerInterface, _ libcontainerd.Client) ContainerManager {
	return &unsupportedContainerManager{}
}

func (m *unsupportedContainerManager) Start() error {
	return fmt.Errorf("Container Manager is unsupported in this build")
}
