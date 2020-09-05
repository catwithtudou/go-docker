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
		cli.StringFlag{
			Name:  "net",
			Usage: "container network",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "docker volume",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "docker env",
		},
		cli.StringSliceFlag{
			Name:  "p",
			Usage: "port mapping",
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

		detach := context.Bool("d")

		if tty && detach {
			return fmt.Errorf("ti and d paramter can not both provided")
		}

		containerName := context.String("name")
		volume := context.String("v")
		net := context.String("net")
		// 要运行的镜像名
		imageName := context.Args().Get(0)
		envs := context.StringSlice("e")
		ports := context.StringSlice("p")

		//启动命令
		Run(cmdArray, tty, res,containerName,imageName,volume,net,envs,ports)
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

var ListCommand = cli.Command{
	Name: "ps",
	Usage: "list all container",
	Action: func(ctx *cli.Context)error{
		container.ListContainerInfo()
		return nil
	},
}

var StopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop the container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing stop container name")
		}
		containerName := ctx.Args().Get(0)
		container.StopContainer(containerName)
		return nil
	},
}
