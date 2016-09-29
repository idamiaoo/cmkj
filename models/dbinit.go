package models

import (
	"go/cmkj_server_go/conf"
	"go/cmkj_server_go/util"

	//"fmt"

	_ "github.com/denisenkom/go-mssqldb" //导入mssql数据库驱动
	"github.com/jinzhu/gorm"
)

type databaseEngine struct {
	Db *gorm.DB
}

//数据库引擎
var (
	DBEngine  *databaseEngine
	HisEngine *databaseEngine
	BetEngine *databaseEngine
)

func connectdb(dburl string) *gorm.DB {
	db, err := gorm.Open("mssql", dburl)
	if err != nil {
		util.Log.Fatal(err)
	}
	if err = db.DB().Ping(); err != nil {
		util.Log.Fatal(err)
	}
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(10)
	db.LogMode(conf.Conf.DefaultBool("orm_logmode", false))

	return db
}

//InitDb 数据库连接初始化
func InitDb() {
	DBEngine = &databaseEngine{
		Db: connectdb(conf.Conf.String("dbServer")),
	}
	HisEngine = &databaseEngine{
		Db: connectdb(conf.Conf.String("hisServer")),
	}
	BetEngine = &databaseEngine{
		Db: connectdb(conf.Conf.String("betServer")),
	}
}
