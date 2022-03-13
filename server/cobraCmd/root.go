package cobraCmd

//子命令的初始化包都可以放到此处，如下
import (
	_ "github.com/CosmicBDry/gocmdb/server/cloud/plugins"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cmdb [web]",
	Short: "Asset management system[CMDB]",
	Long:  "Asset management system[CMDB]",
}

func Execute() {
	RootCmd.Execute()
}
