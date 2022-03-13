package cobraCmd

import (
	"os"
	"time"

	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/CosmicBDry/gocmdb/server/routers"
	"github.com/CosmicBDry/gocmdb/server/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/spf13/cobra"
)

var (
	initial, syncdb, force, verbose, resetPassword bool
)

var WebCmd = &cobra.Command{
	Use:   "web",
	Short: "lauch web program",
	Long:  "lauch web program",
	RunE: func(cmd *cobra.Command, args []string) error {

		orm.RegisterDriver("mysql", orm.DRMySQL)
		orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("dsn"))

		if db, err := orm.GetDB(); err != nil || db.Ping() != nil {
			beego.Error("数据库连接失败")
			os.Exit(-1)
		}
		routers.Register()
		beego.SetLogger("file", `{"level":7,"filename":"logs/web.log","maxdays":15,"maxlines":10000}`)
		switch {
		case initial:
			orm.RunSyncdb("default", force, verbose)
			Birth := time.Now()
			admin := &models.User{Name: "admin", IsSuperman: true, Birthday: &Birth}
			ormer := orm.NewOrm()
			if err := ormer.Read(admin, "Name"); err == orm.ErrNoRows {
				Password := utils.RandString(20)
				admin.SetPassword(Password)
				if _, err := ormer.Insert(admin); err == nil {
					beego.Informational("Database table and user 'admin' have been initialized , admin password is: ", Password)
				} else {

					beego.Informational("admin initialization failed!")

				}

			} else {

				beego.Informational("Database table has been initialized,but admin is already exist,admin does not need to be initialized!")
			}
		case syncdb:
			orm.RunSyncdb("default", force, verbose)
			beego.Informational("sync table success!")
		case resetPassword:
			admin := &models.User{Id: 1}
			password := utils.RandString(20)
			admin.SetPassword(password)
			ormer := orm.NewOrm()
			if rows, err := ormer.Update(admin, "Password"); rows != 0 && err == nil {
				beego.Informational("User 'admin' reset success,admin password is : ", password)
			} else {
				beego.Error("User 'admin' Password reset failed,Ensure that the user 'admin' exists or please check the database connection!")

			}

		default:
			beego.Run()
		}

		return nil
	},
}

func init() {

	RootCmd.AddCommand(WebCmd)
	WebCmd.Flags().BoolVarP(&initial, "init", "I", false, "sync table,according to your needs,Adding '-force' option to reset table and admin password")
	WebCmd.Flags().BoolVarP(&syncdb, "sync", "S", false, "sync table")
	WebCmd.Flags().BoolVarP(&force, "force", "F", false, "force sync db(drop table)")
	WebCmd.Flags().BoolVarP(&verbose, "verbose", "V", false, "display detailed information of sql")
	WebCmd.Flags().BoolVarP(&resetPassword, "resetPassword", "R", false, "only reset admin password!")

}
