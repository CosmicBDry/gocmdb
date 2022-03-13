package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

type GitProject struct {
	Id             int64      `orm:"column(id)" json:"id"`
	ProjectId      int64      `orm:"column(project_id)" json:"project_id"`
	Name           string     `orm:"column(name);size(128)" json:"name"`
	NameSpace      string     `orm:"column(namespace);size(128)" json:"namespace"`
	Branch         string     `orm:"column(branch);size(128)" json:"branch"`
	Description    string     `orm:"column(description);size(512)" json:"-"`
	Author         string     `orm:"column(author);size(64)"  json:"-"`
	AuthorEmail    string     `orm:"column(author_email);size(64)"json:"-"`
	CurrentVersion string     `orm:"column(current_version);size(128)"json:"-"`
	LastVersion    string     `orm:"column(last_version);size(128)"json:"-"`
	EventName      string     `orm:"column(event_name);size(64)"json:"-"`
	Message        string     `orm:"column(message);size(128)"json:"-"`
	CommitTime     string     `orm:"column(commit_time);size(32)"json:"-"`
	HttpUrl        string     `orm:"column(http_url);size(128)"json:"-"`
	SShUrl         string     `orm:"column(ssh_url);size(128)"json:"-"`
	PackageFile    string     `orm:"column(package_file);size(1024)"json:"package_file"`
	ReleaseMachine string     `orm:"column(release_machine);size(64)"json:"-"`
	BackendHost    string     `orm:"column(backend_host);size(1024)"json:"-"`
	AutoDeploy     bool       `orm:"column(auto_deploy);default(false)"json:"auto_deploy"`
	CreatedTime    *time.Time `orm:"column(created_time);auto_now_add"json:"-"`
	UpdatedTime    *time.Time `orm:"column(updated_time);auto_now"json:"-"`
	DeletedTime    *time.Time `orm:"column(deleted_time);null;defualt(null)"json:"-"`
}

func (g *GitProject) GitCreate() error {
	ormer := orm.NewOrm()
	if g.ProjectId == 0 {
		return errors.New("ProjectId is 0,Not effective id")
	}
	gitproject := &GitProject{ProjectId: g.ProjectId}
	if err := ormer.Read(gitproject, "ProjectId"); err != nil {
		if err == orm.ErrNoRows {
			ormer.Insert(g)
			return nil
		}
		return err
	}

	if g.CurrentVersion == gitproject.CurrentVersion {
		return nil
	}
	//fmt.Printf("------------------------------------------->%#v\n", g)
	g.BackendHost = gitproject.BackendHost
	g.ReleaseMachine = gitproject.ReleaseMachine
	g.Id = gitproject.Id
	ormer.Update(g)
	//fmt.Println(num, err)
	return nil
}

func GitLabGetList(start, length int64, colname, colsort, searchvalue string) (int64, int64, []GitProject) {
	var results []GitProject
	ormer := orm.NewOrm()
	cond := orm.NewCondition()
	cond = cond.And("DeletedTime__isnull", true)
	cond1 := orm.NewCondition()
	cond1 = cond1.Or("ProjectId__exact", searchvalue).Or("Name__icontains", searchvalue).Or("NameSpace__icontains", searchvalue)
	queryset := ormer.QueryTable(&GitProject{}).SetCond(cond)
	total, _ := queryset.Count()
	totalFilter := total

	if searchvalue == "" {

		if colsort == "asc" {
			QueryResult := queryset.Limit(int(length)).Offset(int(start))
			QueryResult.OrderBy(colname).All(&results)
			totalFilter, _ = QueryResult.Count()

		} else if colsort == "desc" {
			colname = "-" + colname
			QueryResult := queryset.Limit(int(length)).Offset(int(start))
			QueryResult.OrderBy(colname).All(&results)
			totalFilter, _ = QueryResult.Count()
		}

	} else {
		queryset.SetCond(cond1).All(&results)
		cond = cond.AndCond(cond1)
		totalFilter, _ = queryset.SetCond(cond).Count()
	}

	return total, totalFilter, results

}

func GitGetById(pk int64) (*GitProject, error) {
	ormer := orm.NewOrm()
	gitproject := &GitProject{Id: pk}

	if err := ormer.Read(gitproject); err != nil {
		return nil, err
	}

	return gitproject, nil
}

func ReleaseConfigModify(id int64, releaseHost, backendHost, packagefile string, autodeploy bool) error {
	ormer := orm.NewOrm()
	gitproject := &GitProject{Id: id}
	if err := ormer.Read(gitproject); err == nil {
		if gitproject.ReleaseMachine != releaseHost || gitproject.BackendHost != backendHost || gitproject.PackageFile != packagefile || gitproject.AutoDeploy != autodeploy {
			gitproject.ReleaseMachine = releaseHost
			gitproject.BackendHost = backendHost
			gitproject.PackageFile = packagefile
			gitproject.AutoDeploy = autodeploy
			_, err := ormer.Update(gitproject, "ReleaseMachine", "BackendHost", "PackageFile", "AutoDeploy")
			return err
		}
		return errors.New("未发生更改,无需更新！")
	} else {
		return err
	}

}

func init() {
	orm.RegisterModel(&GitProject{})
}
