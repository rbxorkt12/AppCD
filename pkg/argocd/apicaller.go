package argocd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	structtype "github.com/rbxorkt12/applink/pkg/type"
	"io/ioutil"
	"log"
	"net/http"
)


func Syncall(items []structtype.Item,argoinfo ArgoCDinfo)error{
	argoport:=argoinfo.iport
	argotoken:=argoinfo.token
	for _,app := range items{
		var revisioncaller map[string]string
		appname:=app.Meta.Name
		syncurl:=fmt.Sprintf("https:%s/api/v1/applications/%s/sync",argoport,appname)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		revisioncaller["revision"]= app.Spec.Source.Revision
		jsonform,err:=json.Marshal(revisioncaller)
		if err!=nil { return err}
		req, err := http.NewRequest("GET", syncurl,bytes.NewBuffer(jsonform) )
		if err != nil {
			return err// handle err
		}
		//request header setting; authorization is required(to get in argocd cluster)
		req.Header.Set("Authorization", "Bearer " +argotoken)
		resp, err := client.Do(req)
		//response error 처리
		if err != nil {
			return err// handle err
		}
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err!=nil { return err}
		log.Printf("Sync: %s의 responce는 %s",appname,string(bytes))
	}
	return nil
}

func Createcall(items []structtype.Item,argoinfo ArgoCDinfo) error {
	argoport := argoinfo.iport
	argotoken := argoinfo.token
	createurl := fmt.Sprintf("http://%s/api/v1/applications", argoport)
	for _, app := range items {
		app.Spec.Project = "default"
		b,err :=json.Marshal(app)
		if err!=nil {return err}
		req, err := http.NewRequest("POST", createurl, bytes.NewBuffer(b))
		if err != nil {
			return err// handle err
		}
		req.Header.Add("Authorization", "Bearer " +argotoken)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err!=nil { return err}
		log.Printf("Create: %s의 responce는 %s",app.Meta.Name,string(bytes))

	}
	return nil
}

func Deletecall(item []structtype.Item,argoinfo ArgoCDinfo) error{
	argoport:=argoinfo.iport
	argotoken:=argoinfo.token
	for _,app:= range item {
		url:=fmt.Sprintf("http://%s/api/v1/applications/%s", argoport, app.Meta.Name)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return err// handle err
		}
		//request header setting; authorization is required(to get in argocd cluster)
		req.Header.Add("Authorization", "Bearer " +argotoken)

		//client gen; insecure connection(without certification)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err!=nil { return err}
		log.Printf("Sync: %s의 responce는 %s",app.Meta.Name,string(bytes))
	}
	return nil

}
