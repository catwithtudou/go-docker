package config

/**
 *@Author tudou
 *@Date 2020/9/5
 **/

const (
	Running = "running"
	Stop    = "stopped"
	Exit    = "exited"
)

const (
	WriteLayer = "/writeLayer"
	RootPath   = "/root"
	MntPath    = "/mnt"
	BinPath    = "/bin/"
)

const (
	DefaultContainerInfoPath = "/var/run/go-docker/"
	ContainerInfoFileName    = "config.json"
	ContainerLogFileName     = "container.log"
)

const (
	EnvExecPid = "docker_pid"
	EnvExecCmd = "docker_cmd"
)

const (
	DefaultNetworkPath   = "/var/run/go-docker/network/network/"
	DefaultAllocatorPath = "/var/run/go-docker/network/ipam/subnet.json"
)