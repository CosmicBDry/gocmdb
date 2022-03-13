package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Alert struct {
	Id          int64      `orm:"column(id)"json:"id"`
	AlertName   string     `orm:"column(alertname);size(64)" json:"alertname"`
	Instance    string     `orm:"column(instance);size(64)" json:"instance"`
	Serverity   string     `orm:"column(serverity);size(16)" json:"serverity"`
	Status      string     `orm:"column(status);size(16)" json:"status"`
	FingerPrint string     `orm:"column(finger_print);size(64)" json:"fingerprint"`
	Description string     `orm:"column(description);size(128)" json:"description"`
	Summary     string     `orm:"column(summary);size(128)"json:"summary"`
	StartsAt    *time.Time `orm:"column(starts_at);null"json:"startsat"`
	EndsAt      *time.Time `orm:"column(ends_at);null" json:"endsat"`
	Labels      string     `orm:"column(labels);type(longtext);" json:"-"`
	CreatedTime *time.Time `orm:"column(created_time);auto_now_add;"json:"-"`
	UpdatedTime *time.Time `orm:"column(updated_time);auto_now;"json:"-"`
	DeletedTime *time.Time `orm:"column(deleted_time);null"json:"-"`
}

func (a *Alert) CreateOrUpdate() error {
	alert := &Alert{FingerPrint: a.FingerPrint}
	ormer := orm.NewOrm()
	if err := ormer.Read(alert, "FingerPrint"); err == orm.ErrNoRows {

		if _, err := ormer.Insert(a); err != nil {
			return err
		}
	} else {
		a.Id = alert.Id
		if alert.Status != a.Status {
			if _, err := ormer.Update(a); err != nil {
				return err
			}
		}

	}
	return nil
}

func GetList(start, length int64, colname, colsort, searchvalue string) (int64, int64, []Alert) {
	var results []Alert
	ormer := orm.NewOrm()
	cond := orm.NewCondition()
	cond = cond.And("DeletedTime__isnull", true)
	cond1 := orm.NewCondition()
	cond1 = cond1.Or("AlertName__icontains", searchvalue).Or("Instance__icontains", searchvalue).Or("Serverity__icontains", searchvalue).Or("Status__icontains", searchvalue)
	cond1 = cond1.Or("StartsAt__icontains", searchvalue).Or("EndsAt__icontains", searchvalue).Or("FingerPrint", searchvalue)
	queryset := ormer.QueryTable(&Alert{}).SetCond(cond)
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

func Delete(pk int64) error {

	ormer := orm.NewOrm()

	queryset := ormer.QueryTable(&Alert{})
	times := time.Now()
	_, err := queryset.Filter("DeletedTime__isnull", true).Filter("Id__exact", pk).Update(orm.Params{"DeletedTime": &times})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	orm.RegisterModel(new(Alert))

}
