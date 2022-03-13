package kubernetes

import (
	"github.com/CosmicBDry/gocmdb/server/controllers/auth"
	"github.com/CosmicBDry/gocmdb/server/controllers/layout"
	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/CosmicBDry/gocmdb/server/forms"
	"github.com/astaxie/beego"
	"net/http"
	"strings"
	//"errors"
	"fmt"
)

type K8sPageController struct {
	layout.LayoutController
}

func (c *K8sPageController) DeploymentIndex() {

	c.Data["expand"] = "KubernetesCluster"
	c.Data["menu"] = "DeploymentConfig"
	c.LayoutSections["LayoutScript"] = "kubernetes/deployments/index-js.html"
	c.TplName = "kubernetes/deployments/index.html"

}

type K8sClusterController struct {
	auth.LoginRequireController
}

func (c *K8sClusterController) List() {
	Draw, _ := c.GetInt64("draw")
	StartPos, _ := c.GetInt64("start")
	Length, _ := c.GetInt64("length")

	jsonResponse := map[string]interface{}{
		"code": 400,
		"text": "",
	}
	deployments, err := models.DeploymentList(beego.AppConfig.String("KubeConfig"))
	if err != nil {

		jsonResponse = map[string]interface{}{
			"code": 400,
			"text": "deployment资源获取失败！",
		}
	} else {

		Total, TotalFilter, Results := models.DeploymentPageList(StartPos, Length, deployments)
		if len(Results) > 0 {
			jsonResponse = map[string]interface{}{
				"code":            200,
				"text":            "deployment资源获取成功！",
				"Draw":            Draw,
				"results":         Results,
				"recordsTotal":    Total,
				"recordsFiltered": TotalFilter,
			}
		} else {
			jsonResponse = map[string]interface{}{
				"code": 305,
				"text": "deployment资源获取异常！",
				"Draw":            Draw,
				"results":         Results,
				"recordsTotal":    Total,
				"recordsFiltered": TotalFilter,
			}
		}

	}

	c.Data["json"] = jsonResponse
	c.ServeJSON()

}


func(c *K8sClusterController) DeploymentCreate(){

	if c.Ctx.Input.IsPost(){
		jsonResponse := map[string]interface{}{
			"code":400,
			"text":"",
		}
		deploymentForm := forms.NewDeploymentForm()

		c.ParseForm(deploymentForm)
		deploymentInstance :=deploymentForm.DeploymentFormToModel()

		err :=models.DeploymentInstanceCreate(beego.AppConfig.String("KubeConfig"),deploymentInstance)
		if err !=nil{
			jsonResponse = map[string]interface{}{
				"code":400,
				"text":"创建失败："+err.Error(),
			}
		} else{
			jsonResponse = map[string]interface{}{
				"code":200,
				"text":"创建成功!",
			}
		}
		
		c.Data["json"]=jsonResponse
		c.ServeJSON()
	}


	c.TplName="kubernetes/deployments/create.html"

}


func(c *K8sClusterController) DeploymentDelete(){
	responseJson :=map[string]interface{}{
		"code":400,
		"text":"",
	}

	if c.Ctx.Input.IsPost(){
			fmt.Println(c.GetString("pk"))
		
		pk :=strings.SplitN(c.GetString("pk"),"/",2)
		err :=models.DeploymentInstanceDelete(beego.AppConfig.String("KubeConfig"),pk[0],pk[1])
		if err == nil{
			c.Ctx.Output.SetStatus(200)
			responseJson =map[string]interface{}{
				"code":200,
				"text":"删除操作成功！",
			}
		}else{
			c.Ctx.Output.SetStatus(400)
			responseJson =map[string]interface{}{
				"code":400,
				"text":"删除操作失败: "+err.Error(),
			}
		}
	}else{
		c.Ctx.Output.SetStatus(http.StatusForbidden)//403
		responseJson =map[string]interface{}{
			"code":403,
			"text":"无效的请求方法",
		}	
	}

	c.Data["json"]=responseJson
	c.ServeJSON()

}

func( c *K8sClusterController) DeploymentModify(){
	
	 if c.Ctx.Input.IsPost(){
		responseJson := map[string]interface{}{
			 "code":400,
			 "text":"",
		 }
		deploymentForm := forms.NewDeploymentForm()
		c.ParseForm(deploymentForm)
		deploymentInstance :=deploymentForm.DeploymentFormToModel()
		err :=models.DeploymentInstanceModify(beego.AppConfig.String("KubeConfig"),deploymentInstance)
		if err == nil{
			c.Ctx.Output.SetStatus(http.StatusOK)//200
			responseJson=map[string]interface{}{
				"code":200,
				"text":"控制器更新成功！",
			}
		}else{
			c.Ctx.Output.SetStatus(http.StatusBadRequest)//404
			responseJson=map[string]interface{}{
				"code":400,
				"text":"控制器更新失败: "+err.Error(),
			}
		}

		c.Data["json"]=responseJson
		c.ServeJSON()
	 }
	pk :=strings.SplitN(c.GetString("pk"),"/",2)
	fmt.Println(pk)
	deploymentInstance,_:=models.DeploymentInstanceGet(beego.AppConfig.String("KubeConfig"),pk[0],pk[1])
	fmt.Printf("%#v\n",deploymentInstance)
	c.Data["GetDeployment"]= *deploymentInstance
	c.TplName="kubernetes/deployments/modify.html"

}