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
	fmt.Println("this is clusteritem")
	for _,item1:= range clusteritem{
		fmt.Println(item1)
	}
	fmt.Println("this is configitem")
	for _,item2:= range configitem{
		fmt.Println(item2)
	}


}