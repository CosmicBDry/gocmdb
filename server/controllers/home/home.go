package home

import (
	"github.com/CosmicBDry/gocmdb/server/controllers/auth"
)

type HomeController struct {
	auth.LoginRequireController
}

func (c *HomeController) Page() {
	//c.TplName = "test/test.html"
	c.TplName = "home/homePage.html" //只有最后一个才能生效
}
