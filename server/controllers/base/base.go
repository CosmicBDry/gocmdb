package base

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

//基础控制器获取xsrf的token，此时所有其他控制器html中都可以获取到xsrf
/*func (c *BaseController) Prepare() {

	c.Data["_xsrf"] = c.XSRFToken()
}*/
