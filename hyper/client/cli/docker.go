/*
Copyright © 2022 Hypergiant, LLC

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
package cli

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/term"
)

type DockerClient struct {
	Cli client.Client
	Ctx context.Context
}

func NewDockerClient() *DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	dockerClient := &DockerClient{
		Cli: *cli,
		Ctx: context.Background(),
	}

	return dockerClient
}

func (dockerClient *DockerClient) CreateContainer(
	image, name string, contConfig *container.Config,
	hostConfig *container.HostConfig, pullImage bool,
) (string, error) {

	if pullImage {
		reader, err := dockerClient.Cli.ImagePull(dockerClient.Ctx, image, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}

		defer reader.Close()

		termFd, isTerm := term.GetFdInfo(os.Stderr)
		err = jsonmessage.DisplayJSONMessagesStream(reader, os.Stderr, termFd, isTerm, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

	containerCreatedBody, err := dockerClient.Cli.ContainerCreate(dockerClient.Ctx, contConfig, hostConfig, nil, nil, name)

	return containerCreatedBody.ID, err
}

func (dockerClient *DockerClient) ExecuteContainer(containerID string, attach bool) {
	if err := dockerClient.Cli.ContainerStart(dockerClient.Ctx, containerID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	if attach {
		statusCh, errCh := dockerClient.Cli.ContainerWait(dockerClient.Ctx, containerID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}

		out, err := dockerClient.Cli.ContainerLogs(dockerClient.Ctx, containerID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}

		_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		if err != nil {
			panic(err)
		}
	}
}

func (dockerClient *DockerClient) ListContainers(containerName string) ([]types.Container, error) {
	containerListOptions := types.ContainerListOptions{}
	if containerName != "" {
		containerListOptions.Filters = filters.NewArgs()
		containerListOptions.Filters.Add("name", containerName)
	}
	containers, err := dockerClient.Cli.ContainerList(dockerClient.Ctx, containerListOptions)

	if err != nil {
		panic(err)
	}

	return containers, err
}
func (dockerClient *DockerClient) ListAllRunningContainers() ([]types.Container, error) {
	return dockerClient.ListContainers("")
}

func (dockerClient *DockerClient) ListImages() ([]types.ImageSummary, error) {
	imageListOptions := types.ImageListOptions{}
	images, err := dockerClient.Cli.ImageList(dockerClient.Ctx, imageListOptions)

	if err != nil {
		panic(err)
	}

	return images, err
}

func (dockerClient *DockerClient) InspectContainer(containerId string) types.ContainerJSON {
	containerJSON, _, err := dockerClient.Cli.ContainerInspectWithRaw(dockerClient.Ctx, containerId, false)

	if err != nil {
		panic(err)
	}

	return containerJSON
}

func (dockerClient *DockerClient) RemoveContainer(containerId string) error {
	errStop := dockerClient.Cli.ContainerStop(dockerClient.Ctx, containerId, nil)

	if errStop != nil {
		panic(errStop)
	}

	errRemove := dockerClient.Cli.ContainerRemove(dockerClient.Ctx, containerId, types.ContainerRemoveOptions{})

	if errRemove != nil {
		panic(errRemove)
	}

	return errRemove
}

type HyperPackageDockerfileParameters struct {
	StudyPath string
}

func (dockerClient *DockerClient) CreateDockerFile(studyPath string, savePath string, requirements bool) {
	dockerFileTemplate := ""
	if requirements {
		dockerFileTemplate = `
FROM ghcr.io/gohypergiant/hyperdrive-jupyter:cpu-localstable
COPY requirements.txt /home/jovyan/requirements.txt
RUN pip install -r requirements.txt
`
	} else {
		fastAppApiKey := generateFastAppAPIKey()
		dockerFileTemplate = fmt.Sprintf(`
FROM ubuntu:latest as builder
RUN apt update -y && apt install unzip -y
ADD {{.StudyPath}} study.hyperpackage.zip
RUN unzip ./study.hyperpackage.zip -d /hyperpackage

FROM ghcr.io/gohypergiant/gohypergiant/mlsdk-fast-app:stable
COPY --from=builder /hyperpackage /hyperpackage
ENV fastkey=%s
RUN echo "*** The Fast App API key is: $fastkey ***"
`, fastAppApiKey)
	}

	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	tmpl, err := template.New("dockerfile").Parse(dockerFileTemplate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	params := HyperPackageDockerfileParameters{StudyPath: studyPath}
	err = tmpl.Execute(file, params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = file.Sync()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func (dockerClient *DockerClient) BuildImage(dockerfilePath string, tags []string) {

	dockerBuildContext, err := archive.Tar("./", archive.Uncompressed) // TODO: pass this path in as a flag
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	opts := types.ImageBuildOptions{
		Dockerfile: dockerfilePath,
		Tags:       tags,
		Remove:     true,
	}
	res, err := dockerClient.Cli.ImageBuild(dockerClient.Ctx, dockerBuildContext, opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	err = jsonmessage.DisplayJSONMessagesStream(res.Body, os.Stderr, termFd, isTerm, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
