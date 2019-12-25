package subsystem

//用于传递资源限制配置的结构体、包含内存限制、CPU 时间片权重、CPU 核心数目
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}
