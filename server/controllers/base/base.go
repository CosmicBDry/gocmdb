package base

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prepare() {

	c.Data["_xsrf"] = c.XSRFToken()
}
