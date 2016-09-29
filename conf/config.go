package conf

import (
	"fmt"

	"github.com/astaxie/beego/config"
)

var Conf config.Configer

func Loadconf(configfile string) {
	var err error
	Conf, err = config.NewConfig("ini", configfile)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
func GetSectionkey(section, key string) string {
	return fmt.Sprintf("%s::%s", section, key)
}
