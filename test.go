package main

import (
	"fmt"
	"github.com/rbxorkt12/applink/pkg/argocd"
)

func main(){
	argoinfo,err:=argocd.ArgocdSet()
	if err!=nil {panic(err)}
	fmt.Println(argoinfo)

}