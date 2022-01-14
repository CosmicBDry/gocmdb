package routers

import (
	"github.com/astaxie/beego"

	"github.com/CosmicBDry/gocmdb/server/controllers/auth"
	"github.com/CosmicBDry/gocmdb/server/controllers/cloud"
	"github.com/CosmicBDry/gocmdb/server/controllers/user"
)

func Register() {

	beego.AutoRouter(&auth.AuthController{})

	beego.AutoRouter(&user.UserPageController{})

	beego.AutoRouter(&user.UserController{})
	beego.AutoRouter(&cloud.CloudPlatformPageController{})
	beego.AutoRouter(&cloud.CloudPlatformController{})
	beego.AutoRouter(&cloud.VirtualMachinePageController{})
	beego.AutoRouter(&cloud.VirtualMachineController{})
}
