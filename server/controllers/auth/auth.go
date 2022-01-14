package auth

import (
	"github.com/CosmicBDry/gocmdb/server/controllers/base"
	"github.com/CosmicBDry/gocmdb/server/models"
)

type LoginRequireController struct {
	base.BaseController
	User *models.User
}

func (c *LoginRequireController) Prepare() {
	c.BaseController.Prepare()

	manager := defaultManager

	if user := manager.IsLogin(c); user == nil {
		//fmt.Println("LoginRequireController user is : ", user)
		manager.GoToLogin(c)
	} else {
		c.User = user
		c.Data["user"] = user
		//c.TplName = "test.html"

	}

}

type AuthController struct {
	base.BaseController
}

func (c *AuthController) Login() {

	manager := defaultManager
	manager.Login(c)

}

func (c *AuthController) Logout() {
	manager := defaultManager
	manager.Logout(c)

}
