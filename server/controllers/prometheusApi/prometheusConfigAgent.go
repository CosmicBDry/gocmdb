package prometheusApi

import (
	"encoding/json"
	"fmt"
	"time"

	"strconv"

	"github.com/CosmicBDry/gocmdb/server/controllers/auth"
	"github.com/CosmicBDry/gocmdb/server/controllers/base"
	"github.com/CosmicBDry/gocmdb/server/controllers/layout"
	"github.com/CosmicBDry/gocmdb/server/forms"
	"github.com/CosmicBDry/gocmdb/server/models"
	"gopkg.in/yaml.v2"
)

type AgentController struct {
	base.BaseController
}

func (c *AgentController) Register() {

	jsonResponse := map[string]interface{}{
		"code": 400,
		"text": "",
	}

	//fmt.Println("Input.CopyBody----->", c.Ctx.Input.CopyBody(1024*1024))
	c.Ctx.Input.CopyBody(1024 * 1024)
	agentForm := forms.DefaultAgentForm
	json.Unmarshal(c.Ctx.Input.RequestBody, agentForm)
	//fmt.Printf("%#v\n", agentForm)
	agent := agentForm.ToModel()
	//fmt.Printf("--------------------->%#v\n", agent)
	if err, ok := models.DefaultAgentManager.Register(agent); err != nil {
		jsonResponse = map[string]interface{}{
			"code":  400,
			"text":  "注册失败",
			"error": err,
		}
	} else if ok == "inserted" {
		jsonResponse = map[string]interface{}{
			"code": 200,
			"text": "注册成功",
		}

	} else if ok == "updated" {
		jsonResponse = map[string]interface{}{
			"code": 201,
			"text": "注册已更新",
		}
	} else if ok == "" {
		jsonResponse = map[string]interface{}{
			"code": 300,
			"text": "已注册，无需重复注册",
		}
	}
	c.Data["json"] = jsonResponse
	c.ServeJSON()

}

//心跳检测
func (c *AgentController) HeartBeat() {
	jsons := map[string]interface{}{
		"code": 400,
		"text": "",
	}
	if err := models.DefaultAgentManager.HeartBeat(c.GetString("uuid")); err != nil {
		jsons = map[string]interface{}{
			"code":  400,
			"text":  "心跳检测失败，确定是否注册或数据库连接异常？",
			"error": err.Error(),
		}
	} else {

		jsons = map[string]interface{}{
			"code": 200,
			"text": "心跳检测成功！",
		}
	}

	c.Data["json"] = jsons
	c.ServeJSON()

}

func (c *AgentController) SetConfig() {

}

func (c *AgentController) GetConfig() {
	jsons := map[string]interface{}{
		"code": 400,
		"text": "",
	}

	Config, Configversion, err := models.DefaultAgentManager.GetConfig(c.GetString("uuid"), c.GetString("configversion"))
	if err != nil {
		jsons = map[string]interface{}{
			"code":  400,
			"text":  "配置获取失败",
			"error": err.Error(),
		}

	} else if Config == "" {
		jsons = map[string]interface{}{
			"code": 300,
			"text": "本地配置版本已是最新",
		}

	} else {
		jsons = map[string]interface{}{
			"code":          200,
			"text":          "配置获取成功",
			"config":        Config,
			"configversion": Configversion,
		}

	}

	c.Data["json"] = jsons

	c.ServeJSON()

}

//prometheus代理的页面展示，提供配置、删除按钮操作
type AgentPageController struct {
	layout.LayoutController
}

func (c *AgentPageController) Index() {
	c.Data["expand"] = "SystemManager"
	c.Data["menu"] = "AgentManager"
	c.LayoutSections["LayoutScript"] = "agent/index-js.html"
	c.TplName = "agent/index.html"
}

func (c *AgentPageController) List() {
	Draw, _ := c.GetInt64("draw")
	StartPos, _ := c.GetInt64("start")
	Length, _ := c.GetInt64("length")
	ColName := c.GetString("colName")
	ColSort := c.GetString("colSort")
	SearchValue := c.GetString("searchvalue")

	Total, TotalFilter, Results := models.AgentList(StartPos, Length, ColName, ColSort, SearchValue)
	//	fmt.Println(Draw, StartPos, Length, ColName, ColSort, SearchValue)
	jsons := map[string]interface{}{
		"code":            200,
		"text":            "响应成功",
		"draw":            Draw,
		"results":         Results,
		"recordsTotal":    Total,
		"recordsFiltered": TotalFilter,
	}

	c.Data["json"] = jsons
	c.ServeJSON()
}

//prometheus的dailog代理配置编辑处理页面控制器
type AgentOperationController struct {
	auth.LoginRequireController
}

func (c *AgentOperationController) Modify() {
	jobs := models.Jobs{}
	if c.Ctx.Input.IsPost() {
		agentconfigForm := &forms.AgentConfigForm{}
		c.ParseForm(agentconfigForm)
		err := yaml.UnmarshalStrict([]byte(agentconfigForm.Config), &jobs)
		//fmt.Println(jobs.ScrapeConfigs[0].JobName, err)
		jsons := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		if err == nil {
			id := agentconfigForm.Id
			config := agentconfigForm.Config
			configversion := agentconfigForm.ConfigVersion
			agent := models.DefaultAgentManager.GetAgentById(id)
			if agent.Config != config {
				configversion = strconv.Itoa(int(time.Now().Unix()))
				err1 := models.DefaultAgentManager.SetConfig(id, config, configversion)
				if err1 == nil {
					jsons = map[string]interface{}{
						"code": 200,
						"text": "配置更新成功！",
					}
				} else {
					fmt.Println(err1)
					jsons = map[string]interface{}{
						"code": 400,
						"text": "配置更新失败，请检查数据库连接！",
					}
				}
			} else {
				jsons = map[string]interface{}{
					"code": 304,
					"text": "未发生修改！",
				}
			}
		} else {
			jsons = map[string]interface{}{
				"code": 400,
				"text": "无效的yaml格式！",
			}
		}
		c.Data["json"] = jsons
		c.ServeJSON()
	} else {
		pk, _ := c.GetInt64("pk")
		agent := models.DefaultAgentManager.GetAgentById(pk)
		c.Data["AgentForm"] = agent
		c.TplName = "agent/modify.html"
	}

}

/*func (c *AgentOperationController) Test() {

	c.TplName = "test/test.html"

}*/
