package container

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

//容器执行的第一个进程
//使用mount挂载proc文件系统(方便ps等系统命令)
//
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return errors.New("failed to get user command in run container")
	}
	//使用mount挂载
	err := setUpMount()
	if err != nil {
		logrus.Errorf("set up mount error: %v", err)
		return err
	}

	//获取PATH路径
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("look %s path error: %v", err)
		return err
	}

	//在PATH路径下执行command
	err = syscall.Exec(path, cmdArray[0:], os.Environ())
	if err != nil {
		logrus.Errorf("failed to exec command;err:%v", err)
		return err
	}
	return nil
}

//获取用户输入的参数命令
func readUserCommand() []string {
	//index 3 的文件描述符即 readPipe
	readPipe := os.NewFile(uintptr(3), "pipe")
	//读取读管道
	bs, err := ioutil.ReadAll(readPipe)
	if err != nil {
		logrus.Errorf("read pipe error: %v", err)
		return nil
	}
	msg := string(bs)
	return strings.Split(msg, "")
}

//使用mount挂载proc文件系统
func setUpMount() (err error) {
	//声明新的mount namespaced独立
	//挂载状态要先从共享挂载编程私有挂载
	err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		logrus.Errorf("failed to mount independence;err:%v", err)
		return
	}
	//挂载proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		logrus.Errorf("failed to mount proc;err:%v", err)
		return
	}
	return
}
