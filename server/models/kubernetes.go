package models

import (
	"context"
	"time"

	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubeAppV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
)


type DeploymentInstace struct {
	Name                 string            `json:"app_name" form:"name"`
	NameSpace            string            `json:"namespace" form:"namespace"`
	CreationTimestamp    time.Time         `json:"created_time"`
	RevisionHistoryLimit int32             `json:"revision_history_limit" form:"historyVersionLimit"`
	Replicas             int32             `json:"replicas" form:"replicas"`
	AvailableReplicas    int32             `json:"available_replicas"`
	Labels    			 map[string]string `json:"-" form:"labels"`
	ContainerName string `form:"containerName"`
	ImgUrl string `form:"imageUrl"`
	ContainerPortName    string       `json:"-"form:"containerPortName"`
	ContainerPort    int32 `form:"containerPort"`
	HostPort int32 `form:"hostPort"`
}

//通过kubeconfig配置文件生成一个集群操作客户端clientset（为结构体实例）
func ClienSet(kubeconfig string) (*kubernetes.Clientset,error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}
	return clientset,nil

}


//返回集群中所有的namespace结构体对象
func NamespaceList(kubeconfig string) (*coreV1.NamespaceList, error) {

	clientset,err:=ClienSet(kubeconfig)
	if err !=nil{
		return nil,err
	}
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return namespaceList, nil

}

//所有名称空间的deployment接口集合，每个名称空间都有自己独有的deployment接口
func DeploymentInterface(kubeconfig string) ([]interface{}, error) {

	var deploymentInterface []interface{}
	clientset,err:=ClienSet(kubeconfig)
	if err !=nil{
		return nil,err
	}

	namespaceList, err := NamespaceList(kubeconfig)
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaceList.Items {

		deploymentInterface = append(deploymentInterface, clientset.AppsV1().Deployments(namespace.Name))

	}

	return deploymentInterface, nil

}

//查询集群中所有deployment实例
//调用DeploymentInterface()返回的接口集合中接口的list方法，返回所有存在的deployment实例集合
func DeploymentList(kubeconfig string) ([]DeploymentInstace, error) {
	var deploymentInstances []DeploymentInstace
	DeploymentClients, err := DeploymentInterface(kubeconfig)
	if err != nil {
		return nil, err
	}

	for _, deploymentclient := range DeploymentClients {
		deploymentList, err := deploymentclient.(kubeAppV1.DeploymentInterface).List(context.TODO(), metaV1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, deployment := range deploymentList.Items {
			deploymentInstances = append(deploymentInstances, DeploymentInstace{
				Name:                 deployment.Name,
				NameSpace:            deployment.Namespace,
				CreationTimestamp:    deployment.CreationTimestamp.Time,
				RevisionHistoryLimit: *deployment.Spec.RevisionHistoryLimit,
				Replicas:             *deployment.Spec.Replicas,
				AvailableReplicas:    deployment.Status.AvailableReplicas,
			})
		}
	}
	return deploymentInstances, nil
}

//前端datatable参数的提供，返回数据总数、过滤总数、deployment实例列表
func DeploymentPageList(start, length int64, deploymentInstances []DeploymentInstace) (int64, int64, []DeploymentInstace) {
	var Total int64
	var filterTotal int64
	var filterDeployments []DeploymentInstace

	maxIndex := int64(len(deploymentInstances) - 1)
	Total = maxIndex + 1

	if len(deploymentInstances) > 0 {
		if maxIndex > start+length-1 {
			filterDeployments = deploymentInstances[start : start+length]

		} else if maxIndex <= start+length-1 {
			filterDeployments = deploymentInstances[start:]
		}

	} else {
		return 0, 0, []DeploymentInstace{}
	}
	filterTotal = int64(len(filterDeployments))

	return Total, filterTotal, filterDeployments

}

//Deployment实例创建
func DeploymentInstanceCreate(kubeconfig string,deployment *DeploymentInstace)error{

			deploymentInstance := &appV1.Deployment{
				TypeMeta: metaV1.TypeMeta{ //定义deployment资源的apiVersion、kind类型等
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				ObjectMeta: metaV1.ObjectMeta{ //定义控制器的metadata：名称、namespace、labels等
					Name:      deployment.Name,
					Namespace: deployment.NameSpace,
					Labels: deployment.Labels,
				},
				Spec: appV1.DeploymentSpec{ //定义deployment控制器的spec字段
					Replicas: &deployment.Replicas,
					Selector: &metaV1.LabelSelector{
									MatchLabels: deployment.Labels,
							  },
					RevisionHistoryLimit: &deployment.RevisionHistoryLimit,		
					Template: coreV1.PodTemplateSpec{
						ObjectMeta: metaV1.ObjectMeta{
							//Name: "golang-pod",
							Labels: deployment.Labels,
						},
						Spec: coreV1.PodSpec{
							Containers: []coreV1.Container{
								coreV1.Container{
									Name:            deployment.ContainerName,
									Image:           deployment.ImgUrl,
									ImagePullPolicy: coreV1.PullIfNotPresent,
									Ports: []coreV1.ContainerPort{
										coreV1.ContainerPort{
											Name:          deployment.ContainerPortName,
											ContainerPort:  deployment.ContainerPort,
											HostPort: deployment.HostPort,
											Protocol:      coreV1.ProtocolTCP,
										},
									},
								},
							},
						},
					},
				},
			}

			clientset,err:=ClienSet(kubeconfig)
				if err !=nil{
					return err
				}

			DeploymentClient:=clientset.AppsV1().Deployments(deployment.NameSpace)

			_, err = DeploymentClient.Create(context.TODO(), deploymentInstance, metaV1.CreateOptions{})
			if err !=nil{
				return err
			}
	return nil
}

func  DeploymentInstanceGet(kubeconfig,namespace,deploymentName string) (*DeploymentInstace,error){
	clientset,err:=ClienSet(kubeconfig)
	if err !=nil{
		return nil,err
	}
	DeploymentClient:=clientset.AppsV1().Deployments(namespace)
	deploymentInstance,err:=DeploymentClient.Get(context.TODO(),deploymentName,metaV1.GetOptions{})
	if err !=nil{
		return nil,err
	}

	
	return &DeploymentInstace{
		Name: deploymentInstance.Name,
		NameSpace: deploymentInstance.Namespace,
		CreationTimestamp:    deploymentInstance.CreationTimestamp.Time,
		RevisionHistoryLimit: *deploymentInstance.Spec.RevisionHistoryLimit,
		Replicas:   *deploymentInstance.Spec.Replicas,          
		AvailableReplicas: deploymentInstance.Status.AvailableReplicas,
		Labels: deploymentInstance.Spec.Selector.MatchLabels , 			 
		ContainerName: deploymentInstance.Spec.Template.Spec.Containers[0].Name,
		ImgUrl: deploymentInstance.Spec.Template.Spec.Containers[0].Image,
		ContainerPortName:  deploymentInstance.Spec.Template.Spec.Containers[0].Ports[0].Name  ,
		ContainerPort: deploymentInstance.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort ,
		HostPort: deploymentInstance.Spec.Template.Spec.Containers[0].Ports[0].HostPort,
	},nil

}


func DeploymentInstanceDelete(kubeconfig,namespace,deploymentName string)error{
	var graceTime int64 = 30
	clientset,err:=ClienSet(kubeconfig)
	if err !=nil{
		return err
	}
	DeploymentClient:=clientset.AppsV1().Deployments(namespace)
	err =DeploymentClient.Delete(context.TODO(),deploymentName,metaV1.DeleteOptions{
		GracePeriodSeconds: &graceTime,
	})
	if err !=nil{
		return err
	}
	return nil
}

func DeploymentInstanceModify(kubeconfig string,newDeploymentInstance *DeploymentInstace)error{
	
	
	clientset,err:=ClienSet(kubeconfig)
	if err !=nil{
		return err
	}


	DeploymentClient:=clientset.AppsV1().Deployments(newDeploymentInstance.NameSpace)
	deploymentInstance,err:=DeploymentClient.Get(context.TODO(),newDeploymentInstance.Name,metaV1.GetOptions{})
	if err !=nil{
		return err
	}
	
	deploymentInstance.Spec.Replicas= &newDeploymentInstance.Replicas
	deploymentInstance.Spec.RevisionHistoryLimit =&newDeploymentInstance.RevisionHistoryLimit
	deploymentInstance.Spec.Template.Spec.Containers[0].Image= newDeploymentInstance.ImgUrl
	deploymentInstance.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = newDeploymentInstance.ContainerPort
	deploymentInstance.Spec.Template.Spec.Containers[0].Ports[0].HostPort = newDeploymentInstance.HostPort


	_,err=DeploymentClient.Update(context.TODO(),deploymentInstance,metaV1.UpdateOptions{})
	if err !=nil{
		return err
	}

	return nil

}
