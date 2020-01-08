package lint

import (
	"errors"
	"fmt"
	"github.com/rbxorkt12/applink/pkg/config"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"log"
	"os"
	"strings"
)

//각종 lint function들은 함수로 만들고 이 함수안에 저장해주세요.
func Configvalid(directory string) (bool,error){
	appoconfig,err:=config.GetConfig(directory)
	if (err != nil){
		return false,err
	}
	orders:=appoconfig.Orders
	for _,order:= range orders{
		for _,chart:= range order.Charts{
			if strings.HasSuffix(chart.Repository,"git"){
				return false,errors.New("There is invalid value in Appoconfig, especially in repository value")
			}
		}
	}
	return true,nil
}

func Appoconfigexist(directory string) (bool,error){
	appopath := fmt.Sprintf("%s/Appoconfig.yaml",directory)
	if!(fileExists(appopath)){
		log.Fatalln("There is no Appoconfig file in directory, Please make it")
		return false,errors.New("no Appoconfig")
	}
	Info("Found Appoconfig.yaml")
	return true,nil
}


func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}