package cloud

type CloudManager struct {
	Plugins map[string]CloudI
}

func NewCloudManager() *CloudManager {
	return &CloudManager{
		Plugins: make(map[string]CloudI),
	}

}

func (c *CloudManager) Register(i CloudI) {
	c.Plugins[i.Type()] = i
}

var DefaultCloudManager = NewCloudManager()
