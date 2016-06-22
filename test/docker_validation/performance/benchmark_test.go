/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package docker_performance

import (
	//"sync"

	"k8s.io/kubernetes/pkg/kubelet/dockertools"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/test/e2e/framework"

	docker "github.com/docker/engine-api/client"
	dockertypes "github.com/docker/engine-api/types"
	dockercontainer "github.com/docker/engine-api/types/container"
	dockerstrslice "github.com/docker/engine-api/types/strslice"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Docker configuration
var (
	endpoint = "unix:///var/run/docker.sock"
)

const (
	image = "busybox:1.24"

	totalCount   = 20 // 800
	runningCount = 1  // 200
	listCount    = 10
	inspectCount = 50
	sampleCount  = 10
)

/*type Task func(id int) error

//TODO merge this with comparable utilities in integration/framework/master_utils.go
// RunParallel spawns a goroutine per task in the given queue
func RunParallel(task Task, count int) {
	var wg sync.WaitGroup
	semCh := make(chan struct{}, count)
	wg.Add(count)
	for id := 0; id < count; id++ {
		go func(id int) {
			semCh <- struct{}{}
			task(id)
			<-semCh
			wg.Done()
		}(id)
	}
	wg.Wait()
	close(semCh)
}
*/

func newContainerName() string {
	return "benchmark_container_" + string(util.NewUUID())
}

var _ = Describe("Docker performance Test", func() {
	var client dockertools.DockerInterface

	BeforeEach(func() {
		c, err := docker.NewClient(endpoint, "", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		client = dockertools.NewKubeDockerClient(c)
		Expect(client.PullImage(image, dockertypes.AuthConfig{}, dockertypes.ImagePullOptions{})).To(Succeed())
	})

	Context("when benchmark container operations [Performance]", func() {
		Measure("container operations", func(b Benchmarker) {
			var all, running []string
			for i := 0; i < totalCount; i++ {
				b.Time("create container latency", func() {
					dockerOpts := dockertypes.ContainerCreateConfig{
						Name: newContainerName(),
						Config: &dockercontainer.Config{
							Cmd:   dockerstrslice.StrSlice([]string{"/bin/sh"}),
							Image: image,
						},
					}
					container, err := client.CreateContainer(dockerOpts)
					if err != nil {
						framework.Logf("Create container error: %v", err)
					}
					all = append(all, container.ID)
				})
			}

			Expect(runningCount).To(BeNumerically("<=", totalCount))
			for i := 0; i < runningCount; i++ {
				b.Time("start container latency", func() {
					id := all[i]
					err := client.StartContainer(id)
					if err != nil {
						framework.Logf("Start container error: %v", err)
					}
					running = append(running, id)
				})
			}

			for i := 0; i < sampleCount; i++ {
				b.Time("list all latency", func() {
					_, err := client.ListContainers(dockertypes.ContainerListOptions{All: true})
					if err != nil {
						framework.Logf("List containers error: %v", err)
					}
				})
			}

			for i := 0; i < sampleCount; i++ {
				b.Time("list alive latency", func() {
					_, err := client.ListContainers(dockertypes.ContainerListOptions{All: false})
					if err != nil {
						framework.Logf("List containers error: %v", err)
					}
				})
			}

			for _, id := range running {
				b.Time("stop container latnecy", func() {
					err := client.StopContainer(id, 0)
					if err != nil {
						framework.Logf("Stop container error: %v", err)
					}
				})
			}

			for _, id := range all {
				b.Time("remove container latency", func() {
					err := client.RemoveContainer(id, dockertypes.ContainerRemoveOptions{})
					if err != nil {
						framework.Logf("Remove container error: %v", err)
					}
				})
			}

		}, 3)

		/*


			Measure("start containers in parallel", func(b Benchmarker) {
				b.Time("start in parallel", func() {
					RunParallel(func(id int) error {
						client.StartContainer(ids[id], &docker.HostConfig{})
						return nil
					}, containerCount)
				})
				helpers.StopContainers(client, ids)
			}, concurrencyCount)

			Measure("list containers", func(b Benchmarker) {
				listAll := true
				b.Time("list", func() {
					for id := 0; id < containerCount; id++ {
						client.ListContainers(docker.ListContainersOptions{All: listAll})
					}
				})
			}, concurrencyCount)

			Measure("list containers in parallel", func(b Benchmarker) {
				listAll := true
				b.Time("list in parallel", func() {
					RunParallel(func(id int) error {
						client.ListContainers(docker.ListContainersOptions{All: listAll})
						return nil
					}, containerCount)
				})
			}, concurrencyCount)*/

		/*AfterEach(func() {
			for _, id := range ids {
				client.RemoveContainer(id, dockertypes.ContainerRemoveOptions{})
			}
		})*/
	})
})
