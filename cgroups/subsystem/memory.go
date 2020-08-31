package subsystem

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

const (
	memoryLimitFile = "memory.limit_in_bytes"
)

type MemorySubSystem struct {
}

func (m *MemorySubSystem) Name() string {
	return "memory"
}

func (m *MemorySubSystem) Set(cGroupPath string, res *ResourceLimitConfig) error {
	//获取cGroup中memory的绝对路径
	subsystemCGroupPath, err := GetCGroupPath(m.Name(), cGroupPath, true)
	if err != nil {
		logrus.Errorf("failed to get %s path;err:%v", cGroupPath, err)
		return err
	}
	if res.MemoryLimit != "" {
		// 设置cGroup内存限制，
		// 写入到cGroup memory目录的memory.limit_in_bytes文件中
		err := ioutil.WriteFile(path.Join(subsystemCGroupPath, memoryLimitFile), []byte(res.MemoryLimit), 0644)
		if err != nil {
			logrus.Errorf("failed to write memory;err:%v", err)
			return err
		}
	}
	return nil
}

func (m *MemorySubSystem) Remove(cGroupPath string) error {
	subsystemCGroupPath, err := GetCGroupPath(m.Name(), cGroupPath, true)
	if err != nil {
		logrus.Errorf("failed to remove memory limit;err:%v", err)
		return err
	}
	return os.RemoveAll(subsystemCGroupPath)
}

func (m *MemorySubSystem) Apply(cGroupPath string, pid int) error {
	subsystemCGroupPath, err := GetCGroupPath(m.Name(), cGroupPath, true)
	if err != nil {
		logrus.Errorf("failed to get %s path;err:%v", cGroupPath, err)
		return err
	}
	//将该进程PID写入tasks文件中
	tasksPath := path.Join(subsystemCGroupPath, "tasks")
	err = ioutil.WriteFile(tasksPath, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		logrus.Errorf("failed to write path[%s] pid[%s] to tasks;err:%v", tasksPath, pid, err)
		return err
	}
	return nil
}
