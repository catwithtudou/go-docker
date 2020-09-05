package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go-docker/cgroups/subsystem"
	"go-docker/container"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

//创建namespace隔离的容器进程(启动容器)
var RunCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cGroups limit",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti", //表示是否前台运行
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory list",
		},
		cli.StringFlag{
			Name:  "ch",
			Usage: "cpu share limit",
		},
		cli.StringFlag{
			Name:  "cs",
			Usage: "cpu set limit",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container args")
		}
		//获取ti参数
		tty := context.Bool("ti")

		//获取subSystem资源限制配置
		res := &subsystem.ResourceLimitConfig{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cs"),
			CpuShare:    context.String("ch"),
		}

		//cmdArray是容器运行后执行的第一个命令信息
		//cmdArray[0]:command
		//cmdArray[1:]:args
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		//启动命令
		Run(cmdArray, tty, res)
		return nil
	},
}

//初始化容器内容(mount namespace挂载proc文件系统)，运行用户程序
//mount namespace 用来隔离文件系统的挂载点，这样进程就只能看到自己的 mount namespace 中的文件系统挂载点
var InitCommand = cli.Command{
	Name:  "init",
	Usage: "Init a container process run user's process in container. Don't call it outside",
	Action: func(context *cli.Context) error {
		logrus.Info("init beginning....")
		return container.RunContainerInitProcess()
	},
}

var LogCommand = cli.Command{
	Name: "logs",
	Usage: "looking container log",
	Action: func(ctx *cli.Context)error{
		if len(ctx.Args())<1{
			return fmt.Errorf("missing container name")
		}
		containerName := ctx.Args().Get(0)
		container.LookContainerLog(containerName)
		return nil
	},
}