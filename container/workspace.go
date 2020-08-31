package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
)

/**
 *@Author tudou
 *@Date 2020/8/31
 **/

const (
	WriteLayer = "writeLayer"
	RootPath   = "/root/"
	MntPath    = "/root/mnt"
	BinPath    = "/bin/"
)

//创建容器运行目录
func NewWorkSpace(rootPath string, mntPath string, writeLayer string, volume string, containerName string) error {
	//创建只读层
	err := createReadOnlyLayer(rootPath, containerName)
	if err != nil {
		logrus.Errorf("failed to create read only layer;err: %v", err)
		return err
	}
	//创建读写层
	err = createWriteLayer(rootPath, writeLayer)
	if err != nil {
		logrus.Errorf("failed to create write layer;err: %v", err)
		return err
	}
	//将只读层与读写层指定到创建的挂载点
	err = createMountPoint(rootPath, mntPath, writeLayer, containerName)
	if err != nil {
		logrus.Errorf("create mount point, err: %v", err)
		return err
	}
	//设置宿主机与容器文件映射
	return mountVolume(rootPath, mntPath, volume)
}

func createReadOnlyLayer(rootPath string, containerName string) error {
	containerPath := path.Join(rootPath, containerName)
	_, err := os.Stat(containerPath)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(containerPath, os.ModePerm)
		if err != nil {
			logrus.Errorf("failed to mkdir containerPath[%s];err: %v", containerName, err)
			return err
		}
	}
	//解压
	containerTarPath := path.Join(rootPath, containerName+".tar")
	if _, err = exec.Command("tar", "-xvf", containerTarPath, "-C", containerPath).CombinedOutput(); err != nil {
		logrus.Errorf("failed to tar %s.tar;err: %v", containerName, err)
		return err
	}
	return nil
}

func createWriteLayer(rootPath string, rewriteLayerPath string) error {
	writeLayerPath := path.Join(rootPath, rewriteLayerPath)
	_, err := os.Stat(writeLayerPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(writeLayerPath, os.ModePerm)
		if err != nil {
			logrus.Errorf("failed to mkdir writeLayer[%s];err: %v", writeLayerPath, err)
			return err
		}
	}
	return nil
}

func createMountPoint(rootPath string, reMntPath string, writeLayer string, containerName string) error {
	mntPath := path.Join(rootPath, reMntPath)
	_, err := os.Stat(mntPath)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(mntPath, os.ModePerm)
		if err != nil {
			logrus.Errorf("failed to mkdir mnt path[%s];err: %v", mntPath, err)
			return err
		}
	}

	dirs := fmt.Sprintf("dirs=%s%s:%s%s", rootPath, writeLayer, rootPath, containerName)
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntPath)
	if err := cmd.Run(); err != nil {
		logrus.Errorf("failed to mnt cmd[%s] run;err: %v", cmd, err)
		return err
	}
	return nil
}

func mountVolume(rootPath, mntPath, volume string) error {
	if volume != "" {
		volumes := strings.Split(volume, ":")
		if len(volumes) > 1 {
			// 创建宿主机中文件路径
			parentPath := volumes[0]
			if _, err := os.Stat(parentPath); err != nil && os.IsNotExist(err) {
				if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
					logrus.Errorf("failed to mkdir parent path[%s];err: %v", parentPath, err)
					return err
				}
			}

			// 创建容器中的挂载点
			containerPath := volumes[1]
			containerVolumePath := path.Join(rootPath, mntPath, containerPath)
			if _, err := os.Stat(containerVolumePath); err != nil && os.IsNotExist(err) {
				if err := os.MkdirAll(containerVolumePath, os.ModePerm); err != nil {
					logrus.Errorf("failed to mkdir volume path[%s];err: %v", containerVolumePath, err)
					return err
				}
			}

			// 把宿主机文件目录挂载到容器挂载点中
			dirs := fmt.Sprintf("dirs=%s", parentPath)
			cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumePath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				logrus.Errorf("failed to mount cmd run;err: %v", err)
				return err
			}
		}
	}
	return nil
}

// 删除容器workspace
func DeleteWorkSpace(rootPath, mntPath, writeLayerPath, volume string) error {
	// 卸载挂载点
	err := unMountPoint(rootPath, mntPath)
	if err != nil {
		return err
	}
	// 删除读写层
	err = deleteWriteLayer(rootPath, writeLayerPath)
	if err != nil {
		return err
	}
	// 删除宿主机与文件系统映射
	deleteVolume(rootPath, mntPath, volume)
	return nil
}

func unMountPoint(rootPath, mntPath string) error {
	reMntPath := path.Join(rootPath, mntPath)
	if _, err := exec.Command("umount", reMntPath).CombinedOutput(); err != nil {
		logrus.Errorf("failed to unmount mnt[%s];err: %v", reMntPath, err)
		return err
	}
	err := os.RemoveAll(reMntPath)
	if err != nil {
		logrus.Errorf("failed to remove mnt path[%s];err: %v", reMntPath, err)
		return err
	}
	return nil
}

func deleteWriteLayer(rootPath, writeLayer string) error {
	writerLayerPath := path.Join(rootPath, writeLayer)
	err := os.RemoveAll(writerLayerPath)
	if err != nil {
		logrus.Errorf("failed to remove write layer path[%s];err: %v", writeLayer, err)
		return err
	}
	return nil
}

func deleteVolume(rootPath, mntPath, volume string) {
	if volume != "" {
		volumes := strings.Split(volume, ":")
		if len(volumes) > 1 {
			containerPath := path.Join(rootPath, mntPath, volumes[1])
			if _, err := exec.Command("umount", containerPath).CombinedOutput(); err != nil {
				logrus.Errorf("failed to unmount container path[%s];err: %v", containerPath, err)
			}
		}
	}
}
