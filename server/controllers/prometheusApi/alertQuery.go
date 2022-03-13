package prometheusApi

import (
	"github.com/CosmicBDry/gocmdb/server/controllers/layout"
	"github.com/CosmicBDry/gocmdb/server/models"
)

type AlertQueryController struct {
	layout.LayoutController
}

func (c *AlertQueryController) Index() {

	c.Data["expand"] = "SystemManager"
	c.Data["menu"] = "AlertManager"
	c.LayoutSections["LayoutScript"] = "alert/index-js.html"
	c.TplName = "alert/index.html"
}

func (c *AlertQueryController) List() {
	Draw, _ := c.GetInt64("draw")
	StartPos, _ := c.GetInt64("start")
	Length, _ := c.GetInt64("length")
	ColName := c.GetString("colName")
	ColSort := c.GetString("colSort")
	SearchValue := c.GetString("searchvalue")

	Total, TotalFilter, Results := models.GetList(StartPos, Length, ColName, ColSort, SearchValue)
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

func (c *AlertQueryController) Delete() {

	json := map[string]interface{}{
		"code": 400,
		"text": "",
	}

	if c.Ctx.Input.IsPost() {

		pk, _ := c.GetInt64("pk")

		if err := models.Delete(pk); err == nil {
			json = map[string]interface{}{
				"code": 200,
				"text": "告警删除成功！",
			}

		} else {
			json = map[string]interface{}{
				"code": 405,
				"text": "告警删除失败，请检查数据库连接是否异常！",
			}

		}
	} else {
		json = map[string]interface{}{
			"code": 400,
			"text": "不支持该请求方法，只支持POST请求",
		}

	}

	c.Data["json"] = json
	c.ServeJSON()

}
