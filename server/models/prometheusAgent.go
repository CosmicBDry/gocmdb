package models

import (
	"time"

	//"encoding/json"

	"github.com/astaxie/beego/orm"
	//"gopkg.in/yaml.v2"
)

type Agent struct {
	Id            int64      `orm:"column(id)"json:"id"`
	Hostname      string     `orm:"column(hostname);size(32)"json:"hostname"`
	UUID          string     `orm:"column(uuid);size(32)"json:"uuid"`
	Addr          string     `orm:"column(addr);size(64)"json:"addr"`
	Description   string     `orm:"column(description);size(1024)"json:"-"`
	Config        string     `orm:"column(config);type(text)"json:"-"`
	ConfigVersion string     `orm:"column(config_version);size(64)"json:"-"`
	HeartBeat     *time.Time `orm:"column(heart_beat);null"json:"heartbeat"`
	CreatedTime   *time.Time `orm:"column(created_time);auto_now_add"json:"-"`
	UpdatedTime   *time.Time `orm:"column(updated_time);auto_now"json:"-"`
	DeletedTime   *time.Time `orm:"column(deleted_time);null;default(null)"json:"-"`
}

type AgentManager struct {
}

func NewAgentManager() *AgentManager {

	return &AgentManager{}
}

func (agentmanager *AgentManager) Register(agent *Agent) (error, string) {
	ormer := orm.NewOrm()
	tmpAgent := &Agent{UUID: agent.UUID}
	if err := ormer.Read(tmpAgent, "UUID"); err == orm.ErrNoRows {
		if _, err := ormer.Insert(agent); err != nil {
			return err, ""
		}
		return nil, "inserted"
	} else if err == nil {
		if tmpAgent.Hostname != agent.Hostname || tmpAgent.Addr != agent.Addr {
			tmpAgent.Hostname = agent.Hostname
			tmpAgent.Addr = agent.Addr
			if _, err := ormer.Update(tmpAgent); err != nil {
				return err, ""
			}
			return nil, "updated"
		}

	}
	return nil, ""
}

func (agentmanager *AgentManager) HeartBeat(uuid string) error {

	ormer := orm.NewOrm()
	times := time.Now()
	tmpagent := &Agent{UUID: uuid}

	if err := ormer.Read(tmpagent, "UUID"); err != nil {
		//fmt.Println("read error", err)
		return err

	} else {
		tmpagent.HeartBeat = &times
		tmpagent.DeletedTime = nil
		if _, err := ormer.Update(tmpagent, "HeartBeat", "DeletedTime"); err != nil {
			//fmt.Println("update error", err)
			return err
		}

	}

	return nil

}

func (agentmanager *AgentManager) GetConfig(uuid, configversion string) (string, string, error) {
	//config := Jobs{}
	ormer := orm.NewOrm()
	agent := Agent{}
	queryset := ormer.QueryTable(&Agent{})

	err := queryset.Filter("DeletedTime__isnull", true).Filter("UUID__exact", uuid).One(&agent)
	if err != nil {
		return "", "", err
	}
	if configversion < agent.ConfigVersion {
		//yaml.Unmarshal([]byte(agent.Config), &config)
		//err := json.Unmarshal([]byte(agent.Config), &config)
		//fmt.Println(err)
		return agent.Config, agent.ConfigVersion, nil
	}

	return "", "", nil
}

func AgentList(start, length int64, colname, colsort, searchvalue string) (int64, int64, []Agent) {
	var results []Agent
	ormer := orm.NewOrm()
	cond := orm.NewCondition()
	cond = cond.And("DeletedTime__isnull", true)
	cond1 := orm.NewCondition()
	cond1 = cond1.Or("Hostname__icontains", searchvalue).Or("Addr__icontains", searchvalue).Or("HeartBeat__icontains", searchvalue)

	queryset := ormer.QueryTable(&Agent{}).SetCond(cond)
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

func (agentmanager *AgentManager) SetConfig(id int64, config, configversion string) error {

	agent := &Agent{Id: id}
	ormer := orm.NewOrm()
	ormer.Read(agent)
	agent.Config = config
	agent.ConfigVersion = configversion
	_, err := ormer.Update(agent, "Config", "ConfigVersion")
	if err != nil {
		return err
	}
	return nil

}

func (agentmanager *AgentManager) GetAgentById(id int64) *Agent {
	agent := &Agent{Id: id}
	ormer := orm.NewOrm()
	ormer.Read(agent)
	return agent
}

var DefaultAgentManager = NewAgentManager()

type Jobs struct {
	Global        interface{} `yaml:"global" json:"global"`
	Alerting      interface{} `yaml:"alerting" json:"alerting"`
	RuleFiles     interface{} `yaml:"rule_files" json:"rule_files"`
	ScrapeConfigs []*struct {
		JobName     string `yaml:"job_name" json:"job_name"`
		MetricsPath string `yaml:"metrics_path,omitempty" json:"metrics_path"`
		Scheme      string `yaml:"scheme,omitempty" json:"scheme"`
		BasicAuth   struct {
			Username string `yaml:"username" json:"username"`
			Password string `yaml:"password" json:"password"`
		} `yaml:"basic_auth,omitempty" json:"basic_auth"`
		FileSdConfigs []*struct {
			Files           []string `yaml:"files" json:"files"`
			RefreshInterval string   `yaml:"refresh_interval"json:"refresh_interval"`
		} `yaml:"file_sd_configs"json:"file_sd_configs"`
	} `yaml:"scrape_configs"json:"scrape_configs"`
}

func init() {

	orm.RegisterModel(&Agent{})

}
