package aliyun

import (
	"github.com/CosmicBDry/gocmdb/server/cloud"
)

type Aliyun struct{}

func (a *Aliyun) Type() string {
	return "aliyun"
}

func (a *Aliyun) Name() string {
	return "阿里云"
}

func (a *Aliyun) Init(string, string, string, string) {

}

func (a *Aliyun) TestConnect() error {
	return nil
}

func (a *Aliyun) GetInstances() []*cloud.Instance {
	return nil
}

func (a *Aliyun) StartInstance(string) error {
	return nil
}

func (a *Aliyun) StopInstance(string) error {
	return nil
}

func (a *Aliyun) RestartInstance(string) error {
	return nil
}
func (a *Aliyun) TerminateInstance(string) error {
	return nil
}
func init() {
	cloud.DefaultCloudManager.Register(new(Aliyun))
}
