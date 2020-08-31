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

const cpuShare = "cpu.shares"

type CpuSubSystem struct {
	apply bool
}

func (*CpuSubSystem) Name() string {
	return "cpu"
}

func (c *CpuSubSystem) Set(cGroupPath string, res *ResourceLimitConfig) error {
	subsystemCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, true)
	if err != nil {
		logrus.Errorf("failed to get path[%s];err: %v", cGroupPath, err)
		return err
	}
	if res.CpuShare != "" {
		err = ioutil.WriteFile(path.Join(subsystemCGroupPath, cpuShare), []byte(res.CpuShare), 0644)
		if err != nil {
			logrus.Errorf("failed to write file %s;err: %+v", cpuShare, err)
			return err
		}
		c.apply = true
	}
	return nil
}

func (c *CpuSubSystem) Remove(cGroupPath string) error {
	subsystemCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false)
	if err != nil {
		logrus.Errorf("failed to get path[%s];err: %v", cGroupPath, err)
		return err
	}
	return os.RemoveAll(subsystemCGroupPath)
}

func (c *CpuSubSystem) Apply(cGroupPath string, pid int) error {
	if c.apply {
		subsystemCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false)
		if err != nil {
			logrus.Errorf("failed to get path[%s];err: %v", cGroupPath, err)
			return err
		}

		tasksPath := path.Join(subsystemCGroupPath, "tasks")
		err = ioutil.WriteFile(tasksPath, []byte(strconv.Itoa(pid)), os.ModePerm)
		if err != nil {
			logrus.Errorf("failed to write path[%s] pid[%s] to tasks;err: %v", tasksPath, pid, err)
			return err
		}
	}
	return nil
}
