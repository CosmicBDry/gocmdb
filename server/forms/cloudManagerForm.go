package forms

import (
	"strings"

	"github.com/CosmicBDry/gocmdb/server/cloud"
	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/astaxie/beego/validation"
)

type CloudPlatformForm struct {
	Name      string `form:"name"`
	Types     string `form:"types"`
	Addr      string `form:"addr"`
	Region    string `form:"region"`
	AccessKey string `form:"accesskey"`
	SecretKey string `form:"secretkey"`
	Remark    string `form:"remark"`
}

func (f *CloudPlatformForm) Valid(v *validation.Validation) {

	f.Name = strings.TrimSpace(f.Name)
	f.Types = strings.TrimSpace(f.Types)
	f.Addr = strings.TrimSpace(f.Addr)
	f.Region = strings.TrimSpace(f.Region)
	f.AccessKey = strings.TrimSpace(f.AccessKey)
	f.SecretKey = strings.TrimSpace(f.SecretKey)
	f.Remark = strings.TrimSpace(f.Remark)

	v.MinSize(f.Name, 2, "name.name").Message("名称最小2字符")
	v.MaxSize(f.Name, 32, "name.name").Message("名称最大32字符")
	v.MinSize(f.Region, 2, "region.region").Message("region域最小2个字符")
	v.MaxSize(f.Region, 64, "region.region").Message("region域最大64字符")
	v.MaxSize(f.SecretKey, 1024, "secretkey.secretkey").Message("最大1024字符")
	v.MaxSize(f.AccessKey, 1024, "accessKey.accesskey").Message("最大1024字符")
	v.MaxSize(f.Remark, 1024, "remark.remark").Message("最大1024字符")

	if f.Types != "aliyun" && f.Types != "tencentCloud" {
		v.SetError("types", "该类型未注册")
	} else if models.DefautlCloudPlatformManager.GetbyName(f.Name) != nil {
		v.SetError("name", "云平台名称早已存在，不可重复")
	} else {
		cloud.DefaultCloudManager.Plugins[f.Types].Init(f.Addr, f.Region, f.AccessKey, f.SecretKey)
		if err := cloud.DefaultCloudManager.Plugins[f.Types].TestConnect(); err != nil {
			v.SetError("error", "云平台连接失败，请检查请求token认证、区域等是否正确")
		}
	}
}
