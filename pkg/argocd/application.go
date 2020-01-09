package argocd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/rbxorkt12/applink/pkg/config"
	"go/ast"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)


type Reciver struct{
	Items []Item `json:"items"`
}

type Item struct {
	Meta Metadata `json:"metadata"`
	Spec Spec	`json"spec"`
}

type Metadata struct {
	Name string `json:"name"`
}

type Spec struct {
	Source Source `json:"source"`
	Dest Destination `json:"destination"`
}

type Source struct {
	Url string `json:"repoURL"`
	Path string `json:"path"`
	Revision string `json:"targetRevision"`
}

type Destination struct {
	Server	string `json:"server"`
	Namespace string `json:"namespace"`
}

//diff 많이 느릴겁니다. 알고리즘 고쳐주세용


func GetappsinConfig ( config *config.Appoconfig) []Item{
	return config.ConvertApp()
}

func GetappsinCluster( cluster argoCDinfo ) ([]Item,error){
	//appname := "appsync"
	url:=fmt.Sprintf("http://%s/api/v1/applications", cluster.iport) //API calling for get application list, not completed
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//request generating to get application info, http method is GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil,err// handle err
	}
	//request header setting; authorization is required(to get in argocd cluster)
	req.Header.Set("Authorization", "Bearer " +cluster.token)
	resp, err := client.Do(req)
	if err != nil {
		return nil,err// handle err
	}
	defer resp.Body.Close()
	//data reading from response
	bytes, _ := ioutil.ReadAll(resp.Body)
	var dat *Reciver
	if err := json.Unmarshal(bytes, dat); err != nil {
		return nil,err// handle err
	}
	return dat.Items,nil

	//fmt.Println(string(y))
	//strresp := string(bytes) //바이트를 문자열로
	//fmt.Println(strresp)
}


