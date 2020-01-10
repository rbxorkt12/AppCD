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
	Giturl= "https://github.com/example.git"
	Dir = "./example"
	SettingDuration time.Duration = time.Second*100
	Repoid string = "admin"
	Password string = "password"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	optionhanlder:=&handler.Optionhandler{Change:handler.Wait}
	optionhanlder.Sethandler(e.Group("/api"))
	e.Logger.Debug(e.Start(":/8080"))
	//handler setting

	argoinfo,err:=argocd.ArgocdSet()
	if err!=nil {
		log.Fatalln(err)
		os.Exit(128)
	}
	//argocd info setting

	for {
		time.Sleep(SettingDuration)
		flag,err:=gitchecker.Isrepotobeupdate(Giturl,Dir)
		if(err!=nil){
			if(err.Error()=="NOFILEEXIST"){
				gitchecker.AuthGitclone(Giturl,Dir,Repoid,Password)
			} else {
			log.Fatalln(err)
			os.Exit(238) //내맘임
			}
		}
		if flag == false {
			continue
		}
		// repo가 업데이트 된 상황
		err = gitchecker.Gitupdate(Giturl,Dir,Repoid,Password)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(238)
		}
		// 최근 커밋으로 바꿈.
		flag,err = lint.Configvalid(Dir)
		if flag==false {
			log.Println("It is not invalid Appoconfig")
			continue
		}
		// appoconfig가 valid?
		slack.SendSlackNotification("webhookurl","Do you want change?")
		for {
			if optionhanlder.Change == handler.Wait {
				time.Sleep(time.Second)
			} else if optionhanlder.Change == handler.No {
				optionhanlder.Change = handler.Wait
				log.Println("You choice not change")
				slack.SendSlackNotification("webhookurl", "You choice not change")
				break
			} else if optionhanlder.Change == handler.YES {
				optionhanlder.Change = handler.Wait
				log.Println("You choice change")
				slack.SendSlackNotification("webhookurl", "You choice change")
				con, err := config.GetConfig(Dir)
				if err != nil {
					log.Fatalln(err)
					os.Exit(577)
				}
				applist := argocd.GetappsinConfig(con)
				clusterlist, err := argocd.GetappsinCluster(*argoinfo)
				create, delete, update := argocd.Appdiff(clusterlist, applist)
				var wg sync.WaitGroup
				wg.Add(3)
				go func() {
					argocd.Createcall(create)
					defer wg.Done()
				}()
				go func() {
					argocd.Deletecall(delete)
					defer wg.Done()
				}()
				go func() {
					argocd.Updatecall(update)
					defer wg.Done()
				}()
				wg.Wait()
				argocd.Syncall()
				log.Println("All setting succeed")
				slack.SendSlackNotification("webhookurl", "All setting succeed")
				break
			} else {
				log.Fatalln("Option controller are not error")
				os.Exit(788)
			}
		}
	}
}


