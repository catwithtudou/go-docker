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
func Run(cmdArray []string, tty bool, res *subsystem.ResourceLimitConfig) {
	//获取容器父进程
	//此进程是被资源限制后的进程
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Error("failed to new parent process")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("failed to start parent;err: %v", err)
		return
	}
	//添加资源限制
	cGroupManager := cgroups.NewCGroupManager("gocker")
	//删除资源限制
	defer cGroupManager.Destroy()
	//设置资源限制
	cGroupManager.Set(res)
	//将容器进程，加入到各个subsystem挂载对应的cGroup中
	cGroupManager.Apply(parent.Process.Pid)

	//
	err := sendInitCommand(cmdArray, writePipe)
	if err != nil {
		panic(err)
	}
	err = parent.Wait()
	if err != nil {
		panic(err)
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
