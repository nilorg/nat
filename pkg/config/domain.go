package config

import (
	"io/fs"
	"os"
	"path"
	"runtime"
)

// windows写配置文件到用户目录下隐藏文件,~/.nilorg/nat/domain/www.nilorg.com.conf
// linux写配置文件到/etc/nilorg/nat/domain/www.nilorg.com.conf

// 文件后缀
const FileSuffix = ".conf"

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func getConfigPath() (dir string, err error) {
	// 判断系统类型
	if isWindows() {
		var homeDir string
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return
		}
		dir = path.Join(homeDir, ".nilorg", "nat", "domain")
	} else {
		dir = "/etc/nilorg/nat/domain/"
	}
	return
}

type Domain struct {
	Domain string
	Ip     string
}

func SetDomain(d *Domain) (err error) {
	var confPath string
	confPath, err = getConfigPath()
	if err != nil {
		return
	}
	var exists bool
	exists, err = isDirExists(confPath)
	if err != nil {
		return
	}
	if !exists {
		err = os.MkdirAll(confPath, 0777)
		if err != nil {
			return
		}
	}
	// 写入配置文件
	confPath = path.Join(confPath, d.Domain+FileSuffix)
	err = os.WriteFile(confPath, []byte(d.Ip), 0777)
	return
}

func GetDomain(domain string) (d *Domain, err error) {
	var confPath string
	confPath, err = getConfigPath()
	if err != nil {
		return
	}
	var exists bool
	exists, err = isDirExists(confPath)
	if err != nil {
		return
	}
	if !exists {
		return
	}
	// 读取配置文件
	confPath = path.Join(confPath, domain+FileSuffix)
	var buf []byte
	buf, err = os.ReadFile(confPath)
	if err != nil {
		return
	}
	d = &Domain{
		Domain: domain,
		Ip:     string(buf),
	}
	return
}

func DelDomain(domain string) (err error) {
	var confPath string
	confPath, err = getConfigPath()
	if err != nil {
		return
	}
	var exists bool
	exists, err = isDirExists(confPath)
	if err != nil {
		return
	}
	if !exists {
		return
	}
	// 删除配置文件
	confPath = path.Join(confPath, domain+".conf")
	err = os.Remove(confPath)
	return
}

func GetDomains() (ds []*Domain, err error) {
	var confPath string
	confPath, err = getConfigPath()
	if err != nil {
		return
	}
	var exists bool
	exists, err = isDirExists(confPath)
	if err != nil {
		return
	}
	if !exists {
		return
	}
	// 读取配置文件
	var files []fs.DirEntry
	files, err = os.ReadDir(confPath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		var buf []byte
		buf, err = os.ReadFile(path.Join(confPath, file.Name()))
		if err != nil {
			return
		}
		ds = append(ds, &Domain{
			Domain: file.Name()[:len(file.Name())-len(FileSuffix)],
			Ip:     string(buf),
		})
	}
	return
}

// 判断目录是否存在
func isDirExists(dir string) (exists bool, err error) {
	var fileInfo os.FileInfo
	fileInfo, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			exists = false
			err = nil
		}
		return
	}
	exists = fileInfo.IsDir()
	return
}
