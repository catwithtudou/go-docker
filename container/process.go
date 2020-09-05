package container

import (
	"github.com/sirupsen/logrus"
	"go-docker/config"
	"os"
	"os/exec"
	"path"
	"syscall"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

const (
	//通过访问/proc/self/目录来获取自己的系统信息
	//获取指定进程的信息，例如内存映射、CPU绑定信息等等
	procSelfExe = "/proc/self/exe"
)

//创建一个被namespace隔离的进程command
func NewParentProcess(tty bool,containerName, imageName,volume string, envs []string) (*exec.Cmd, *os.File) {
	//调用syscall包的Pipe()函数
	//获取读和写阻塞管道
	readPipe, writePipe, _ := os.Pipe()
	//获取系统信息传入init参数
	cmd := exec.Command(procSelfExe, "init")
	//限制资源
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	//前台进程
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}else{
		logDir := path.Join(config.DefaultContainerInfoPath,containerName)
		if _,err := os.Stat(logDir);err!=nil && os.IsNotExist(err){
			err:=os.MkdirAll(logDir,os.ModePerm)
			if err!=nil{
				logrus.Errorf("failed to mkdir log file;err: %v",err)
			}
		}

		logFileName := path.Join(logDir,config.ContainerLogFileName)
		file,err := os.Create(logFileName)
		if err!=nil{
			logrus.Errorf("failed to crate log file;err: %v",err)
		}
		cmd.Stdout = file
	}


	cmd.ExtraFiles = []*os.File{
		readPipe, //读管道
	}

	//设置环境变量
	cmd.Env = append(os.Environ(),envs...)
	err:=NewWorkSpace(volume,containerName,imageName)
	if err!=nil{
		logrus.Errorf("failed to create work space,err: %v",err)
	}

	cmd.Dir = path.Join(config.RootPath,config.MntPath,containerName)

	return cmd, writePipe
}
