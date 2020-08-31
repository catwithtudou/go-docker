package cgroups

import (
	"github.com/sirupsen/logrus"
	"go-docker/cgroups/subsystem"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

//资源管理器
type CGroupManager struct {
	Path string
}

func NewCGroupManager(path string) *CGroupManager {
	return &CGroupManager{Path: path}
}

func (c *CGroupManager) Set(res *subsystem.ResourceLimitConfig) {
	for _, reSubsystem := range subsystem.Subsystems {
		err := reSubsystem.Set(c.Path, res)
		if err != nil {
			logrus.Errorf("failed to set %s;err:%v", reSubsystem.Name(), err)
		}
	}
}

func (c *CGroupManager) Destroy() {
	for _, reSubsystem := range subsystem.Subsystems {
		err := reSubsystem.Remove(c.Path)
		if err != nil {
			logrus.Errorf("failed to remove %s;err:%v", reSubsystem.Name(), err)
		}
	}
}

func (c *CGroupManager) Apply(pid int) {
	for _, reSubsystem := range subsystem.Subsystems {
		err := reSubsystem.Apply(c.Path, pid)
		if err != nil {
			logrus.Errorf("failed to apply %s;err:%v", reSubsystem.Name(), err)
		}
	}
}
