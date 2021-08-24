package yum

import (
	"fmt"
	"github.com/linuxsuren/http-downloader/pkg/exec"
	"runtime"
	"strings"
	"time"
)

// DockerInstallerInCentOS is the installer of Docker in CentOS
type DockerInstallerInCentOS struct {
	count int
}

// Available check if support current platform
func (d *DockerInstallerInCentOS) Available() (ok bool) {
	if runtime.GOOS == "linux" {
		_, err := exec.LookPath("yum")
		ok = err == nil
	}
	return
}

// Install installs the Docker
func (d *DockerInstallerInCentOS) Install() (err error) {
	if err = exec.RunCommand("yum", "install", "-y",
		"yum-utils"); err != nil {
		return
	}
	if err = exec.RunCommand("yum-config-manager", "--add-repo",
		"https://download.docker.com/linux/centos/docker-ce.repo"); err != nil {
		return
	}
	if err = exec.RunCommand("yum", "install", "-y",
		"docker-ce",
		"docker-ce-cli",
		"containerd.io"); err != nil {
		return
	}
	return
}

// Uninstall uninstalls the Docker
func (d *DockerInstallerInCentOS) Uninstall() (err error) {
	err = exec.RunCommand("yum", "remove", "-y",
		"docker",
		"docker-client",
		"docker-client-latest",
		"docker-common",
		"docker-latest",
		"docker-latest-logrotate",
		"docker-logrotate",
		"docker-engine",
		"docker-ce",
		"docker-ce-cli",
		"containerd.io")
	return
}

// WaitForStart waits for the service be started
func (d *DockerInstallerInCentOS) WaitForStart() (ok bool, err error) {
	var result string
	if result, err = exec.RunCommandAndReturn("systemctl", "", "status", "docker"); err != nil {
		return
	} else if strings.Contains(result, "Unit docker.service could not be found") {
		err = fmt.Errorf("unit docker.service could not be found")
	} else if strings.Contains(result, "Active: active") {
		ok = true
	} else {
		if d.count > 0 {
			fmt.Println("waiting for Docker service start")
		} else if d.count > 9 {
			return
		}

		d.count ++
		time.Sleep(time.Second * 1)
		return d.WaitForStart()
	}
	return
}

// Start starts the Docker service
func (d *DockerInstallerInCentOS) Start() error {
	return exec.RunCommand("systemctl", "start", "docker")
}

// Stop stops the Docker service
func (d *DockerInstallerInCentOS) Stop() error {
	return exec.RunCommand("systemctl", "stop", "docker")
}