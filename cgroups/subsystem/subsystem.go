package subsystem

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

//资源限制配置
type ResourceLimitConfig struct {
	//内存限制
	MemoryLimit string
	//CPU时间片权重
	CpuShare string
	//CPU核数
	CpuSet string
}

//在hierarchy中cGroup便是虚拟的路径地址
type Subsystem interface {
	//返回subsystem名字
	Name() string
	//设置资源限制
	Set(cGroupPath string, res *ResourceLimitConfig) error
	//删除资源限制
	Remove(cGroup string) error
	//进程应用此cGroup配置
	Apply(cGroupPath string, pid int) error
}

var (
	Subsystems = []Subsystem{
		&MemorySubSystem{},
		&CpuSubSystem{},
		&CpuSetSubSystem{},
	}
)
