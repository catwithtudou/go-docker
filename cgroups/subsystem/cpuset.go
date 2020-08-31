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

const cpuSet = "cpuset.cpus"

type CpuSetSubSystem struct {
	apply bool
}

func (*CpuSetSubSystem) Name() string {
	return "cpuset"
}

func (c *CpuSetSubSystem) Set(cGroupPath string, res *ResourceLimitConfig) error {
	subsystemCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, true)
	if err != nil {
		logrus.Errorf("failed to get path[%s];err: %v", cGroupPath, err)
		return err
	}

	if res.CpuSet != "" {
		err := ioutil.WriteFile(path.Join(subsystemCGroupPath, cpuSet), []byte(res.CpuSet), 0644)
		if err != nil {
			logrus.Errorf("failed to write file %s;err: %+v", cpuSet, err)
			return err
		}
		c.apply = true
	}
	return nil
}

func (c *CpuSetSubSystem) Remove(cGroupPath string) error {
	subsystemCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsystemCGroupPath)
}

func (c *CpuSetSubSystem) Apply(cGroupPath string, pid int) error {
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
