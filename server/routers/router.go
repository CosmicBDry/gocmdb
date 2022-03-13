package routers

import (
	"github.com/astaxie/beego"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/CosmicBDry/gocmdb/server/controllers/auth"
	"github.com/CosmicBDry/gocmdb/server/controllers/cloud"
	"github.com/CosmicBDry/gocmdb/server/controllers/gitlab"
	"github.com/CosmicBDry/gocmdb/server/controllers/kubernetes"
	"github.com/CosmicBDry/gocmdb/server/controllers/prometheusApi"
	"github.com/CosmicBDry/gocmdb/server/controllers/user"
	"github.com/CosmicBDry/gocmdb/server/filters"
)

func Register() {

	beego.AutoRouter(&auth.AuthController{})

	beego.AutoRouter(&user.UserPageController{})

	beego.AutoRouter(&user.UserController{})
	beego.AutoRouter(&cloud.CloudPlatformPageController{})
	beego.AutoRouter(&cloud.CloudPlatformController{})
	beego.AutoRouter(&cloud.VirtualMachinePageController{})
	beego.AutoRouter(&cloud.VirtualMachineController{})
	//过滤的插入
	beego.InsertFilter("/*", beego.BeforeExec, filters.BeforeExec)
	//控制器执行完后一定有输出，则不会执行后续过滤器操作beego.AfterExec，因此将第四参数改为false允许后续操作
	beego.InsertFilter("/*", beego.AfterExec, filters.AfterExec, false)
	//自定义一个handler处理器，提供于prometheus指标采集
	beego.Handler("/metrics", promhttp.Handler())

	//beego接收altermanager告警的api控制器路由注册
	beego.AutoRouter(&prometheusApi.PrometheusController{})
	beego.AutoRouter(&prometheusApi.AlertQueryController{})

	beego.AutoRouter(&prometheusApi.AgentController{})
	beego.AutoRouter(&prometheusApi.AgentPageController{})
	beego.AutoRouter(&prometheusApi.AgentOperationController{})

	//gitlab的webhook触发获取
	beego.AutoRouter(&gitlab.GitLabHookController{})
	beego.AutoRouter(&gitlab.GitLabPageController{})

	//kubernetes集群
	beego.AutoRouter(&kubernetes.K8sPageController{})

	beego.AutoRouter(&kubernetes.K8sClusterController{})

}
