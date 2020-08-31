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

//
//var buildstamp = ""
//var githash = ""
//var goversion = ""

func main() {
	//args := os.Args
	//if len(args) == 2 && (args[1] == "--version" || args[1] == "-v") {
	//	fmt.Printf("Git Commit Hash: %s\n", githash)
	//	fmt.Printf("UTC Build Time : %s\n", buildstamp)
	//	fmt.Printf("Golang Version : %s\n", goversion)
	//	return
	//}

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
