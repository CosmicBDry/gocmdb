package cloud

const (
	StatusPending     = "创建中"
	StatusLauchFailed = "创建失败"
	StatusRunning     = "运行中"
	StatusStarting    = "开机中"
	StatusStoping     = "关机中"
	StatusStopped     = "已关机"
	StatusRestarting  = "重启中"
	StatusTerminating = "销毁中"
	StatusShutdown    = "停止待销毁"
	StatusUnknow      = "未知状态"
)

type Instance struct {
	InstanceId   string
	Name         string
	CPU          int64
	OS           string
	Mem          int64
	PublicAddrs  []string
	PrivateAddrs []string
	Status       string
	CreatedTime  string
	ExpiredTime  string
}

type CloudI interface {
	Type() string
	Name() string
	Init(string, string, string, string)
	TestConnect() error
	GetInstances() []*Instance
	StartInstance(string) error
	StopInstance(string) error
	RestartInstance(string) error
	TerminateInstance(string) error
}
