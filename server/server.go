package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rbxorkt12/applink/pkg/argocd"
	"github.com/rbxorkt12/applink/pkg/config"
	"github.com/rbxorkt12/applink/pkg/handler"
	"github.com/rbxorkt12/applink/pkg/gitchecker"
	"github.com/rbxorkt12/applink/pkg/lint"
	"github.com/rbxorkt12/applink/pkg/slack"
	"log"
	"os"
	"sync"
	"time"
)

// Const를 나중에 다른방법으로 넣을수있게 바꿔주세요.
const (
	Giturl= "https://github.com/rbxorkt12/appcd_example.git"
	Dir = "./example"
	SettingDuration time.Duration = time.Second*100
	Repoid string = ""
	Password string = ""
)

func main() {
	//나중에는 이 부분만 프로그래밍되고 밑부분은 전부다 파이프라인
//	e := echo.New()
//	e.Use(middleware.Logger())
//	e.Use(middleware.Recover())
//	optionhanlder:=&handler.Optionhandler{Change:handler.Wait}
//	optionhanlder.Sethandler(e.Group("/api"))
//	e.Logger.Debug(e.Start(":/8080"))
	//handler setting

	//1. argocd info setting
	argoinfo,err:=argocd.ArgocdSet()
	if err!=nil {
		log.Fatalln(err)
		os.Exit(128)
	}

	//2. repo info 찾기 in k8s secret

	for {

		//3. git이 update 되었는지 확
		time.Sleep(SettingDuration)
		flag,err:=gitchecker.Isrepotobeupdate(Giturl,Dir)
		if(err!=nil){
			if(err.Error()=="NOFILEEXIST"){
				_,err:=gitchecker.AuthGitclone(Giturl,Dir,Repoid,Password)
				if err!=nil {
					log.Fatalln(err)
					os.Exit(555)
				}
			} else {
			log.Fatalln(err)
			os.Exit(238) //내맘임
			}
		}
		if flag == false {
			continue
		}
		// 3-2 update 되었으면 최근커밋을 불러옴.
		err = gitchecker.Gitupdate(Giturl,Dir,Repoid,Password)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(238)
		}
		// 3-3 불러온 최근 커밋의 config를 lint.
		flag,err = lint.Configvalid(Dir)
		if flag==false {
			log.Println("It is not invalid Appoconfig")
			continue
		}

		slack.SendSlackNotification("webhookurl","Do you want change?")
		autocon, err := config.GetConfig(Dir,"Autoconfig")
		manualconf, err := config.GetConfig(Dir,"Manualconfig")
		if err != nil {
			log.Fatalln(err)
			os.Exit(577)
		}
		autoapplist := argocd.GetappsinConfig(autocon)
		manualapplist := argocd.GetappsinConfig(manualconf)
		_,_,sameapp := argocd.Appdiff(autoapplist,manualapplist)
		if sameapp == nil {
			log.Fatalln("Duplicated application between Autoconfig.yaml and Manualconfig.yaml")
			continue
		}
		clusterlist, err := argocd.GetappsinCluster(*argoinfo)
		create_auto, delete1, update1 := argocd.Appdiff(clusterlist, autoapplist)
		create_manual, delete2, update2 := argocd.Appdiff(clusterlist, manualapplist)

		for _,app:= range create_auto{
			app.Spec.Project = "default"
			app.Spec.Sync = "Auto"
		}

		for _,app:= range create_manual{
			app.Spec.Project = "default"
			app.Spec.Sync = "manual"
		}

		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			argocd.Createcall(create_auto,*argoinfo)
			argocd.Createcall(create_manual,*argoinfo)
			defer wg.Done()
		}()
		go func() {
			argocd.Deletecall(delete1,*argoinfo)
			argocd.Deletecall(delete2,*argoinfo)
			defer wg.Done()
		}()
		go func() {
			argocd.Updatecall(update1,*argoinfo)
			argocd.Updatecall(update2,*argoinfo)
			defer wg.Done()
		}()
		wg.Wait()
		log.Println("All setting succeed")
		slack.SendSlackNotification("webhookurl", "All setting succeed")
	}
}


