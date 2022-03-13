package prometheusApi

import (
	"encoding/json"

	"github.com/CosmicBDry/gocmdb/server/controllers/base"
	"github.com/CosmicBDry/gocmdb/server/forms"
)

type PrometheusController struct {
	base.BaseController
}

func (c *PrometheusController) Alert() {
	//fmt.Println(string(c.Ctx.Input.CopyBody(1024 * 1024))) //CopyBody与RequestBody输出结果相同
	c.Ctx.Input.CopyBody(1024 * 1024)
	//fmt.Println(strings.Repeat("->", 36))
	alertform := forms.AlertForm{}

	json.Unmarshal(c.Ctx.Input.RequestBody, &alertform)
	json := map[string]interface{}{
		"code": 400,
		"text": "",
	}

	if c.Ctx.Input.IsPost() {
		if alertform.Alerts != nil {
			for _, alerts := range alertform.Alerts {

				if err := alerts.ToModel().CreateOrUpdate(); err != nil {
					//fmt.Println(err)
					json = map[string]interface{}{
						"code":  400,
						"text":  "告警请求操作失败",
						"error": err,
					}
				} else {
					json = map[string]interface{}{
						"code": 200,
						"text": "告警请求操作成功",
					}

				}

				c.Data["json"] = json
				c.ServeJSON()
			}
		} else {
			json = map[string]interface{}{
				"code": 405,
				"text": "告警请求不可为空",
			}
		}
	} else {
		json = map[string]interface{}{
			"code": 406,
			"text": "请通过POST方法请求",
		}
	}

	c.Data["json"] = json
	c.ServeJSON()

}
