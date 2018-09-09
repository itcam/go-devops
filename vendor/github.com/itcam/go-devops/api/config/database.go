package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

type DBPool struct {
	Uic *gorm.DB
}

var (
	dbp DBPool
)

func Con() DBPool {
	return dbp
}

func SetLogLevel(loggerlevel bool) {
	dbp.Uic.LogMode(loggerlevel)
}

func InitDB(loggerlevel bool, vip *viper.Viper) (err error) {

	var u *sql.DB
	uicd, err := gorm.Open("mysql", vip.GetString("db.bihu"))
	uicd.Dialect().SetDB(u)
	uicd.LogMode(loggerlevel)
	if err != nil {
		return fmt.Errorf("connect to bihu: %s", err.Error())
	}
	uicd.SingularTable(true)
	dbp.Uic = uicd

	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return "arrow_" + defaultTableName
	//}
	SetLogLevel(loggerlevel)
	return
}

func CloseDB() (err error) {

	err = dbp.Uic.Close()
	if err != nil {
		return
	}
	return
}
