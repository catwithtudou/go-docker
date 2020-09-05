package main

import (
	"github.com/sirupsen/logrus"
	"go-docker/cgroups"
	"go-docker/cgroups/subsystem"
	"go-docker/container"
	"os"
	"strings"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

//运行指令
func Run(cmdArray []string, tty bool, res *subsystem.ResourceLimitConfig,containerName, imageName, volume, net string, envs, ports []string) {
	//获取容器父进程
	//此进程是被资源限制后的进程
	parent, writePipe := container.NewParentProcess(tty,containerName,imageName,volume,envs)
	if parent == nil {
		logrus.Error("failed to new parent process")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("failed to start parent;err: %v", err)
		return
	}
	//获取容器ID
	containerID := container.GenContainerID(10)
	if containerName == "" {
		containerName = containerID
	}

	//记录容器信息
	err:= container.RecordContainerInfo(parent.Process.Pid,cmdArray,containerName,containerID)


	//添加资源限制
	cGroupManager := cgroups.NewCGroupManager("gocker")
	//删除资源限制
	defer cGroupManager.Destroy()
	//设置资源限制
	cGroupManager.Set(res)
	//将容器进程，加入到各个subsystem挂载对应的cGroup中
	cGroupManager.Apply(parent.Process.Pid)

	// 设置初始化命令
	err = sendInitCommand(cmdArray, writePipe)
	if err != nil {
		panic(err)
	}

	if tty {
		// 等待父进程结束
		err := parent.Wait()
		if err != nil {
			logrus.Errorf("parent wait, err: %v", err)
		}
		// 删除容器工作空间
		err = container.DeleteWorkSpace(containerName, volume)
		if err != nil {
			logrus.Errorf("delete work space, err: %v", err)
		}
		// 删除容器信息
		container.DeleteContainerInfo(containerName)
	}
}

//初始化容器命令
func sendInitCommand(cmdArray []string, writePipe *os.File) error {
	command := strings.Join(cmdArray, " ")
	logrus.Info("all command is " + command)
	_, err := writePipe.WriteString(command)
	if err != nil {
		panic(err)
		return err
	}
	return writePipe.Close()
}
