package container

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go-docker/config"
	"io/ioutil"
	"path"
	"strconv"
	"syscall"
)

/**
 *@Author tudou
 *@Date 2020/9/5
 **/

// 停止容器，修改容器状态
func StopContainer(containerName string){
	info,err:=getContainerInfo(containerName)
	if err !=nil{
		logrus.Errorf("failed to get container info;err: %v",err)
		return
	}
	if info.Pid!=""{
		pid,_ := strconv.Atoi(info.Pid)
		//杀死进程
		if err:= syscall.Kill(pid,syscall.SIGTERM);err!=nil{
			logrus.Errorf("failed to stop container pid[%s];err: %v",info.Pid,err)
			return
		}

		//修改容器状态
		info.Status = config.Stop
		info.Pid = ""
		jsonInfo,_ :=json.Marshal(info)
		fileName:=path.Join(config.DefaultContainerInfoPath,containerName,config.ContainerInfoFileName)
		err:=ioutil.WriteFile(fileName,jsonInfo,0622)
		if err!=nil{
			logrus.Errorf("failed to write container config.json;err: %v",err)
			return
		}
	}

}
