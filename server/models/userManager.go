package models

import (
	"strings"
	"time"

	"github.com/CosmicBDry/gocmdb/server/utils"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type User struct {
	Id              int              `orm:"column(id);" form:"id"  json:"id"`
	Name            string           `orm:"column(name);size(32);" form:"name" json:"name"`
	Password        string           `orm:"column(password);size(1024);" form:"password"  json:"_"`
	ConfirmPassword string           `orm:"-"  form:"confirmpassword" json:"-"`
	Gender          int              `orm:"column(gender);default(0);" form:"gender" json:"gender"`
	Birthday        *time.Time       `orm:"column(birthday);type(date);null" form:"birthday" json:"birthday"`
	Tel             string           `orm:"column(tel);size(11);" form:"tel" json:"tel"`
	Email           string           `orm:"column(email);size(64);" form:"email" json:"email"`
	Addr            string           `orm:"column(addr);size(128);" form:"addr" json:"addr"`
	Remark          string           `orm:"column(remark);size(1024);" form:"remark" json:"remark"`
	IsSuperman      bool             `orm:"column(issuperman)";default(false);  json:"_"`
	Status          int              `orm:"column(status);" form:"status" json:"status"`
	CreatedTime     *time.Time       `orm:"column(created_time);auto_now_add;" form:"created_time" json:"created_time"`
	UpdatedTime     *time.Time       `orm:"column(updated_time);auto_now;" form:"updated_time" json:"updated_time"`
	DeletedTime     *time.Time       `orm:"column(deleted_time);null;default(null);"  json:"-"`
	Token           *Token           `orm:"reverse(one);"` //user表和token表的反引用与引用的关系，建立一对一的关系
	CloudPlatforms  []*CloudPlatform `orm:"reverse(many);" json:"cloud_platform" `
}

func (u *User) SetPassword(Passwd string) {

	u.Password = utils.Md5Salt(Passwd, "")

}

func (u *User) ValidPassword(Passwd string) bool {

	salt, _ := utils.SplitMd5Salt(u.Password)

	return utils.Md5Salt(Passwd, salt) == u.Password
}

func (u *User) IsLock() bool {

	return u.Status == StatusLock
}

//user表中的属性的值有效性判断
func (u *User) Valid(v *validation.Validation) {

	u.Name = strings.TrimSpace(u.Name)
	v.MinSize(u.Name, 2, "name.name.name").Message("用户名至少2个字符")
	v.MaxSize(u.Name, 32, "name.name.name").Message("用户名最多32个字符")

	v.Required(u.Birthday, "birth.birth.birth").Message("日期不可为空")

	u.Tel = strings.TrimSpace(u.Tel)

	if u.Tel != "" {

		if v.Numeric(u.Tel, "tel.tel.tel").Message("手机号码必须全为数字"); v.ErrorMap()["tel"] == nil {
			v.Length(u.Tel, 11, "tel.tel.tel").Message("手机号码非11位有效数字")
		}

	}

	//fmt.Println(v.ErrorMap())
	//fmt.Println(v.ErrorMap()["tel"])

	//通过u.Id == 0判断是否为新建用户
	if u.Id == 0 {

		u.Password = strings.TrimSpace(u.Password)
		v.MinSize(u.Password, 6, "password.password.password").Message("密码至少6个字符")
		v.MaxSize(u.Password, 32, "password.password.password").Message("密码最多32个字符")

		if DefautlUserManager.GetUserByName(u.Name) != nil {
			v.SetError("name", "用户名已存在")
		}
	} else {
		ormer := orm.NewOrm()
		cond := orm.NewCondition()
		cond = cond.And("Name__exact", u.Name).And("Gender__exact", u.Gender).And("Birthday__exact", u.Birthday)
		cond = cond.And("Tel__exact", u.Tel).And("Email__exact", u.Email).And("Addr__exact", u.Addr).And("Remark__exact", u.Remark)
		queryset := ormer.QueryTable(&User{})

		if count, _ := queryset.Exclude("Id__exact", u.Id).Filter("Name__exact", u.Name).Count(); count != 0 {
			v.SetError("name", "此用户已存在，无法修改")
		} else if count, _ := queryset.SetCond(cond).Count(); count != 0 {

			v.SetError("NotChanged", "当前未发生任何修改，若无需修改请关闭或重新修改")
		}

	}

	if u.ConfirmPassword != u.Password {
		v.SetError("password", "确认密码与输入密码不一致")
	}

}

type UserManager struct{}

func NewUserManager() *UserManager {
	return &UserManager{}
}

func (u *UserManager) GetUserByName(name string) *User {

	user := &User{Name: name}
	ormer := orm.NewOrm()
	if err := ormer.Read(user, "Name"); err == nil {
		return user
	}
	return nil
}
func (u *UserManager) GetUserById(id int) *User {
	user := &User{Id: id}
	ormer := orm.NewOrm()

	if err := ormer.Read(user); err == nil {
		ormer.LoadRelated(user, "Token") //开启表关联
		return user
	}
	return nil

}

//GetUsers获取所有用户的用户列表
func (u *UserManager) GetUsers(Draw, start, length int, colname, colsort, searchvalue string) (int64, int64, []User) {
	var userlist []User
	var Total, TotalFilter int64
	ormer := orm.NewOrm()
	queryset := ormer.QueryTable(&User{})
	Total, _ = queryset.Count()
	TotalFilter = Total
	if searchvalue == "" {
		if colsort == "desc" {
			colname = "-" + colname
			queryset.Limit(length).Offset(start).OrderBy(colname).All(&userlist)

		} else {
			queryset.Limit(length).Offset(start).OrderBy(colname).All(&userlist)
		}

	} else {
		con := orm.NewCondition()
		con = con.Or("Name__icontains", searchvalue).Or("addr__icontains", searchvalue).Or("created_time__icontains", searchvalue)
		con = con.Or("updated_time__icontains", searchvalue).Or("email__icontains", searchvalue)
		queryset.SetCond(con).All(&userlist)
		TotalFilter, _ = queryset.SetCond(con).Count()

	}
	//fmt.Println("database:", Total, TotalFilter)
	return Total, TotalFilter, userlist
}

func (u *UserManager) Create(user *User) error {

	ormer := orm.NewOrm()
	if _, err := ormer.Insert(user); err != nil {
		return err
	}
	return nil
}

func (u *UserManager) Modify(users *User) {
	ormer := orm.NewOrm()
	ormer.Update(users, "Name", "Gender", "Birthday", "Tel", "Email", "Addr", "Remark")

}

func (u *UserManager) Delete(users *User) {
	ormer := orm.NewOrm()
	ormer.Delete(users, "Id")

}

func (u *UserManager) Lock(users *User) {
	ormer := orm.NewOrm()
	if users.Status == StatusUnLock {
		users.Status = StatusLock
	} else if users.Status == StatusLock {
		users.Status = StatusUnLock
	}
	ormer.Update(users, "Status")
}

type Token struct {
	Id          int        `orm:"column(id)"`
	User        *User      `orm:"column(user);rel(one);"` //user表和token表的反引用与引用的关系，建立一对一的关系
	AccessKey   string     `orm:"column(accesskey);size(1024);"`
	SecretKey   string     `orm:"column(secretkey);size(1024);"`
	CreatedTime *time.Time `orm:"column(created_time);auto_now_add;"`
	UpdatedTime *time.Time `orm:"column(updated_time);auto_now;"`
}

type TokenManager struct{}

func NewTokenManager() *TokenManager {
	return &TokenManager{}
}

func (t *TokenManager) GetByKey(accesskey, secretkey string) *Token {
	token := &Token{AccessKey: accesskey, SecretKey: secretkey}
	ormer := orm.NewOrm()
	if err := ormer.Read(token, "AccessKey", "SecretKey"); err == nil {

		ormer.LoadRelated(token, "User") //通过ormer.LoadRelated来启用token表和user的关联关系
		//token.User.Id 从而可以直接通过token表来访问user表中的内容
		return token
	}

	return nil
}

func (t *TokenManager) GencertToken(users *User) string {
	tokenUser := &Token{User: users}
	ormer := orm.NewOrm()

	if ormer.Read(tokenUser, "User") == orm.ErrNoRows {
		tokenUser.AccessKey = utils.RandString(32)
		tokenUser.SecretKey = utils.RandString(32)
		ormer.Insert(tokenUser)
		return "inserted"

	} else { //先读后更新才能确保数据不被覆盖掉
		tokenUser.AccessKey = utils.RandString(32)
		tokenUser.SecretKey = utils.RandString(32)
		ormer.Update(tokenUser)
		return "updated"
	}
	return ""

}

var DefaultTokenManager = NewTokenManager()
var DefautlUserManager = NewUserManager()

func init() {

	orm.RegisterModel(new(User), new(Token))

}
