package lint

import (
	"errors"
	"fmt"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"log"
	"os"
	config "github.com/rbxorkt12/applink/pkg/config"
)


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