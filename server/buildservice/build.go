package buildservice

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/CosmicBDry/gocmdb/server/rwFile"
	"github.com/CosmicBDry/gocmdb/server/systemCommand"

	"github.com/astaxie/beego"
)

type BuildService struct {
	Mutex     sync.RWMutex   //读写锁
	ProjectId map[int64]bool //判断项目是否再构建中
}

func NewBuildService() *BuildService {

	return &BuildService{
		Mutex:     sync.RWMutex{},
		ProjectId: make(map[int64]bool),
	}
}

var DefaultBuildService = NewBuildService()

func (b *BuildService) Build(gitproject *models.GitProject) error {
	b.Mutex.RLock() //添加一个读锁,每次只允许一个进程读取b.ProjectId映射
	if _, ok := b.ProjectId[gitproject.ProjectId]; ok {
		b.Mutex.RUnlock() //读完之后将锁释放，使得其他进程可读
		beego.Warning(gitproject.SShUrl + ": 项目正在构建中...,请不要重复构建!")
		return fmt.Errorf("Is Building")
	}
	b.Mutex.RUnlock() ////读完之后将锁释放
	b.ProjectId[gitproject.ProjectId] = true

	//构建脚本：------------------------------------------------------------------------------->
	WorkDir := beego.AppConfig.DefaultString("builder::BuildDir", "/tmp/")
	ProjectDir := fmt.Sprintf("%sProject_%s/", WorkDir, gitproject.Name)
	PackageDir := beego.AppConfig.DefaultString("builder::PackageDir", "/tmp/")
	Now := time.Now().Format("2006-01-02_15-04-05")
	PackageCmd := fmt.Sprintf("cd %s%s/%s/target && tar czf %s%s/%s.tar.gz ./* && wait && sleep 5\n", ProjectDir, gitproject.Branch, gitproject.Name, PackageDir, gitproject.Name, Now)
	if gitproject.PackageFile != "" {
		PackageCmd = fmt.Sprintf("cd %s%s/%s && cp -a {%s,} ./target/ && wait && sleep 5 &&cd ./target && tar czf %s%s/%s.tar.gz ./* && wait && sleep 5\n", ProjectDir, gitproject.Branch, gitproject.Name, gitproject.PackageFile, PackageDir, gitproject.Name, Now)
		fmt.Println("----->: ", PackageCmd)
	}
	cmd := []string{
		fmt.Sprintf("rm -rf %s && wait && mkdir -p %s%s && echo `%s` >> %s%s/build.log && wait\n", ProjectDir, ProjectDir, gitproject.Branch, "date +%Y-%m-%d_%H:%M:%S", ProjectDir, gitproject.Branch),
		fmt.Sprintf("mkdir -p %s%s\n", PackageDir, gitproject.Name),
		fmt.Sprintf("cd %s%s &&git clone -b %s %s %s && wait && sleep 5\n", ProjectDir, gitproject.Branch, gitproject.Branch, gitproject.SShUrl, gitproject.Name),
		fmt.Sprintf("cd %s%s/%s && go build -o ./target/%s && wait && sleep 5\n", ProjectDir, gitproject.Branch, gitproject.Name, gitproject.Name),
		PackageCmd,
	}

	for _, v := range cmd {
		var result string
		var err error
		path := fmt.Sprintf("%s%s/build.log", ProjectDir, gitproject.Branch)
		if result, err = systemCommand.RunCmd(v); err != nil {
			beego.Error(gitproject.SShUrl + ": 项目代码构建失败,详情请查看build.log构建记录")
			if errs := rwFile.WriteFile(result+err.Error()+"\n", path); errs != nil {
				beego.Error(path + " :代码构建记录写入失败," + errs.Error())
			}
			delete(b.ProjectId, gitproject.ProjectId)
			return err
		}

		if result == "" {
			result = "执行成功: " + v
		}

		if err = rwFile.WriteFile(result, path); err != nil {
			beego.Error(path + " :代码构建记录写入失败," + err.Error())
			continue
		}
	}
	beego.Informational(gitproject.SShUrl + ": 项目构建成功!")
	delete(b.ProjectId, gitproject.ProjectId)
	return nil
}

func GetBuildLog(gitproject *models.GitProject) string {
	WorkDir := beego.AppConfig.DefaultString("builder::BuildDir", "/tmp/")
	ProjectDir := fmt.Sprintf("%sProject_%s/", WorkDir, gitproject.Name)
	Path := filepath.Join(ProjectDir, gitproject.Branch, "/", "build.log")
	return rwFile.ReadFile(Path)
}
