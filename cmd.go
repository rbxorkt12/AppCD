package main

import (
	"encoding/json"
	"fmt"
	"github.com/rbxorkt12/applink/pkg/argocd"
	structtype "github.com/rbxorkt12/applink/pkg/type"
	"io/ioutil"
	"log"
	"os"
	"github.com/rbxorkt12/applink/pkg/config"
	"gopkg.in/yaml.v2"
)


type Namelist struct{
	Name string `json:"name"`
}

func main(){
	firstArg := os.Args[1]
	switch firstArg {
	case "convert" :
		config,err:=Readstdinandunmarshalconfig()
		if err!=nil{
			os.Exit(111)
			log.Println(err)
		}
		apps:=config.ConvertApp()
		for _, app := range apps {
			if os.Args[2]=="auto" {
				app.Meta.Annotations.AppCDoption = "Auto"
				app.Spec.Sync = &structtype.Syncpolicy{Automated: &structtype.SyncPolicyAutomated{}}
			} else {
				app.Meta.Annotations.AppCDoption = "Manual"
			}
		}
		json_byte,err:=json.MarshalIndent(apps,"","    ")
		if err!=nil{
			os.Exit(111)
			log.Println(err)}
		fmt.Println(string(json_byte))
	case "names" :
		apps,err:=ReadStdinAndUnmarshalApp()
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		var strings []Namelist
		for _,app:= range apps{
			var name Namelist
			name.Name = app.Meta.Name
			strings=append(strings,name)
		}
		json_byte,err:=json.MarshalIndent(strings,"","    ")
		if err!=nil{
			os.Exit(111)
			log.Println(err)}
		fmt.Println(string(json_byte))
	case "argoinfo" :
		argoinfo,err:=argocd.ArgocdSet(os.Args[2],os.Args[3])
		if err!=nil {
			log.Fatalln(err)
			os.Exit(111)
		}
		byte,err:=json.Marshal(argoinfo)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(111)
		}
		fmt.Println(string(byte))
	case "find" :
		target:= os.Args[2]
		apps,err:=ReadStdinAndUnmarshalApp()
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		for _,app := range apps {
			if app.Meta.Name == target{
				jsonbody,err:=json.MarshalIndent(app,"","    ")
				if err!=nil {
					log.Fatalln(err)
					os.Exit(23)
				}
				fmt.Println(string(jsonbody))
			}
		}
	case "split" :
		a := os.Args[2]
		recieve,err:=ReadStdinAndUnmarshalReciver()
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		for _,item:= range recieve.Items{
			if item.Meta.Annotations.AppCDoption == "Auto" && a== "Auto" {
				onebody, err := json.MarshalIndent(item, "", "   ")
				if err != nil {
					log.Fatalln(err)
					os.Exit(23)
				}
				fmt.Println(string(onebody))
			}
			if item.Meta.Annotations.AppCDoption == "Manual" && a== "Manual"{
				onebody,err:=json.MarshalIndent(item,"","   ")
				if err!=nil {
					log.Fatalln(err)
					os.Exit(23)
				}
				fmt.Println(string(onebody))
			}
		}
	case "diff":
		before := os.Args[2]
		after := os.Args[3]
		flag:=Exists(before)
		if flag == false {
			fmt.Errorf("no file %s",before)
			os.Exit(123)
		}
		flag=Exists(after)
		if flag == false {
			fmt.Errorf("no file %s",after)
			os.Exit(123)
		}
		dat, err := ioutil.ReadFile(before)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		var beforeitems []structtype.Item
		err=json.Unmarshal(dat,&beforeitems)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		dat2, err := ioutil.ReadFile(after)
		var afteritems []structtype.Item
		err=json.Unmarshal(dat2,&afteritems)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		create,delete,update:=argocd.Appdiff(beforeitems,afteritems)
		writeitems(create,"/diff/CREATE")
		writeitems(delete,"/diff/DELETE")
		writeitems(update,"/diff/UPDATE")

	default :
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

func ReadStdinAndUnmarshalApp() ([]*structtype.Item,error){
	data,err:= ioutil.ReadAll(os.Stdin)
	if err!=nil {
		return nil,err
	}
	var apps []*structtype.Item
	err=json.Unmarshal(data,&apps)
	if err!=nil {
		return nil,err
	}
	return apps,nil
}

func ReadStdinAndUnmarshalReciver() (*structtype.Reciver,error){
	data,err:= ioutil.ReadAll(os.Stdin)
	if err!=nil {
		return nil,err
	}
	var reciver *structtype.Reciver
	err=json.Unmarshal(data,reciver)
	if err!=nil {
		return nil,err
	}
	return reciver,nil
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func writeitems(qq []structtype.Item,dest string){
	f, err := os.Create(dest)
	defer f.Close()
	if err!=nil {
		log.Fatalln(err)
		os.Exit(123)
	}
	for _,item:=range qq{
		byte,err:=json.MarshalIndent(item,"","   ")
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		f.WriteString(string(byte))
	}
}