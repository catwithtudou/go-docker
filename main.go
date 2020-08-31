package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

const (
	usage = "gocker"
	name  = "gocker"
)

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage

	//定义运行命令
	app.Commands = []cli.Command{RunCommand, InitCommand}

	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}

}
