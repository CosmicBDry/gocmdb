package forms

import (
	"strings"

	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/astaxie/beego/validation"
)

type LoginForm struct {
	Name     string `form:"username"`
	Password string `form:"password"`
	User     *models.User
}

func (f *LoginForm) Valid(v *validation.Validation) {

	f.Name = strings.TrimSpace(f.Name)
	f.Password = strings.TrimSpace(f.Password)

	if f.Name == "" || f.Password == "" {
		v.SetError("error", "用户名或密码不能为空!")
	} else if user := models.DefautlUserManager.GetUserByName(f.Name); user == nil || !user.ValidPassword(f.Password) {
		v.SetError("error", "用户名或密码错误!")
	} else if user.IsLock() {
		v.SetError("error", "用户被锁定!")
	} else {
		f.User = user
	}
}
