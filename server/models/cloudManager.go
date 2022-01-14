package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/CosmicBDry/gocmdb/server/cloud"
	"github.com/astaxie/beego/orm"
)

type CloudPlatform struct {
	Id              int64             `orm:"column(id);" json:"id"`
	Name            string            `orm:"column(name);size(64)"  json:"name"`
	VirtualMachines []*VirtualMachine `orm:"reverse(many);on_delete(cascade)" json:"virtual_machines"`
	Types           string            `orm:"column(types);size(64)"  json:"types"`
	Addr            string            `orm:"column(addr);size(128)"  json:"addr"`
	AccessKey       string            `orm:"column(access_key);size(1024)" json:"-"`
	SecretKey       string            `orm:"column(secret_key);size(1024)" json:"-"`
	Region          string            `orm:"column(region);size(64)" json:"region"`
	Remark          string            `orm:"column(remark);size(1024)" json:"remark"`
	CreatedTime     *time.Time        `orm:"column(created_time);auto_now_add" json:"created_time"`
	DeletedTime     *time.Time        `orm:"column(deleted_time);null;default(null)" json:"-"`
	SyncTime        *time.Time        `orm:"column(sync_time);null;default(null)" json:"sync_time"`
	User            *User             `orm:"column(user);rel(fk);null;" json:"-"`
	Status          int               `orm:"column(status);" json:"status"`
}

type CloudPlatformManager struct{}

//datatable页面通过数据库的limit、offset返回数据，或返回全部数据
func (c *CloudPlatformManager) QueryList(draw, start, length int64, colname, colsort, searchvalue, platformsearch string) (int64, int64, []CloudPlatform) {

	var result []CloudPlatform

	ormer := orm.NewOrm()

	queryset := ormer.QueryTable(&CloudPlatform{}).Filter("DeletedTime__isnull", true).Filter("Name__icontains", platformsearch)

	Total, _ := queryset.Count()
	TotalFilter := Total

	if searchvalue == "" {
		if colsort == "desc" { //判断为逆序
			colname = "-" + colname
			queryset.Limit(length).Offset(start).OrderBy(colname).All(&result)

		} else if colsort == "asc" { //判断为升序
			queryset.Filter("DeletedTime__isnull", true).Limit(length).Offset(start).OrderBy(colname).All(&result)
		} else {
			queryset.Filter("DeletedTime__isnull", true).All(&result)
		}

	} else {
		con := orm.NewCondition()
		con = con.Or("Name__icontains", searchvalue).Or("created_time__icontains", searchvalue).Or("SyncTime__icontains", searchvalue)
		con = con.Or("Region__icontains", searchvalue)
		queryset.SetCond(con).All(&result)
		TotalFilter, _ = queryset.SetCond(con).Count()

	}

	return Total, TotalFilter, result
}

func (c *CloudPlatformManager) GetbyId(id int64) *CloudPlatform {
	obj := &CloudPlatform{}
	ormer := orm.NewOrm()
	err := ormer.QueryTable(&CloudPlatform{}).Filter("DeletedTime__isnull", true).Filter("Id__exact", id).One(obj)
	if err != nil {

		return nil
	}
	return obj
}

func (c *CloudPlatformManager) GetbyName(name string) *CloudPlatform {
	obj := &CloudPlatform{}
	ormer := orm.NewOrm()
	err := ormer.QueryTable(&CloudPlatform{}).Filter("DeletedTime__isnull", true).Filter("Name__exact", name).One(obj)
	if err != nil {

		return nil
	}
	return obj
}

func (c *CloudPlatformManager) Create(name, types, addr, region, accesskey, secretkey, remark string) bool {
	ormer := orm.NewOrm()
	obj := &CloudPlatform{
		Name:      name,
		Types:     types,
		Addr:      addr,
		Region:    region,
		AccessKey: accesskey,
		SecretKey: secretkey,
		Remark:    remark,
	}
	if rows, err := ormer.Insert(obj); rows > 0 && err == nil {
		return true
	}
	return false
}

func (c *CloudPlatformManager) Modify(pk int64, name, types, addr, region, accesskey, secretkey, remark string) {

	ormer := orm.NewOrm()
	platform := &CloudPlatform{Id: pk}
	if err := ormer.Read(platform); err != nil {
		fmt.Println(err)
	}
	platform.Name = name
	platform.Types = types
	platform.Addr = addr
	platform.Region = region
	if accesskey != "" {
		platform.AccessKey = accesskey
	}
	if secretkey != "" {
		platform.SecretKey = secretkey
	}

	platform.Remark = remark

	if rows, err := ormer.Update(platform); err != nil {
		fmt.Println(rows, err)
	}
}

func (c *CloudPlatformManager) Disable(pk int64) {
	ormer := orm.NewOrm()

	cond := orm.NewCondition()
	cond = cond.And("DeletedTime__isnull", true)

	queryset := ormer.QueryTable(&CloudPlatform{}).SetCond(cond)

	queryset.Filter("Id__exact", pk).Update(orm.Params{"Status": StatusLock})

}

func (c *CloudPlatformManager) Enable(pk int64) {

	ormer := orm.NewOrm()

	cond := orm.NewCondition()
	cond = cond.And("DeletedTime__isnull", true)

	queryset := ormer.QueryTable(&CloudPlatform{}).SetCond(cond)

	queryset.Filter("Id__exact", pk).Update(orm.Params{"Status": StatusUnLock})

}

func (c *CloudPlatformManager) Delete(pk int64) bool {

	ormer := orm.NewOrm()

	cond := orm.NewCondition()
	cond = cond.And("DeletedTime__isnull", true)

	queryset := ormer.QueryTable(&CloudPlatform{}).SetCond(cond)

	now := time.Now()
	platform := &CloudPlatform{Id: pk}

	fmt.Println("platform one: ")
	if num, err := queryset.Filter("Id__exact", pk).Update(orm.Params{"DeletedTime": &now}); num > 0 && err == nil {
		ormer.Read(platform)
		DefaultVirtualMachineManager.DeleteByCloudPlatform(platform)
		return true
	} else {
		//fmt.Println(num, err)
		return false
	}

	return true
}

/*func (c *CloudPlatformManager) Delete(pk int64) bool {

	ormer := orm.NewOrm()

	result := ormer.Raw("delete from  cloud_platform  where id =?", pk)
	fmt.Println(result.Exec())

	return true
}*/

func (c *CloudPlatformManager) CloudSyncTime(platform *CloudPlatform, now time.Time) {
	platform.SyncTime = &now
	ormer := orm.NewOrm()
	ormer.Update(platform)
}

func NewCloudPlatformManager() *CloudPlatformManager {
	return &CloudPlatformManager{}
}

var DefautlCloudPlatformManager = NewCloudPlatformManager()

type VirtualMachine struct {
	Id            int64          `orm:"column(id)" json:"id" `
	InstanceId    string         `orm:"column(instance_id);size(1024)" json:"instance_id"`
	Name          string         `orm:"column(name);size(64)" json:"name"`
	CloudPlatform *CloudPlatform `orm:"column(cloud_platform);rel(fk)" json:"cloud_platform"`
	Status        string         `orm:"column(status);size(1024)" json:"status"`
	Mem           int64          `orm:"column(mem)" json:"mem"`
	OS            string         `orm:"column(os);size(128)" json:"os"`
	CPU           int64          `orm:"column(cpu)" json:"cpu"`
	PublicAddrs   string         `orm:"column(public_addrs);size(1024)" json:"public_addrs"`
	PrivateAddrs  string         `orm:"column(private_addrs);size(1024)" json:"private_addrs"`
	VmCreatedTime string         `orm:"column(vm_created_time);size(1024)" json:"vm_created_time"`
	VmExpiredTime string         `orm:"column(vm_expired_time);size(1024)" json:"-"`
	CreatedTime   *time.Time     `orm:"column(created_time);auto_now_add"`
	UpdatedTime   *time.Time     `orm:"column(updated_time);auto_now"`
	DeletedTime   *time.Time     `orm:"column(deleted_time);null;default(null)"`
}

type VirtualMachineManager struct {
}

func NewVirtualMachineManager() *VirtualMachineManager {
	return &VirtualMachineManager{}
}

func (c *VirtualMachineManager) QueryList(draw, start, length, platform int64, colname, colsort, searchvalue string) (int64, int64, []VirtualMachine) {

	var result []VirtualMachine

	ormer := orm.NewOrm()
	con := orm.NewCondition()
	con = con.And("DeletedTime__isnull", true)
	queryset := ormer.QueryTable(&VirtualMachine{}).RelatedSel("CloudPlatform").SetCond(con)
	//RelatedSel来启用表的关联关系查询
	if platform > 0 {
		con = con.And("CloudPlatform__exact", platform)
		queryset = queryset.SetCond(con)
		//fmt.Println("platform > 0")
	}

	Total, _ := queryset.Count()
	TotalFilter := Total

	if searchvalue == "" {
		if colsort == "desc" { //判断为逆序
			colname = "-" + colname
			queryset.Limit(length).Offset(start).OrderBy(colname).All(&result)

		} else if colsort == "asc" { //判断为升序
			queryset.Limit(length).Offset(start).OrderBy(colname).All(&result)
		} else {
			queryset.All(&result)
		}

	} else {
		con1 := orm.NewCondition()
		con1 = con1.Or("Name__icontains", searchvalue).Or("VmCreatedTime__icontains", searchvalue).Or("Status__icontains", searchvalue)
		con = con.AndCond(con1) //con最好不要直接加Or，如con.Or(),否则con最初的条件可能会被忽略掉，可通过AndCond(con1)来解决此问题
		queryset.SetCond(con).All(&result)

		TotalFilter, _ = queryset.SetCond(con).Count()

	}
	return Total, TotalFilter, result
}

func (c *VirtualMachineManager) Delete(id int64) error {

	ormer := orm.NewOrm()
	now := time.Now()
	rows, err := ormer.QueryTable(&VirtualMachine{}).Filter("DeletedTime__isnull", true).Filter("Id__exact", id).Update(orm.Params{"DeletedTime": &now})

	if rows <= 0 || err != nil {
		return errors.New("删除失败")
	}

	return nil
}

func (c *VirtualMachineManager) DeleteByCloudPlatform(platform *CloudPlatform) error {
	ormer := orm.NewOrm()

	now := time.Now()
	rows, err := ormer.QueryTable(&VirtualMachine{}).Filter("DeletedTime__isnull", true).Filter("CloudPlatform__exact", platform).Update(orm.Params{"DeletedTime": &now})

	if rows <= 0 || err != nil {

		return errors.New("删除失败")

	}
	return nil
}

func (v *VirtualMachineManager) SyncVmByCloudPlatform(C *CloudPlatform, instance *cloud.Instance) error {
	ormer := orm.NewOrm()
	//不能查询到就insert创建，然后再update更新
	//能查到则直接update更新

	vm := &VirtualMachine{CloudPlatform: C, InstanceId: instance.InstanceId}
	if err := ormer.Read(vm, "CloudPlatform", "InstanceId"); err == orm.ErrNoRows {
		if rows, err := ormer.Insert(vm); rows <= 0 || err != nil {
			//fmt.Println(rows, err)
			return errors.New("虚拟机同步创建失败")
		}
	}
	//vm.InstanceId = instance.InstanceId
	vm.Name = instance.Name
	vm.Status = instance.Status
	vm.Mem = instance.Mem
	vm.OS = instance.OS
	vm.CPU = instance.CPU
	vm.PublicAddrs = strings.Join(instance.PublicAddrs, ",")
	vm.PrivateAddrs = strings.Join(instance.PrivateAddrs, ",")
	vm.VmCreatedTime = instance.CreatedTime
	if rows, err := ormer.Update(vm); rows <= 0 || err != nil {
		return errors.New("虚拟机同步更新失败")
	}

	return nil
}

func (v *VirtualMachineManager) SyncVmStatus(C *CloudPlatform, now time.Time) {
	ormer := orm.NewOrm()
	queryset := ormer.QueryTable(&VirtualMachine{})
	queryset.Filter("CloudPlatform__exact", C).Filter("UpdatedTime__lt", &now).Update(orm.Params{"DeletedTime": &now})
	queryset.Filter("CloudPlatform__exact", C).Filter("UpdatedTime__gte", &now).Update(orm.Params{"DeletedTime": nil})
}

func (v *VirtualMachineManager) GetById(pk int64) VirtualMachine {
	machine := VirtualMachine{}
	ormer := orm.NewOrm()
	queryset := ormer.QueryTable(&VirtualMachine{}).RelatedSel("CloudPlatform").Filter("DeletedTime__isnull", true)
	err := queryset.Filter("Id__exact", pk).One(&machine)
	fmt.Printf("Virtualmachine one %#v\n", err)
	return machine
}

var DefaultVirtualMachineManager = NewVirtualMachineManager()

func init() {

	orm.RegisterModel(new(CloudPlatform), new(VirtualMachine))

}
