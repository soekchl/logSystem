package models

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	m_orm     orm.Ormer
	slave_orm orm.Ormer
)

func init() {
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModel(
		new(User),
		new(Log),
	)
	orm.RunSyncdb("default", false, true)

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	m_orm = orm.NewOrm()

	initSlave()
}

func initSlave() {
	dbhost := beego.AppConfig.String("db.slave.host")
	dbport := beego.AppConfig.String("db.slave.port")
	dbuser := beego.AppConfig.String("db.slave.user")
	dbpassword := beego.AppConfig.String("db.slave.password")
	dbname := beego.AppConfig.String("db.slave.name")
	timezone := beego.AppConfig.String("db.slave.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	orm.RegisterDataBase("slave", "mysql", dsn)

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	slave_orm = orm.NewOrm()
}

func TableName(name string) string {
	return name
}
