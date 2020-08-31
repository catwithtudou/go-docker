package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

func main() {
	//这里需要导入的是exec_linux.go(即环境变为linux)
	cmd := exec.Command("sh")
	//Linux 对线程提供了六种隔离机制，分别为：uts pid user mount network ipc
	//uts: 用来隔离主机名
	//pid：用来隔离进程 PID 号的
	//user: 用来隔离用户的
	//mount：用来隔离各个进程看到的挂载点视图
	//network: 用来隔离网络
	//ipc：用来隔离 System V IPC 和 POSIX message queues
	//docker会分别进行隔离
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNET,
		UidMappings: []syscall.SysProcIDMap{
			{
				//容器UID
				ContainerID: 1,
				//宿主机UID
				HostID: 0,
				Size:   1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				//容器GID
				ContainerID: 1,
				//宿主机GID
				HostID: 0,
				Size:   1,
			},
		},
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
