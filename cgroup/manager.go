package cgroup

import (
	"github.com/sirupsen/logrus"
	"tank/cgroup/subsystem"
)

/**
管理不同的 subsystem
*/

type CgroupManager struct {
	Path     string
	Resource *subsystem.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) Apply(pid int) error {
	for _, subsysIns := range subsystem.SubsystemsIns {
		subsysIns.Apply(c.Path, pid)
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystem.ResourceConfig) error {
	for _, subsysIns := range subsystem.SubsystemsIns {
		subsysIns.Set(c.Path, res)
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _, subsysIns := range subsystem.SubsystemsIns {
		if err := subsysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v ", err)
		}
	}
	return nil
}
