package argocd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/rbxorkt12/applink/pkg/config"
	structtype "github.com/rbxorkt12/applink/pkg/type"
	"io/ioutil"
	"net/http"
)


func GetappsinConfig ( config *config.Appoconfig) []structtype.Item{
	return config.ConvertApp()

}

func GetappsinCluster( cluster ArgoCDinfo ) ([]structtype.Item,error){
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
	var dat *structtype.Reciver
	if err := json.Unmarshal(bytes, dat); err != nil {
		return nil,err// handle err
	}
	return dat.Items,nil

	//fmt.Println(string(y))
	//strresp := string(bytes) //바이트를 문자열로
	//fmt.Println(strresp)
}


