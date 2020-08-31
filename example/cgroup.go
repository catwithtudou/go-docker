package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

const (
	//挂载memory subsystem的hierarchy的位置
	cGroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"
	//通过访问/proc/self/目录来获取自己的系统信息
	//获取指定进程的信息，例如内存映射、CPU绑定信息等等
	procSelfExe = "/proc/self/exe"
	//创建的子cGroup文件名
	memoryCGroup = "cgroup-demo-memory"
	//限制内存使用量
	limitMemory = "100m"
)

func main() {

	if os.Args[0] == procSelfExe {
		//容器进程
		containerPid := syscall.Getpid()

		fmt.Printf("current pid %d \n", containerPid)

		cmd := exec.Command("sh", "-c", "stress --vm-bytes 200m --vm-keep -m 1")
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	cmd := exec.Command(procSelfExe)
	//限制NameSpace资源
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	//获取fork的进程映射在外部命名空间的pid
	processPid := cmd.Process.Pid
	fmt.Printf("process pid : %+v", processPid)

	//创建子cGroup
	childCGroup := path.Join(cGroupMemoryHierarchyMount, memoryCGroup)
	if err := os.Mkdir(childCGroup, 0755); err != nil {
		panic(err)
	}

	//将容器进程PID放到子cGroup中
	if err := ioutil.WriteFile(path.Join(childCGroup, "tasks"), []byte(strconv.Itoa(processPid)), 0644); err != nil {
		panic(err)
	}

	//限制其memory内存使用
	if err := ioutil.WriteFile(path.Join(childCGroup, "memory.limit_in_bytes"), []byte(limitMemory), 0644); err != nil {
		panic(err)
	}

	_, err = cmd.Process.Wait()
	if err != nil {
		panic(err)
	}

}
