package cobraCmd

import (
	"fmt"
	"os"
	"time"

	"github.com/CosmicBDry/gocmdb/server/cloud"
	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/spf13/cobra"
)

var (
	cloudsync bool
)

var SyncCmd = &cobra.Command{
	Use:   "cloudsync",
	Short: "lauch a cloudsync process ",
	Long:  "lauch a cloudsync process ",
	RunE: func(cmd *cobra.Command, args []string) error {

		orm.RegisterDriver("mysql", orm.DRMySQL)
		orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("dsn"))

		if db, err := orm.GetDB(); err != nil || db.Ping() != nil {
			beego.Error("数据库连接失败")
			os.Exit(-1)
		}

		beego.SetLogger("file", `{"level":7,"filename":"logs/cloudsync.log","maxdays":15,"maxlines":1000}`)

		//通过遍历time.Tick(n * time.second)函数实现每间隔n秒后执行以下for循环中的操作
		for now := range time.Tick(5 * time.Second) {
			fmt.Println(now)
			_, _, Platforms := models.DefautlCloudPlatformManager.QueryList(0, 0, 0, "", "", "", "")
			for _, platform := range Platforms {
				if platform.Status == models.StatusLock {
					continue
				}
				//fmt.Println(platform)
				sdk := cloud.DefaultCloudManager.Plugins[platform.Types]
				sdk.Init(platform.Addr, platform.Region, platform.AccessKey, platform.SecretKey)
				instances := sdk.GetInstances()
				fmt.Printf("%#v\n", instances)
				if len(instances) == 0 {

					models.DefaultVirtualMachineManager.SyncVmStatus(&platform, now) //同步虚拟机状态
					models.DefautlCloudPlatformManager.CloudSyncTime(&platform, now) //云平台最后一次同步时间
				}
				for _, instance := range instances {
					if err := models.DefaultVirtualMachineManager.SyncVmByCloudPlatform(&platform, instance); err == nil {
						models.DefaultVirtualMachineManager.SyncVmStatus(&platform, now) //同步虚拟机状态
						models.DefautlCloudPlatformManager.CloudSyncTime(&platform, now) //云平台最后一次同步时间
					}
				}
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(SyncCmd)

	SyncCmd.Flags().BoolVarP(&cloudsync, "cloudsync", "s", false, "sync cloud machine")

}
