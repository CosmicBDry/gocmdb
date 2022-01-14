package cloud

import (
	"fmt"

	"github.com/CosmicBDry/gocmdb/server/cloud"
	"github.com/CosmicBDry/gocmdb/server/controllers/auth"
	"github.com/CosmicBDry/gocmdb/server/controllers/layout"
	"github.com/CosmicBDry/gocmdb/server/forms"
	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/astaxie/beego/validation"
)

type CloudPlatformPageController struct {
	layout.LayoutController
}

func (c *CloudPlatformPageController) Index() {

	if c.User.Id == 1 && c.User.Name == "admin" {
		c.Data["menu"] = "CloudPlatformManager"
		c.Data["expand"] = "CloudManager"
		c.LayoutSections["LayoutScript"] = "cloud/CloudPlatformPage/index-js.html"
		c.TplName = "cloud/CloudPlatformPage/index.html"
	}

}

type CloudPlatformController struct {
	auth.LoginRequireController
}

func (c *CloudPlatformController) List() {
	Draw, _ := c.GetInt64("draw")
	Start, _ := c.GetInt64("start")
	Length, _ := c.GetInt64("length")
	ColName := c.GetString("colName")
	ColSort := c.GetString("colSort")
	SearchValue := c.GetString("searchvalue")
	PlatformSearch := c.GetString("querySearch")

	//fmt.Println(Draw, Start, Length, ColName, ColSort, SearchValue,PlatformSearch)
	Total, TotalFilter, result := models.DefautlCloudPlatformManager.QueryList(Draw, Start, Length, ColName, ColSort, SearchValue, PlatformSearch)

	json := map[string]interface{}{
		"code":            200,
		"text":            "获取成功",
		"recordsTotal":    Total,
		"recordsFiltered": TotalFilter,
		"result":          result,
	}
	c.Data["json"] = json
	c.ServeJSON()
}

func (c *CloudPlatformController) Create() {

	if c.Ctx.Input.IsPost() {

		json := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		valid := &validation.Validation{}
		var CloudForm forms.CloudPlatformForm
		c.ParseForm(&CloudForm)
		if ok, err := valid.Valid(&CloudForm); err == nil && ok {
			name := CloudForm.Name
			types := CloudForm.Types
			addr := CloudForm.Addr
			region := CloudForm.Region
			accesskey := CloudForm.AccessKey
			secretkey := CloudForm.SecretKey
			remark := CloudForm.Remark

			if ok := models.DefautlCloudPlatformManager.Create(name, types, addr, region, accesskey, secretkey, remark); ok {
				cloudplatform := models.DefautlCloudPlatformManager.GetbyName(name)
				plugin := cloud.DefaultCloudManager.Plugins[types]
				plugin.Init(addr, region, accesskey, secretkey)
				Instances := plugin.GetInstances()

				for _, instance := range Instances {
					if err := models.DefaultVirtualMachineManager.SyncVmByCloudPlatform(cloudplatform, instance); err == nil {
						json = map[string]interface{}{
							"code": 200,
							"text": "创建成功",
						}
					} else {
						fmt.Println(err)
						json = map[string]interface{}{
							"code": 400,
							"text": "虚拟机实例同步失败，请检查数据库连接是否异常",
						}

					}

				}

			} else {
				json = map[string]interface{}{
					"code": 400,
					"text": "创建失败",
				}
			}

		} else {

			json = map[string]interface{}{
				"code":   400,
				"text":   "创建失败",
				"result": valid.Errors,
			}

		}
		c.Data["json"] = json
		c.ServeJSON()

	} else {

		c.Data["SelectType"] = cloud.DefaultCloudManager.Plugins
		c.TplName = "cloud/CloudPlatform/create.html"
	}
}

func (c *CloudPlatformController) Modify() {
	pk, _ := c.GetInt64("pk")
	if c.Ctx.Input.IsPost() {
		platform := forms.CloudPlatformForm{}
		c.ParseForm(&platform)
		pk, _ := c.GetInt64("Fromid")
		name := platform.Name
		types := platform.Types
		addr := platform.Addr
		region := platform.Region
		accesskey := platform.AccessKey
		secretkey := platform.SecretKey
		remark := platform.Remark
		models.DefautlCloudPlatformManager.Modify(pk, name, types, addr, region, accesskey, secretkey, remark)
		json := map[string]interface{}{
			"code": 200,
			"text": "修改成功",
		}
		c.Data["json"] = json
		c.ServeJSON()

	} else {
		cloudPlatform := models.DefautlCloudPlatformManager.GetbyId(pk)
		c.Data["CloudPlatform"] = cloudPlatform
		c.Data["SelectType"] = cloud.DefaultCloudManager.Plugins
		c.TplName = "cloud/CloudPlatform/modify.html"

	}

}

func (c *CloudPlatformController) Disable() {

	if c.Ctx.Input.IsPost() {

		json := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		pk, _ := c.GetInt64("pk")

		if status := models.DefautlCloudPlatformManager.GetbyId(pk); status.Status == 0 {
			models.DefautlCloudPlatformManager.Disable(pk)
			json = map[string]interface{}{
				"code": 200,
				"text": "禁用成功",
			}

		} else if status.Status == 1 {
			json = map[string]interface{}{
				"code": 403,
				"text": "已是禁用状态，无需重复禁用",
			}

		}

		c.Data["json"] = json
		c.ServeJSON()

	} else {

		c.TplName = "cloud/CloudPlatform/Getxsrf.html"

	}

}

func (c *CloudPlatformController) Enable() {
	if c.Ctx.Input.IsPost() {

		json := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		pk, _ := c.GetInt64("pk")

		if status := models.DefautlCloudPlatformManager.GetbyId(pk); status.Status == 1 {
			models.DefautlCloudPlatformManager.Enable(pk)
			json = map[string]interface{}{
				"code": 200,
				"text": "启用成功",
			}

		} else if status.Status == 0 {
			json = map[string]interface{}{
				"code": 403,
				"text": "已是启用状态，无需重复启用",
			}

		}

		c.Data["json"] = json
		c.ServeJSON()

	} else {

		c.TplName = "cloud/CloudPlatform/Getxsrf.html"

	}

}

func (c *CloudPlatformController) Delete() {
	if c.Ctx.Input.IsPost() {
		json := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		pk, _ := c.GetInt64("pk")
		if result := models.DefautlCloudPlatformManager.Delete(pk); result {
			json = map[string]interface{}{
				"code": 200,
				"text": "删除成功",
			}

		} else if result == false {
			json = map[string]interface{}{
				"code": 400,
				"text": "早已删除或数据不存在",
			}
		}

		c.Data["json"] = json
		c.ServeJSON()

	} else {

		c.TplName = "cloud/CloudPlatform/Getxsrf.html"

	}

}

type VirtualMachinePageController struct {
	layout.LayoutController
}

func (c *VirtualMachinePageController) Index() {

	c.Data["menu"] = "VirtualMachineManager"
	c.Data["expand"] = "CloudManager"
	c.LayoutSections["LayoutScript"] = "cloud/VirtualMachinePage/index-js.html"
	c.TplName = "cloud/VirtualMachinePage/index.html"
}

func (c *VirtualMachinePageController) List() {

	Draw, _ := c.GetInt64("draw")
	Start, _ := c.GetInt64("start")
	Length, _ := c.GetInt64("length")
	ColName := c.GetString("colName")
	ColSort := c.GetString("colSort")
	SearchValue := c.GetString("searchvalue")
	platform, _ := c.GetInt64("platform")

	json := map[string]interface{}{
		"code": 400,
		"text": "",
	}

	Total, TotalFilter, result := models.DefaultVirtualMachineManager.QueryList(Draw, Start, Length, platform, ColName, ColSort, SearchValue)

	json = map[string]interface{}{
		"code":            200,
		"text":            "获取成功",
		"recordsTotal":    Total,
		"recordsFiltered": TotalFilter,
		"result":          result,
	}

	c.Data["json"] = json
	c.ServeJSON()

}

type VirtualMachineController struct {
	auth.LoginRequireController
}

func (v *VirtualMachineController) Delete() {

	if v.Ctx.Input.IsPost() {
		json := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		pk, _ := v.GetInt64("pk")

		if err := models.DefaultVirtualMachineManager.Delete(pk); err == nil {
			json = map[string]interface{}{
				"code": 200,
				"text": "删除成功",
			}
		} else {
			json = map[string]interface{}{
				"code": 400,
				"text": "删除失败，检查数据库连接是否正常",
			}
		}
		v.Data["json"] = json
		v.ServeJSON()
	} else {
		v.TplName = "cloud/VirtualMachine/Getxsrf.html"
	}

}

func (v *VirtualMachineController) Stop() {

	if v.Ctx.Input.IsPost() {
		pk, _ := v.GetInt64("pk")
		json := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		machine := models.DefaultVirtualMachineManager.GetById(pk)
		addr := machine.CloudPlatform.Addr
		region := machine.CloudPlatform.Region
		accesskey := machine.CloudPlatform.AccessKey
		secretkey := machine.CloudPlatform.SecretKey
		instanceId := machine.InstanceId
		plugin := cloud.DefaultCloudManager.Plugins[machine.CloudPlatform.Types]
		plugin.Init(addr, region, accesskey, secretkey)
		if err := plugin.StopInstance(instanceId); err == nil {

			json = map[string]interface{}{
				"code": 200,
				"text": "实例关闭请求已发送",
			}
		} else if err != nil {
			json = map[string]interface{}{
				"code": 400,
				"text": "实例关闭请求发送失败，请检查云平台连接是否正常或api认证是否正确",
			}
		}

		v.Data["json"] = json
		v.ServeJSON()

	} else {
		v.TplName = "cloud/VirtualMachine/Getxsrf.html"
	}

}

func (v *VirtualMachineController) Start() {

	if v.Ctx.Input.IsPost() {
		pk, _ := v.GetInt64("pk")
		json := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		machine := models.DefaultVirtualMachineManager.GetById(pk)
		addr := machine.CloudPlatform.Addr
		region := machine.CloudPlatform.Region
		accesskey := machine.CloudPlatform.AccessKey
		secretkey := machine.CloudPlatform.SecretKey
		instanceId := machine.InstanceId
		plugin := cloud.DefaultCloudManager.Plugins[machine.CloudPlatform.Types]
		plugin.Init(addr, region, accesskey, secretkey)
		if err := plugin.StartInstance(instanceId); err == nil {

			json = map[string]interface{}{
				"code": 200,
				"text": "实例启动请求已发送",
			}
		} else if err != nil {
			json = map[string]interface{}{
				"code": 400,
				"text": "实例启动请求发送发送失败，请检查云平台连接是否正常或api认证是否正确",
			}
		}

		v.Data["json"] = json
		v.ServeJSON()

	} else {
		v.TplName = "cloud/VirtualMachine/Getxsrf.html"
	}

}

func (v *VirtualMachineController) Reboot() {

	if v.Ctx.Input.IsPost() {
		pk, _ := v.GetInt64("pk")
		json := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		machine := models.DefaultVirtualMachineManager.GetById(pk)
		addr := machine.CloudPlatform.Addr
		region := machine.CloudPlatform.Region
		accesskey := machine.CloudPlatform.AccessKey
		secretkey := machine.CloudPlatform.SecretKey
		instanceId := machine.InstanceId
		plugin := cloud.DefaultCloudManager.Plugins[machine.CloudPlatform.Types]
		plugin.Init(addr, region, accesskey, secretkey)
		if err := plugin.RestartInstance(instanceId); err == nil {

			json = map[string]interface{}{
				"code": 200,
				"text": "实例重启请求已发送",
			}
		} else if err != nil {
			json = map[string]interface{}{
				"code": 400,
				"text": "实例重启请求发送发送失败，请检查云平台连接是否正常或api认证是否正确",
			}
		}

		v.Data["json"] = json
		v.ServeJSON()

	} else {
		v.TplName = "cloud/VirtualMachine/Getxsrf.html"
	}

}

func (v *VirtualMachineController) Terminate() {

	if v.Ctx.Input.IsPost() {
		pk, _ := v.GetInt64("pk")
		json := map[string]interface{}{
			"code": 400,
			"text": "",
		}
		machine := models.DefaultVirtualMachineManager.GetById(pk)
		addr := machine.CloudPlatform.Addr
		region := machine.CloudPlatform.Region
		accesskey := machine.CloudPlatform.AccessKey
		secretkey := machine.CloudPlatform.SecretKey
		instanceId := machine.InstanceId
		plugin := cloud.DefaultCloudManager.Plugins[machine.CloudPlatform.Types]
		plugin.Init(addr, region, accesskey, secretkey)
		if err := plugin.TerminateInstance(instanceId); err == nil {

			json = map[string]interface{}{
				"code": 200,
				"text": "实例销毁请求已发送",
			}
		} else if err != nil {
			json = map[string]interface{}{
				"code": 400,
				"text": "实例销毁请求发送发送失败，请检查云平台连接是否正常或api认证是否正确",
			}
		}

		v.Data["json"] = json
		v.ServeJSON()

	} else {
		v.TplName = "cloud/VirtualMachine/Getxsrf.html"
	}

}
