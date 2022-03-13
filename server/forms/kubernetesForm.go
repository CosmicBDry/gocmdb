package forms 

import(

	"github.com/CosmicBDry/gocmdb/server/models"
	
	"strings"
)

type  DeploymentForm struct{
	Name string `form:"name"`
	NameSpace string `form:"namespace"`
	HistoryVersionLimit int`form:"historyVersionLimit"`
	Replicas int`form:"replicas"`
	Labels string`form:"labels"`
	ContainerName string`form:"containerName"`
	ImageUrl string`form:"imageUrl"`
	ContainerPortName string`form:"containerPortName"`
	ContainerPort int`form:"containerPort"`
	HostPort int `form:"hostPort"`
}

func NewDeploymentForm() *DeploymentForm{
	return &DeploymentForm{}
}

func (f *DeploymentForm) DeploymentFormToModel() *models.DeploymentInstace {

	 LabelsMap :=make(map[string]string)
	 if f.Labels !=""{
		for _,label:=range strings.Split(f.Labels,","){
			label = strings.TrimSpace(label)
			result :=strings.SplitN(label,"=",2)
			LabelsMap[result[0]]=result[1]
		}
	 }
	 
	return &models.DeploymentInstace{
		Name: f.Name,
		NameSpace: f.NameSpace,
		RevisionHistoryLimit: int32(f.HistoryVersionLimit),
		Replicas: int32(f.Replicas),            
		Labels: LabelsMap,   			
		ContainerName: f.ContainerName,
		ImgUrl: f.ImageUrl,
		ContainerPortName: f.ContainerPortName,
		ContainerPort: int32(f.ContainerPort),   
		HostPort: int32(f.HostPort), 
	}

}





