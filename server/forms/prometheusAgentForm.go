package forms

import (
	"github.com/CosmicBDry/gocmdb/server/models"
)

type AgentForm struct {
	UUID     string `json:"uuid"`
	Addr     string `json:"addr"`
	Hostname string `json:"hostname"`
}

func NewAgentForm() *AgentForm {
	return &AgentForm{}
}
func (a *AgentForm) ToModel() *models.Agent {

	return &models.Agent{
		UUID:     a.UUID,
		Addr:     a.Addr,
		Hostname: a.Hostname,
	}

}

type AgentConfigForm struct {
	Id            int64  `form:"id"`
	Config        string `form:"config"`
	ConfigVersion string `form:"configversion"`
}

var DefaultAgentForm = NewAgentForm()
