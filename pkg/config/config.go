package config

import (
	"fmt"
	structtype "github.com/rbxorkt12/applink/pkg/type"
	"github.com/spf13/viper"


)

//argo Rollout 기능은 좀더 고민해봐야 할것 같다.

type Appoconfig struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	Destination string `json:"destination"`
	Charts []Chart `json:"charts"`
}

type Chart struct {
	Repository string `json:"repository"`
	Revision	string `json:"revision"`
	Subpaths []Subpath `json:"subpaths"`
}

type Subpath struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Namespace string `json:"namespace"`
	Chartvalue string `json:"chartvalue"`
	Chartdeploystrategy string `json:"chartdeploystrategy"`
}

func GetConfig(path string,filename string) (*Appoconfig,error) {
	config := &Appoconfig{}
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	if path != "" {
		viper.AddConfigPath(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		return nil,fmt.Errorf("Can't read config file: %s \n", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil,fmt.Errorf("config file format error: %s \n", err)
	}
	return config,nil
}

func (app *Appoconfig)ConvertApp()([]*structtype.Item){
	var list []*structtype.Item
	for _,order:= range app.Orders{
		for _,chart:= range order.Charts{
			for _,subpath := range chart.Subpaths {
				item := &structtype.Item{}
				item.Spec.Dest.Namespace = subpath.Namespace
				item.Spec.Dest.Server=order.Destination
				item.Spec.Source.Revision = chart.Revision
				item.Spec.Source.Path = subpath.Path
				item.Spec.Source.Url = chart.Repository
				item.Spec.Project = "default"
				item.Meta.Name = subpath.Name
				list = append(list, item)
			}
		}
	}
	return list

}
