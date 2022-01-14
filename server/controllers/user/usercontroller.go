package user

import (
	"fmt"

	"github.com/CosmicBDry/gocmdb/server/controllers/auth"
	"github.com/CosmicBDry/gocmdb/server/controllers/layout"
	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/CosmicBDry/gocmdb/server/utils"
	"github.com/astaxie/beego/validation"
)

type UserPageController struct {
	layout.LayoutController
}

func (c *UserPageController) Index() {
	c.Data["menu"] = "UserManager"     //自动选择打开用户管理的页面菜单
	c.Data["expand"] = "SystemManager" //自动选择打开系统管理的页面菜单
	c.LayoutSections["LayoutScript"] = "user/userPage-js.html"
	c.TplName = "user/userPage.html"

}

type UserController struct {
	auth.LoginRequireController
}

func (c *UserController) List() {

	Draw, _ := c.GetInt("draw")
	Start, _ := c.GetInt("start")
	Length, _ := c.GetInt("length")
	ColName := c.GetString("colName")
	ColSort := c.GetString("colSort")
	SearchValue := c.GetString("searchvalue")
	//fmt.Println(Draw, Start, Length, ColName, ColSort, SearchValue)
	Total, TotalFilter, users := models.DefautlUserManager.GetUsers(Draw, Start, Length, ColName, ColSort, SearchValue)

	json := map[string]interface{}{
		"code":            200,
		"text":            "获取成功",
		"recordsTotal":    Total,
		"recordsFiltered": TotalFilter,
		"result":          users,
	}
	c.Data["json"] = json
	c.ServeJSON()

}

func (c *UserController) Create() {

	if c.Ctx.Input.IsPost() {

		valid := &validation.Validation{}
		users := &models.User{}
		c.ParseForm(users)

		if c.GetString("birthday") != "" {
			users.Birthday = utils.TimeParse(c.GetString("birthday")) //将字符串解析成时间格式*time.Time
		}

		if ok, err := valid.Valid(users); err == nil && ok {
			users.SetPassword(users.Password)
			if err := models.DefautlUserManager.Create(users); err == nil {
				c.Data["json"] = map[string]interface{}{
					"code": 200,
					"text": "用户" + `"` + users.Name + `"` + "创建成功！",
				}
			}

		} else {
			c.Data["json"] = map[string]interface{}{
				"code":   400,
				"result": valid.Errors,
			}
			fmt.Println()

		}
		c.ServeJSON()
	}

	c.TplName = "user/Create.html"

}

func (c *UserController) Modify() {
	pk, _ := c.GetInt("pk")
	getUser := models.DefautlUserManager.GetUserById(pk)

	if c.Ctx.Input.IsPost() {
		json := map[string]interface{}{
			"code":   400,
			"text":   "",
			"result": nil,
		}

		users := &models.User{}
		valid := &validation.Validation{}
		c.ParseForm(users)
		if c.GetString("birthday") != "" {
			users.Birthday = utils.TimeParse(utils.TimeStringSplit(c.GetString("birthday"))) //先将字串时间年月日分离出来，再将分离出来字符串解析成时间格式*time.Time
		}

		//fmt.Println(c.GetString("birthday"))
		//fmt.Println(users.Birthday)
		//fmt.Println(users)
		if ok, err := valid.Valid(users); err == nil && ok {

			models.DefautlUserManager.Modify(users)

			json = map[string]interface{}{
				"code":   200,
				"text":   "用户" + `"` + users.Name + `"` + "修改成功！",
				"result": nil,
			}

		} else if valid.ErrorMap()["NotChanged"] != nil {
			json = map[string]interface{}{
				"code":   304,
				"text":   "未修改！",
				"result": valid.Errors,
			}
		} else {
			json = map[string]interface{}{
				"code":   400,
				"text":   "修改失败！",
				"result": valid.Errors,
			}

		}
		c.Data["json"] = json
		c.ServeJSON()
	}
	c.Data["userModify"] = getUser
	c.TplName = "user/Modify.html"

}

func (c *UserController) Delete() {

	if c.Ctx.Input.IsPost() {
		pk, _ := c.GetInt("pk")

		getUser := models.DefautlUserManager.GetUserById(pk)

		if pk == c.GetSession("userID") { //判断进行锁定、解锁的用户是否为当前用户

			c.Data["json"] = map[string]interface{}{
				"code": 403,
				"text": "操作不被允许",
			}

		} else {
			models.DefautlUserManager.Delete(getUser)
			c.Data["json"] = map[string]interface{}{
				"code": 200,
				"text": "用户" + `"` + getUser.Name + `"` + "已删除！",
			}
		}

		c.ServeJSON()

	} else {

		c.TplName = "layout/GetXsrf.html"

	}

}

func (c *UserController) Lock() {

	//fmt.Println(c.User.Name),可直接打印LoginRequreController基础控制器的c.Data["user"]= users
	if c.Ctx.Input.IsPost() {

		pk, _ := c.GetInt("pk")
		getUser := models.DefautlUserManager.GetUserById(pk)

		if pk == c.GetSession("userID") { //判断进行锁定、解锁的用户是否为当前用户

			c.Data["json"] = map[string]interface{}{
				"code": 403,
				"text": "操作不被允许",
			}

		} else {

			if getUser.Status == 0 {
				c.Data["json"] = map[string]interface{}{
					"code": 200,
					"text": "用户" + `"` + getUser.Name + `"` + "锁定！",
				}
				models.DefautlUserManager.Lock(getUser)

			} else if getUser.Status == 1 {
				c.Data["json"] = map[string]interface{}{
					"code": 403,
					"text": "用户" + `"` + getUser.Name + `"` + "早已是锁定状态，无需再次锁定！",
				}
			}

		}

		c.ServeJSON()
	} else {
		c.TplName = "layout/GetXsrf.html"

	}

}

func (c *UserController) UnLock() {

	if c.Ctx.Input.IsPost() {
		pk, _ := c.GetInt("pk")
		getUser := models.DefautlUserManager.GetUserById(pk)

		if pk == c.GetSession("userID") { //判断进行锁定、解锁的用户是否为当前用户

			c.Data["json"] = map[string]interface{}{
				"code": 403,
				"text": "操作不被允许",
			}

		} else {

			if getUser.Status == 1 {
				c.Data["json"] = map[string]interface{}{
					"code": 200,
					"text": "用户" + `"` + getUser.Name + `"` + "已解锁！",
				}
				models.DefautlUserManager.Lock(getUser)

			} else if getUser.Status == 0 {
				c.Data["json"] = map[string]interface{}{
					"code": 403,
					"text": "用户" + `"` + getUser.Name + `"` + "早已是解锁状态，无需再解锁！",
				}
			}

		}

		c.ServeJSON()

	} else {
		c.TplName = "layout/GetXsrf.html"

	}

}

//用户token创建和更新请求接收处理的控制器方法
func (c *UserController) Token() {

	pk, _ := c.GetInt("pk")

	if c.Ctx.Input.IsPost() {

		json := map[string]interface{}{
			"code": 400,
			"text": "",
		}

		if rel := models.DefaultTokenManager.GencertToken(models.DefautlUserManager.GetUserById(pk)); rel == "inserted" {
			json = map[string]interface{}{
				"code": 200,
				"text": "token创建成功",
			}

		} else if rel == "updated" {
			json = map[string]interface{}{
				"code": 200,
				"text": "token更新成功",
			}
		} else {
			json = map[string]interface{}{
				"code": 403,
				"text": "token操作失败",
			}

		}

		c.Data["json"] = json
		c.ServeJSON()

	} else {

		users := models.DefautlUserManager.GetUserById(pk)
		//fmt.Printf("%#v\n", users)
		c.Data["usersToken"] = users
		//fmt.Println(template.HTML(c.XSRFFormHTML()))
		c.TplName = "user/token.html"

	}

}

//创建一个api测试的控制器方法，开启了xsrf功能后，post、put、delete等方法请求将被拒绝，只能用get、head
//curl -k -XGET https://www.xiaojimi.cn/user/api?name=rose  -H "Authentication: Token" -H "AccessKey: " -H "SecretKey: "
func (c *UserController) Api() { //api测试
	json := map[string]interface{}{}

	if c.GetString("name") == "rose" {

		json = map[string]interface{}{
			"code": 200,
			"text": "your name is rose，测试成功！",
		}
	} else if c.GetString("name") == "jack" {

		json = map[string]interface{}{
			"code": 200,
			"text": "your name is jack，测试成功！",
		}
	} else {
		json = map[string]interface{}{
			"code": 200,
			"text": "who you are?",
		}

	}

	c.Data["json"] = json
	c.ServeJSON()
}
