package forms

import (
	"bytes"
	"strings"

	"github.com/CosmicBDry/gocmdb/server/models"
	"github.com/astaxie/beego/validation"
)

type ProjectInfo struct {
	Project_id  int64  `json:"id"`
	ProjectName string `json:"name"`
	Description string `json:"description"`
	NameSpace   string `json:"namespace"`
	GitHttpUrl  string `json:"git_http_url"`
	GitSShUrl   string `json:"git_ssh_url"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type Commmits struct {
	Title        string `json:"title"`
	Message      string `json:"message"`
	CommitTime   string `json:"timestamp"`
	CommitAuthor Author `json:"author"`
}

type ProjectHook struct {
	Version       string      `json:"after"`
	LastVersion   string      `json:"before"`
	Branch        string      `json:"ref"`
	EventName     string      `json:"event_name"`
	Project       ProjectInfo `json:"project"`
	ProjectCommit []Commmits  `json:"commits"`
}

func (p *ProjectHook) GitFormToModel() *models.GitProject {
	var author, email, message, committime string
	for _, commits := range p.ProjectCommit {
		author = commits.CommitAuthor.Name
		email = commits.CommitAuthor.Email
		message = commits.Message
		committime = commits.CommitTime
	}

	return &models.GitProject{
		ProjectId:      p.Project.Project_id,
		Name:           p.Project.ProjectName,
		NameSpace:      p.Project.NameSpace,
		Branch:         strings.TrimPrefix(p.Branch, "refs/heads/"),
		Description:    p.Project.Description,
		Author:         author,
		AuthorEmail:    email,
		CurrentVersion: p.Version,
		LastVersion:    p.LastVersion,
		EventName:      p.EventName,
		Message:        message,
		CommitTime:     committime,
		HttpUrl:        p.Project.GitHttpUrl,
		SShUrl:         p.Project.GitSShUrl,
	}

}

type ReleaseConfigForm struct {
	Id             int64  `form:"id"`
	ReleaseMachine string `form:"releasemachine"`
	BackendHost    string `form:"backendhost"`
	AutoDeploy     bool   `form:"auto_deploy"`
	PackageFile    string `form:"package_file"`
}

func (f *ReleaseConfigForm) Valid(v *validation.Validation) {

	f.ReleaseMachine = strings.TrimSpace(f.ReleaseMachine)
	f.BackendHost = strings.TrimSpace(f.BackendHost)

	var buffers bytes.Buffer
	var result string
	var err error
	buffers.WriteString(f.BackendHost + "," + f.ReleaseMachine)
	for {
		if result, err = buffers.ReadString(','); err != nil {

			v.IP(strings.TrimSpace(strings.Trim(result, ",")), "error.error.error").Message("多个ip间只能是逗号分隔或ip地址格式无效!")
			return
		}
		v.IP(strings.TrimSpace(strings.Trim(result, ",")), "error.error.error").Message("多个ip间只能是逗号分隔或ip地址格式无效!")
		if v.HasErrors() {
			return
		}

	}
}
