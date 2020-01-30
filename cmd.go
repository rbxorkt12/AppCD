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
	"io"
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
		var autolist []structtype.Item
		var manuallist []structtype.Item
		for _,item:= range recieve.Items{
			if item.Meta.Annotations.AppCDoption == "Auto"  {
				if err != nil {
					log.Fatalln(err)
					os.Exit(23)
				}
				autolist=append(autolist,item)
			}
			if item.Meta.Annotations.AppCDoption == "Manual"{
				if err!=nil {
					log.Fatalln(err)
					os.Exit(23)
				}
				manuallist=append(manuallist,item)
			}
		}
		if a== "Auto"{
			byte,err:=json.MarshalIndent(autolist,"","   ")
                        if err!=nil {
               	                log.Fatalln(err)
                                os.Exit(23)
                                }
				fmt.Println(string(byte))
			}else if a=="Manual"{
				byte,err:=json.MarshalIndent(manuallist,"","   ")
                                if err!=nil {
                                        log.Fatalln(err)
                                        os.Exit(23)
                                }
				fmt.Println(string(byte))
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
		fi,err:=os.Stat(before)
		len:=fi.Size()
		fi2,err:=os.Stat(after)
		len2:=fi2.Size()
		dat, err := ioutil.ReadFile(before)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		var beforeitems []structtype.Item
                dat2, err := ioutil.ReadFile(after)
                var afteritems []structtype.Item
		if len<=10{
			fmt.Println("sss")
			if len2<=10{
				fmt.Println("both emptry")
				os.Exit(123)
			}else{
				err=json.Unmarshal(dat2,&afteritems)
				writeitems(afteritems,"/diff/CREATE")
				os.Create("/diff/DELETE")
				os.Create("/diff/UPDATE")
				return
			}
		}else{
			if len2<=10{
				err=json.Unmarshal(dat,&beforeitems)
				writeitems(beforeitems,"/diff/DELETE")
				os.Create("/diff/CREATE")
                                os.Create("/diff/UPDATE")
				return
			}
		}
		err=json.Unmarshal(dat,&beforeitems)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		err=json.Unmarshal(dat2,&afteritems)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		create,delete,update:=argocd.Appdiff(beforeitems,afteritems)
		writeitems(create,"/diff/CREATE")
		writeitems(delete,"/diff/DELETE")
		writeitems(update,"/diff/UPDATE")

	case "appstoparam":
		var resultlist []map[string]string
		applist,err:=ReadStdinAndUnmarshalApp()
		if err!=nil {
			log.Fatalln(err)
			os.Exit(233)
		}
		for _,app := range applist{
			result := make(map[string]string)
			result["appcdoption"]=app.Meta.Annotations.AppCDoption
			result["name"]=app.Meta.Name
			result["project"]=app.Spec.Project
			result["path"]=app.Spec.Source.Path
			result["revision"]=app.Spec.Source.Revision
			result["url"]=app.Spec.Source.Url
			result["namespace"]=app.Spec.Dest.Namespace
			result["server"]=app.Spec.Dest.Server
			resultlist= append(resultlist,result )
		}
		byte,err:=json.MarshalIndent(resultlist,"","   ")
		if err!=nil {
			log.Fatalln(err)
			os.Exit(233)
		}
		fmt.Println(string(byte))


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
	if err!=nil {return nil,err }
	if string(data) =="\n"{
		os.Exit(123)
	}
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

func ReadStdinAndUnmarshalReciver() (structtype.Reciver,error){
	data,err:= ioutil.ReadAll(os.Stdin)
	if err!=nil {
		return structtype.Reciver{},err
	}
	var reciver structtype.Reciver
	err=json.Unmarshal(data,&reciver)
	if err!=nil {
		return structtype.Reciver{},err
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
	byte,err:=json.MarshalIndent(qq,"","   ")
	if err!=nil {
		log.Fatalln(err)
		os.Exit(123)
	}
	f.WriteString(string(byte))
	
}
func IsDirEmpty(name string) bool {
        f, err := os.Open(name)
        if err != nil {
                return false
        }
        defer f.Close()

        // read in ONLY one file
        _, err = f.Readdir(1)

        // and if the file is EOF... well, the dir is empty.
        if err == io.EOF {
                return true
        }
        return false
}
