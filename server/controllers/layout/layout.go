package layout

import (
	"github.com/CosmicBDry/gocmdb/server/controllers/auth"
)

type LayoutController struct {
	auth.LoginRequireController
}

func (c *LayoutController) Prepare() {
	c.LoginRequireController.Prepare()
	c.Data["menu"] = ""
	c.Data["expand"] = ""
	c.Layout = "layout/base.html"
	c.LayoutSections = make(map[string]string) //初始化LayoutSections，因映射使用前必须先初始化
	c.LayoutSections["LayoutStyle"] = ""
	c.LayoutSections["LayoutScript"] = ""
}
