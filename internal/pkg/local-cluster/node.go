/*
Copyright 2020 Rafael Fernández López <ereslibre@ereslibre.es>

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

package localcluster

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

type Node struct {
	Name    string
	Cluster *Cluster
}

func (node *Node) Create() error {
	if err := node.createRuntimeDirectory(); err != nil {
		return err
	}
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	return exec.Command(
		"docker", "run", "-d", "--privileged",
		"--name", fmt.Sprintf("%s-%s", node.Cluster.Name, node.Name),
		"-v", fmt.Sprintf("%s:/containerd-socket", node.runtimeDirectory()),
		"-e", fmt.Sprintf("CONTAINERD_SOCK_UID=%s", currentUser.Uid),
		"-e", fmt.Sprintf("CONTAINERD_SOCK_GID=%s", currentUser.Gid),
		"oneinfra/containerd:latest",
	).Run()
}

func (node *Node) Destroy() error {
	exec.Command(
		"docker", "rm", "-f", fmt.Sprintf("%s-%s", node.Cluster.Name, node.Name),
	).Run()
	return os.RemoveAll(node.runtimeDirectory())
}

func (node *Node) createRuntimeDirectory() error {
	return os.MkdirAll(node.runtimeDirectory(), 0755)
}

func (node *Node) runtimeDirectory() string {
	return filepath.Join(node.Cluster.directory(), node.Name)
}
