package subsystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubsystem struct {
}

func (m *MemorySubsystem) Set(cgroupPath string, res *ResourceConfig) (err error) {
	if subsysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			//设置这个 cgroup 的内存限制，即将限制写入到 cgroup 对应目录的 memory.limit in bytes 文件中。
			if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
	}
	return
}

func (m *MemorySubsystem) Remove(cgroupPath string) (err error) {
	if subsysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false); err == nil {
		return os.Remove(subsysCgroupPath)
	} else {
		return err
	}
}

func (m *MemorySubsystem) Apply(cgroupPath string, pid int) (err error) {
	if subsysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false); err == nil {
		if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v ", err)
		}
		return err
	} else {
		return fmt.Errorf("get cgroup %s error : %v", cgroupPath, err)
	}
}

func (m *MemorySubsystem) Name() (name string) {
	return "memory"
}
