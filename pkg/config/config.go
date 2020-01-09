package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/rbxorkt12/applink/pkg/argocd"
)

//argo Rollout 기능은 좀더 고민해봐야 할것 같다.

type Appoconfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Orders []Order `json:"orders"`
}

type Order struct {
	Destination string `json:"destination"`
	Charts []Chart `json:"charts"`
}

type Chart struct {
	Repository string `json:"repository"`
	Path 	string `json:"path"`
	Revision	string `json:"revision"`
	Namespace string `json:"namespace"`
}

func GetConfig(path string) (*Appoconfig,error) {
	config := &Appoconfig{}
	viper.SetConfigName("Appoconfig")
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

func (app *Appoconfig)ConvertApp()([]argocd.Item){
	var list []argocd.Item
	for _,order:= range app.Orders{
		for _,chart:= range order.Charts{
			item:=&argocd.Item{}
			item.Spec.Dest.Namespace=chart.Namespace
			item.Spec.Dest.Server=order.Destination
			item.Spec.Source.Revision=chart.Revision
			item.Spec.Source.Path=chart.Path
			item.Spec.Source.Url=chart.Repository
			item.Meta.Name=fmt.Sprintf("%s_%s_%s_%s",order.Destination,chart.Namespace,chart.Repository,chart.Path)
			list = append(list, *item)
		}
	}
	return list
}