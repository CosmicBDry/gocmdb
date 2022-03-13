package gitlab

import (
	"encoding/json"
	"net/http"

	//"github.com/CosmicBDry/gocmdb/server/controllers/base"
	"github.com/CosmicBDry/gocmdb/server/buildservice"
	"github.com/CosmicBDry/gocmdb/server/controllers/auth"
	"github.com/CosmicBDry/gocmdb/server/controllers/layout"
	"github.com/CosmicBDry/gocmdb/server/forms"
	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/astaxie/beego/validation"
)

type GitLabPageController struct {
	layout.LayoutController
}

//gitab项目page操作页面--------------------->
func (c *GitLabPageController) Index() {
	c.Data["menu"] = "ProdRelease"
	c.Data["expand"] = "ReleaseManager"
	c.LayoutSections["LayoutScript"] = "gitlab/index-js.html"
	c.TplName = "gitlab/index.html"

}

type GitLabHookController struct {
	//base.BaseController
	auth.LoginRequireController
}

//接收gitlab的webhook请求--------------------->
func (c *GitLabHookController) WebHook() {

	jsonResponse := map[string]interface{}{
		"code": 400,
		"text": "",
	}
	if c.Ctx.Input.IsPost() {
		c.Ctx.Input.CopyBody(1024 * 1024)
		//fmt.Println(string(c.Ctx.Input.RequestBody))
		hookForm := &forms.ProjectHook{}
		json.Unmarshal(c.Ctx.Input.RequestBody, hookForm)
		gitproject := hookForm.GitFormToModel()
		//fmt.Println("X-Gitlab-Token: ", c.Ctx.Input.Header("X-Gitlab-Token"))
		//_, ok := c.Ctx.Request.Header["X-Gitlab-Token"]

		if err := gitproject.GitCreate(); err != nil {
			//自定义服务端的http的响应状态码，http.StatusBadRequest为400状态码，也可以直接写上400状态码
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			jsonResponse = map[string]interface{}{
				"code": 400,
				"text": err.Error(),
			}
		} else {
			jsonResponse = map[string]interface{}{
				"code": 200,
				"text": "请求成功！",
			}
		}

	} else {
		//自定义服务端的http的响应状态码，http.StatusMethodNotAllowed为405，即请求的方法不被允许
		c.Ctx.Output.SetStatus(http.StatusMethodNotAllowed)
		jsonResponse = map[string]interface{}{
			"code": 405,
			"text": "请求的方法不被允许",
		}
	}
	c.Data["json"] = jsonResponse
	c.ServeJSON()
}

func (c *GitLabHookController) List() {
	Draw, _ := c.GetInt64("draw")
	StartPos, _ := c.GetInt64("start")
	Length, _ := c.GetInt64("length")
	ColName := c.GetString("colName")
	ColSort := c.GetString("colSort")
	SearchValue := c.GetString("searchvalue")
	Total, TotalFilter, Results := models.GitLabGetList(StartPos, Length, ColName, ColSort, SearchValue)
	//	fmt.Println(Draw, StartPos, Length, ColName, ColSort, SearchValue)
	json := map[string]interface{}{
		"code":            200,
		"text":            "响应成功",
		"draw":            Draw,
		"results":         Results,
		"recordsTotal":    Total,
		"recordsFiltered": TotalFilter,
	}
	c.Data["json"] = json
	c.ServeJSON()
}

//git项目版本、提交者、项目地址等信息页面查询--------------------->
func (c *GitLabHookController) ProjectQuery() {
	pk, _ := c.GetInt64("pk")
	gitproject, _ := models.GitGetById(pk)
	c.Data["gitProject"] = gitproject
	c.TplName = "gitlab/projectQuery.html"
}

//项目构建发布配置查询------------------------------------->
func (c *GitLabHookController) ReleaseConfigQuery() {
	pk, _ := c.GetInt64("pk")
	gitproject, _ := models.GitGetById(pk)
	c.Data["gitProject"] = gitproject
	c.TplName = "gitlab/releaseConfigQuery.html"
}

//项目构建发布前的配置------------------------------------->
func (c *GitLabHookController) ReleaseConfigModify() {
	if c.Ctx.Input.IsPost() {
		releaseform := &forms.ReleaseConfigForm{}
		//客户端发送的form中字符串数字可以直接解析到结构提对应的bool类型中
		c.ParseForm(releaseform)
		jsonResponse := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		valid := &validation.Validation{}

		if ok, err := valid.Valid(releaseform); err == nil && ok {
			//autodeploy, _ := strconv.ParseBool(releaseform.AutoDeploy)
			if err := models.ReleaseConfigModify(int64(releaseform.Id), releaseform.ReleaseMachine, releaseform.BackendHost, releaseform.PackageFile, releaseform.AutoDeploy); err == nil {
				jsonResponse = map[string]interface{}{
					"code": 200,
					"text": "配置更新成功！",
				}
			} else if err.Error() == "未发生更改,无需更新！" {
				jsonResponse = map[string]interface{}{
					"code": 304,
					"text": err.Error(),
				}
			} else {
				jsonResponse = map[string]interface{}{
					"code": 400,
					"text": err.Error(),
				}
			}
		} else {
			jsonResponse = map[string]interface{}{
				"code":   400,
				"text":   "数据验证错误",
				"result": valid.Errors,
			}

		}

		c.Data["json"] = jsonResponse
		c.ServeJSON()
	}
	pk, _ := c.GetInt64("pk")
	gitproject, _ := models.GitGetById(pk)
	c.Data["gitProject"] = gitproject
	c.TplName = "gitlab/releaseConfigModify.html"
}

//项目编译构建------------------------------------->
func (c *GitLabHookController) Build() {

	jsonResponse := map[string]interface{}{
		"code": 400,
		"text": "",
	}

	if c.Ctx.Input.IsPost() {

		pk, _ := c.GetInt64("pk")

		gitproject, _ := models.GitGetById(pk)
		if err := buildservice.DefaultBuildService.Build(gitproject); err != nil {
			//c.Ctx.Output.SetStatus(http.StatusBadRequest)
			jsonResponse = map[string]interface{}{
				"code": 400,
				"text": gitproject.Name + "构建失败,详情请查看构建记录！",
			}
		} else {
			jsonResponse = map[string]interface{}{
				"code": 200,
				"text": gitproject.Name + "已完成构建！",
			}
		}
		c.Data["json"] = jsonResponse
		c.ServeJSON()

	} else if c.Ctx.Input.IsGet() {
		pk, _ := c.GetInt64("pk")
		gitproject, _ := models.GitGetById(pk)
		buildLog := buildservice.GetBuildLog(gitproject)
		c.Data["build"] = buildLog
		c.TplName = "gitlab/buildLog.html"
	}

}

//项目部署待开发------------------------------------->
func (c *GitLabHookController) Deploy() {}

//项目回滚待开发------------------------------------->
func (c *GitLabHookController) Rollout() {}
