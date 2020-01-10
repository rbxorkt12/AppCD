package main

import (
	"fmt"
	"github.com/rbxorkt12/applink/pkg/argocd"
	"github.com/rbxorkt12/applink/pkg/config"
)

func main(){
	config,err:=config.GetConfig("example")
	if err!= nil {panic(err)}
	argoinfo,err:=argocd.ArgocdSet()
	if err!= nil {panic(err)}
	configitem,err:=argocd.GetappsinCluster(*argoinfo)
	if err!= nil {panic(err)}
	clusteritem :=argocd.GetappsinConfig(config)
	create,delete,update:=argocd.Appdiff(configitem,clusteritem)
	fmt.Println("this is create")
	for item1 := range create {
		fmt.Println(item1)
	}
	for item2 := range delete {
		fmt.Println(item2)
	}
	fmt.Println("this is update")
	for item3 := range update {
		fmt.Println(item3)
	}



}