package main

import (
	"encoding/json"
	"fmt"
	structtype "github.com/rbxorkt12/applink/pkg/type"
	"io/ioutil"
	"log"
	"os"
	"github.com/rbxorkt12/applink/pkg/config"
	"gopkg.in/yaml.v2"
)

func main(){
	firstArg := os.Args[1]
	if firstArg == "convert" {
		config,err:=Readstdinandunmarshalconfig()
		if err!=nil{
			os.Exit(111)
			log.Println(err)
		}
		apps:=config.ConvertApp()
		if os.Args[2]=="auto" {
			for _, app := range apps {
				app.Spec.Sync = &structtype.Syncpolicy{Automated: &structtype.SyncPolicyAutomated{}}
			}
		}
		for _,app := range apps{
			json_byte,err:=json.MarshalIndent(app,"","    ")
			if err!=nil{
				os.Exit(111)
				log.Println(err)}
			fmt.Println(string(json_byte))
		}
	} else {
		log.Println("That is not implemented")
	}
}

func Readstdinandunmarshalconfig() (config.Appoconfig,error){
	data,err:= ioutil.ReadAll(os.Stdin)
	if err!=nil { return config.Appoconfig{}, err}
	config:= config.Appoconfig{}
	err = yaml.Unmarshal(data,&config)
	if err!=nil{ return config,err}
	return config,nil
}