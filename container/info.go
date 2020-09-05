package container

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go-docker/config"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

/**
 *@Author tudou
 *@Date 2020/9/5
 **/

var(
	letterBytes = "0123456789"
)

type ContainerInfo struct{
	Pid string `json:"pid"` //宿主机上的PID
	Id string `json:"id"` //容器ID
	Command string `json:"command"` //init进程内的运行命令
	Name string `json:"name"`
	CreateTime string `json:"creat_time"`
	Status string `json:"status"`
	Volume string `json:"volume"` //数据卷
	PartMapping []string`json:"part_mapping"` //端口映射
}

// 记录容器信息
func RecordContainerInfo(containerPID int,cmdArray []string,containerName,containerID string)error{
	info:= &ContainerInfo{
		Pid:         strconv.Itoa(containerPID),
		Id:          containerID,
		Command:     strings.Join(cmdArray,""),
		Name:        containerName,
		CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		Status:      config.Running,
	}

	dir:= path.Join(config.DefaultContainerInfoPath,containerName)
	_,err:=os.Stat(dir)
	if err!=nil && os.IsNotExist(err){
		err:=os.MkdirAll(dir,os.ModePerm)
		if err!=nil{
			logrus.Errorf("failed to mkdir container dir[%s];err: %v",dir,err)
			return err
		}
	}

	fileName:=fmt.Sprintf("%s/%s",dir,config.ContainerInfoFileName)
	file,err:=os.Create(fileName)
	if err!=nil{
		logrus.Errorf("failed to create config.json[%s];err: %s",fileName,err)
		return err
	}

	jsonInfo,_ := json.Marshal(info)
	_,err = file.WriteString(string(jsonInfo))
	if err!=nil{
		logrus.Errorf("failed to write config.json[%s];err: %v",fileName,err)
		return err
	}

	return nil
}

// 获取容器ID
func GenContainerID(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// 删除容器Info
func DeleteContainerInfo(containerName string) {
	dir := path.Join(config.DefaultContainerInfoPath, containerName)
	err := os.RemoveAll(dir)
	if err != nil {
		logrus.Errorf("failed to remove container info;err: %v", err)
	}
}

// 获取容器信息Info
func getContainerInfo(containerName string) (*ContainerInfo, error) {
	filePath := path.Join(config.DefaultContainerInfoPath, containerName, config.ContainerInfoFileName)
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("failed to read file path[%s];err: %v", filePath, err)
		return nil, err
	}
	info := &ContainerInfo{}
	err = json.Unmarshal(bs, info)
	return info, err
}


// 遍历容器信息
func ListContainerInfo() {
	files, err := ioutil.ReadDir(config.DefaultContainerInfoPath)
	if err != nil {
		logrus.Errorf("failed to read info dir;err: %v", err)
	}

	var infos []*ContainerInfo
	for _, file := range files {
		info, err := getContainerInfo(file.Name())
		if err != nil {
			logrus.Errorf("failed to get container info name[%s];err: %v", file.Name(), err)
			continue
		}
		infos = append(infos, info)
	}

	// 打印
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 2, ' ', 0)
	_, _ = fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, info := range infos {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", info.Id, info.Name, info.Pid, info.Status, info.Command, info.CreateTime)
	}

	// 刷新标准输出流缓存区，将容器列表打印出来
	if err := w.Flush(); err != nil {
		logrus.Errorf("failed to flush info;err: %v", err)
	}
}