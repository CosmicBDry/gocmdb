package auth

import (
	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/astaxie/beego/context"
)

type AuthPlugin interface {
	Name() string
	Is(c *context.Context) bool
	IsLogin(c *LoginRequireController) *models.User
	GoToLogin(c *LoginRequireController)
	Login(c *AuthController) bool
	Logout(c *AuthController)
}

type Manager struct {
	plugin map[string]AuthPlugin
}

func NewManager() *Manager {
	return &Manager{
		plugin: map[string]AuthPlugin{},
	}
}

//注册插件，key为插件Name()方法返回的名称，值为插件（有session、token等插件）------------------------->
func (m *Manager) Register(p AuthPlugin) {

	m.plugin[p.Name()] = p

}

//通过用户提交的数据来判断插件类型，是session插件还是token插件------------------------->
func (m *Manager) GetPlugin(c *context.Context) AuthPlugin {
	for _, plugin := range m.plugin { //遍历m.plugin中注册的多个插件
		if plugin.Is(c) { //如其中一个插件判断为true，则证明为此插件
			return plugin //返回插件
		}
	}

	return nil
}

func (m *Manager) IsLogin(c *LoginRequireController) *models.User {
	if plugin := m.GetPlugin(c.Ctx); plugin != nil {
		return plugin.IsLogin(c)
	}
	return nil

}

func (m *Manager) GoToLogin(c *LoginRequireController) {
	if plugin := m.GetPlugin(c.Ctx); plugin != nil {
		plugin.GoToLogin(c)
	}

}

func (m *Manager) Login(c *AuthController) bool {
	if plugin := m.GetPlugin(c.Ctx); plugin != nil {
		return plugin.Login(c)
	}

	return false
}

func (m *Manager) Logout(c *AuthController) {
	if plugin := m.GetPlugin(c.Ctx); plugin != nil {
		plugin.Logout(c)
	}

}

var defaultManager = NewManager() //实现只定义一个manager，实现所有controller共用一个manager
