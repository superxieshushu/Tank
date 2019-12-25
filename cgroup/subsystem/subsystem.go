package subsystem

/**
1、Subsystem 接口，每个 Subsystem 可以实现下面 4 个方法
2、这里将 cgroup 抽象成了 path ，原因是 cgroup 在 hierarchy 的路径，便是虚拟文件系统中的虚拟路径
*/
type Subsystem interface {
	//返回 subsystem 的名字
	Name() string

	//设置某个 cgroup 在这个 Subsystem 中的资源限制
	Set(path string, res *ResourceConfig) error

	//将进程添加到某个 cgroup 中
	Apply(path string, pid int) error

	//移除某个 cgroup
	Remove(path string) error
}

var (
	SubsystemsIns = []Subsystem{&MemorySubsystem{}}
)
