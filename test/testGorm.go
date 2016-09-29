// testGorm
package main

import (
	"go/cmkj_server_go/conf"
	"go/cmkj_server_go/models"

	"fmt"

	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mssql"
)

func main() {
	conf.Loadconf("test.conf")
	models.InitDb()
	models.DBEngine.Db.LogMode(true)
	tep := models.ReadOrders(11, 3, 15, 84)
	if tep != nil {
		fmt.Println(tep)
	}

}
