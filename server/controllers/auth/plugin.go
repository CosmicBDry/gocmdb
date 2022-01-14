package auth

import (
	"net/http"
	"strings"

	"github.com/CosmicBDry/gocmdb/server/forms"
	"github.com/CosmicBDry/gocmdb/server/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
)

//session插件
type Session struct {
}

func (s *Session) Name() string {
	return "session"
}

func (s *Session) Is(c *context.Context) bool {
	return c.Input.Header("Authentication") == ""
}

func (s *Session) Login(c *AuthController) bool {

	valid := &validation.Validation{}
	form := &forms.LoginForm{}
	c.ParseForm(form)
	if c.Ctx.Input.IsPost() {

		if ok, err := valid.Valid(form); err == nil && ok {
			//fmt.Println("form userid: ", form.User.Id)
			c.SetSession("userID", form.User.Id)
			c.Redirect(c.URLFor(beego.AppConfig.String("Home")), http.StatusFound)
			return true
		}

	}
	c.Data["form"] = form
	c.Data["error"] = valid.Errors
	c.TplName = "login/login.html"
	return false
}

func (s *Session) Logout(c *AuthController) {
	c.DestroySession()
	c.Redirect(c.URLFor(beego.AppConfig.String("Login")), http.StatusFound)
}

func (s *Session) IsLogin(c *LoginRequireController) *models.User {
	if session := c.GetSession("userID"); session != nil {
		if user := models.DefautlUserManager.GetUserById(session.(int)); user != nil {
			return user
		}
	}

	return nil

}

func (s *Session) GoToLogin(c *LoginRequireController) {

	c.Redirect(c.URLFor(beego.AppConfig.String("Login")), http.StatusFound)
}

type Token struct {
}

//Token插件

func (t *Token) Name() string {

	return "token"

}
func (t *Token) Is(c *context.Context) bool {

	return c.Input.Header("Authentication") == "Token"
}

func (t *Token) IsLogin(c *LoginRequireController) *models.User {
	accesskey := strings.TrimSpace(c.Ctx.Input.Header("AccessKey"))
	secretkey := strings.TrimSpace(c.Ctx.Input.Header("SecretKey"))
	if token := models.DefaultTokenManager.GetByKey(accesskey, secretkey); token != nil {
		return token.User
	}
	return nil
}

func (t *Token) Login(c *AuthController) bool {
	json := map[string]interface{}{
		"code":   400,
		"result": "该请求路径为非Token验证方式!",
	}
	c.Data["json"] = json
	c.ServeJSON()
	return false
}

func (t *Token) Logout(c *AuthController) {
	json := map[string]interface{}{
		"code":   400,
		"result": "该请求路径为非Token验证方式!",
	}
	c.Data["json"] = json
	c.ServeJSON()
}

func (s *Token) GoToLogin(c *LoginRequireController) {
	json := map[string]interface{}{
		"code":   403,
		"result": "没有权限，请提供有效的Token!",
	}
	c.Data["json"] = json
	c.ServeJSON()
}

func init() {
	defaultManager.Register(new(Session))
	defaultManager.Register(new(Token))
}
