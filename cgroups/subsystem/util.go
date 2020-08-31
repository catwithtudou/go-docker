package subsystem

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

const (
	procSelfMountInfo = "/proc/self/mountinfo"
)

//获取cGroup的绝对路径
func GetCGroupPath(subsystem string, cGroupPath string, autoCreate bool) (resultPath string, err error) {
	rootPath, err := findCGroupMountPoint(subsystem)
	if err != nil {
		logrus.Errorf("failed to find cGroup path;err:%v", err)
		return "", err
	}
	resultPath = path.Join(rootPath, cGroupPath)
	_, err = os.Stat(resultPath)
	//判断是否存在
	if err != nil && os.IsNotExist(err) {
		//若不存在则创建文件夹
		if autoCreate {
			err = os.MkdirAll(resultPath, 0755)
		} else {
			logrus.Errorf("failed to stat the path[%s];err:%v", resultPath, err)
		}
		return resultPath, err
	}
	return resultPath, err
}

//查询mount被挂载proc的subsystem中cGroup根节点的目录
func findCGroupMountPoint(subsystem string) (path string, err error) {
	file, err := os.Open(procSelfMountInfo)
	if err != nil {
		return
	}
	defer file.Close()

	//TODO:/proc/self/mountinfo信息
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, option := range strings.Split(fields[len(fields)-1], ",") {
			if option == subsystem && len(fields) > 4 {
				return fields[4], nil
			}
		}
	}
	return
}
