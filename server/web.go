package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/CosmicBDry/gocmdb/server/cloud/plugins"
	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/CosmicBDry/gocmdb/server/routers"
	"github.com/CosmicBDry/gocmdb/server/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	//fmt.Println(beego.AppConfig.String("dsn"))
	fmt.Println(beego.AppConfig.String("runmode"))

	h := flag.Bool("h", false, "help manual")
	init := flag.Bool("init", false, "sync table,according to your needs,Adding '-force' option to reset table and admin password")
	syncdb := flag.Bool("sync", false, "sync table")
	force := flag.Bool("force", false, "force sync db(drop table)")
	verbose := flag.Bool("v", false, "display detailed information of sql")
	resetPassword := flag.Bool("resetPassword", false, "only reset admin password!")
	flag.Usage = func() {
		fmt.Println("usage: go web.go -h")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *h {

		flag.Usage()
		os.Exit(-1)

	}

	if !*verbose {
		//orm.Debug = true
		beego.BeeLogger.DelLogger("console")
	}

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("dsn"))

	if db, err := orm.GetDB(); err != nil || db.Ping() != nil {
		beego.Error("数据库连接失败")
		os.Exit(-1)
	}

	routers.Register()

	beego.SetLogger("file", `{"level":7,"filename":"logs/web.log","maxdays":15,"maxlines":1000}`)

	switch {
	case *init:
		orm.RunSyncdb("default", *force, *verbose)
		Birth := time.Now()
		admin := &models.User{Name: "admin", IsSuperman: true, Birthday: &Birth}
		ormer := orm.NewOrm()
		if err := ormer.Read(admin, "Name"); err == orm.ErrNoRows {
			Password := utils.RandString(20)
			admin.SetPassword(Password)
			if _, err := ormer.Insert(admin); err == nil {
				beego.Informational("Database table and user 'admin' have been initialized , admin password is: ", Password)
			} else {
				fmt.Println(err)
				beego.Informational("admin initialization failed!")
			}

		} else {

			beego.Informational("Database table has been initialized,but admin is already exist,admin does not need to be initialized!")
		}

		//fmt.Println("init webserver!")
	case *syncdb:
		orm.RunSyncdb("default", *force, *verbose)
		beego.Informational("sync table success!")
	case *resetPassword:
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
		//models.DefaultTokenManager.GetByKey("abcd", "qwer")
		//fmt.Printf("%#v\n", models.DefautlUserManager.GetUsers()[0])
		/*a := time.Now()
		createuser := &models.User{Name: "generalman123", IsSuperman: false, Birthday: &a, Addr: "北京市昌平区北七家镇宏福苑社区99号", Email: "general123@fox.mail.cn"}
		Password := utils.RandString(20)
		ormer := orm.NewOrm()
		createuser.SetPassword(Password)
		fmt.Println(Password)
		ormer.Insert(createuser)*/
		beego.Run()
	}

}
