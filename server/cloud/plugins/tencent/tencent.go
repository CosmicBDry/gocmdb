package tencent

import (
	"fmt"

	"github.com/CosmicBDry/gocmdb/server/cloud"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type TencentCloud struct {
	Addr       string
	Region     string
	AccessKey  string
	SecretKey  string
	Credential *common.Credential
	Profile    *profile.ClientProfile
}

func (t *TencentCloud) InstanceStatusTransfer(status string) string {
	Maps := map[string]string{
		"PENDING":       cloud.StatusPending,
		"LAUNCH_FAILED": cloud.StatusLauchFailed,
		"RUNNING":       cloud.StatusRunning,
		"STOPPED":       cloud.StatusStopped,
		"STARTING":      cloud.StatusStarting,
		"STOPPING":      cloud.StatusStoping,
		"REBOOTING":     cloud.StatusRestarting,
		"SHUTDOWN":      cloud.StatusShutdown,
		"TERMINATING":   cloud.StatusTerminating,
	}
	if result, ok := Maps[status]; ok {
		return result
	}

	return cloud.StatusUnknow
}

func (t *TencentCloud) Type() string {

	return "tencentCloud"

}

func (t *TencentCloud) Name() string {
	return "腾讯云"
}

func (t *TencentCloud) Init(addr, region, accesskey, secretkey string) {
	t.Addr = addr
	t.Region = region
	t.AccessKey = accesskey
	t.SecretKey = secretkey
	t.Credential = common.NewCredential(accesskey, secretkey)
	t.Profile = profile.NewClientProfile()
	t.Profile.HttpProfile.Endpoint = t.Addr
}
func (t *TencentCloud) TestConnect() error {
	client, err := cvm.NewClient(t.Credential, t.Region, t.Profile)
	if err != nil {
		fmt.Println("client error", err)
		return err
	}
	request := cvm.NewDescribeRegionsRequest()

	_, err = client.DescribeRegions(request)

	if err != nil {
		fmt.Println("response error", err)
		return err
	}

	return nil
}

func (t *TencentCloud) GetInstanceByOffsetLimit(offset, limit int64) (int64, []*cloud.Instance) {

	client, err := cvm.NewClient(t.Credential, t.Region, t.Profile)
	if err != nil {
		fmt.Println("client error:", err)
		return 0, nil
	}
	request := cvm.NewDescribeInstancesRequest()

	request.Offset = &offset
	request.Limit = &limit

	response, err := client.DescribeInstances(request)
	if err != nil {
		fmt.Println("reposne error:", err)
		return 0, nil
	}

	Total := response.Response.TotalCount
	Instances := response.Response.InstanceSet
	rt := make([]*cloud.Instance, len(Instances))

	for index, instance := range Instances {

		pubilicaddrs := make([]string, len(instance.PublicIpAddresses))
		for index, value := range instance.PublicIpAddresses {
			pubilicaddrs[index] = *value
		}
		privateaddrs := make([]string, len(instance.PrivateIpAddresses))
		for index, value := range instance.PrivateIpAddresses {
			privateaddrs[index] = *value
		}

		rt[index] = &cloud.Instance{
			InstanceId:   *instance.InstanceId,
			Name:         *instance.InstanceName,
			Mem:          *instance.Memory * 1024,
			CPU:          *instance.CPU,
			OS:           *instance.OsName,
			PublicAddrs:  pubilicaddrs,
			PrivateAddrs: privateaddrs,
			Status:       t.InstanceStatusTransfer(*instance.InstanceState),
			CreatedTime:  *instance.CreatedTime,
			//ExpiredTime:  *instance.ExpiredTime,注意腾讯云后付费模式没有过期时间，为null
		}

	}

	//fmt.Printf("%v    %#v\n", *Total, rt[0])
	return *Total, rt
}

func (t *TencentCloud) GetInstances() []*cloud.Instance {
	var (
		offset int64 = 0
		limit  int64 = 100
		total  int64 = 1
		rt     []*cloud.Instance
	)

	for offset < total {
		var instances []*cloud.Instance
		total, instances = t.GetInstanceByOffsetLimit(offset, limit)
		if offset == 0 {
			rt = make([]*cloud.Instance, 0, total)
		}

		rt = append(rt, instances...)
		offset += limit
		//fmt.Printf("%#v\n", instances[0])
	}

	return rt
}

func (t *TencentCloud) StartInstance(instanceid string) error {
	client, err := cvm.NewClient(t.Credential, t.Region, t.Profile)
	if err != nil {
		fmt.Println("client", err)
		return err
	}
	request := cvm.NewStartInstancesRequest()
	//request.InstanceIds = common.StringPtrs([]string{instanceid})
	request.InstanceIds = []*string{&instanceid}
	_, err = client.StartInstances(request)
	//reqId := response.Response.RequestId
	if err != nil {
		fmt.Println("response", err)
		return err
	}
	return nil
}

func (t *TencentCloud) StopInstance(instanceid string) error {

	client, err := cvm.NewClient(t.Credential, t.Region, t.Profile)
	if err != nil {
		fmt.Println("client", err)
		return err
	}
	request := cvm.NewStopInstancesRequest()
	//request.InstanceIds = common.StringPtrs([]string{instanceid})
	request.InstanceIds = []*string{&instanceid}
	_, err = client.StopInstances(request)
	//reqId := response.Response.RequestId
	if err != nil {
		fmt.Println("response", err)
		return err
	}
	return nil
}

func (t *TencentCloud) RestartInstance(instanceid string) error {
	client, err := cvm.NewClient(t.Credential, t.Region, t.Profile)
	if err != nil {
		fmt.Println("client", err)
		return err
	}
	request := cvm.NewRebootInstancesRequest()
	request.InstanceIds = []*string{&instanceid}
	_, err = client.RebootInstances(request)
	if err != nil {
		fmt.Println("response", err)
		return err
	}
	return nil
}

func (t *TencentCloud) TerminateInstance(instanceid string) error {
	client, err := cvm.NewClient(t.Credential, t.Region, t.Profile)
	if err != nil {
		fmt.Println("client", err)
		return err
	}
	request := cvm.NewTerminateInstancesRequest()
	request.InstanceIds = []*string{&instanceid}
	_, err = client.TerminateInstances(request)
	if err != nil {
		fmt.Println("response", err)
		return err
	}
	return nil
}

func init() {

	cloud.DefaultCloudManager.Register(new(TencentCloud))
}
