package main

import(
	"fmt"
	"github.com/rbxorkt12/applink/pkg/argocd"
	"log"
)

func main(){
	url,err:=argocd.ArgocdCallurl()
	if err!=nil {
		log.Println(url)
	}
	fmt.Println(url)
}