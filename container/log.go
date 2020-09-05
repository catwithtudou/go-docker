package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go-docker/config"
	"io/ioutil"
	"os"
	"path"
)

/**
 *@Author tudou
 *@Date 2020/9/5
 **/


// 查看容器中的日志消息
func LookContainerLog(containerName string){
	logFileName:=path.Join(config.DefaultContainerInfoPath,containerName,config.ContainerLogFileName)
	file,err:=os.Open(logFileName)
	if err!=nil{
		logrus.Errorf("failed to open the log file path[%s];err: %v",logFileName,err)
		return
	}
	bs,err:= ioutil.ReadAll(file)
	if err!=nil{
		logrus.Errorf("failed to read log file;err: %v",err)
	}
	_,_ = fmt.Fprintf(os.Stdout,string(bs))
	return
}