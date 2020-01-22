package argocd

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"flag"
	"k8s.io/client-go/rest"
	"os"
	"log"
	"path/filepath"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"encoding/json"
	"fmt"
	"errors"
)

var kubeconfig *string

type ArgoCDinfo struct {
	Username string `json:"username"`
	Password string	`json:"password"`
	Iport string `json:"iport"`
	Token string `json:"token"`
}

//id,password 짜는 알고리즘

func ArgocdSet(id string, password string) (*ArgoCDinfo,error){
	var argoinfo ArgoCDinfo
	argoinfo.Password=password
	argoinfo.Username=id
	url,err:=ArgocdCallurl()
	if err!=nil {
		panic(err)
		return nil,err
	}
	argoinfo.Iport = url
	gettoken(&argoinfo)
	return &argoinfo,nil

}

func gettoken(cluster *ArgoCDinfo) string{
	//인증서 없이 접근
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//accountmap default id/pw setting: admin/password
	accountmap := map[string] string{"username" : cluster.Username, "password" : cluster.Password}
	tokenmap := map[string]string{"token":"None"}
	bodyjson, _ :=json.Marshal(accountmap)
	url:=fmt.Sprintf("http://%s/api/v1/session", cluster.Iport)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyjson))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	tokenbytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(tokenbytes, &tokenmap)
	if err != nil {
		panic(err)
	}
	cluster.Token = tokenmap["token"]
	return cluster.Token
}

func ArgocdCallurl() (string,error){
	argocdport,err := ArgocdNodePortgetter()
	if err!=nil {
		log.Fatalln(err)
		return "",err
	}
	clusterip,err :=K8sclusterIp()
	if err!=nil {
		log.Fatalln(err)
		return "",err
	}
	url:= fmt.Sprintf("%s:%s",clusterip,argocdport)
	return url,nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func Connect() (*kubernetes.Clientset, error) {
	var config, err = rest.InClusterConfig()
	if err == nil {
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return clientset, nil
	} else {
		if kubeconfig == nil {
			if home := homeDir(); home != "" {
				kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
			flag.Parse()
		}

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
			return nil, err
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
			return nil, err
		}
		return clientset, nil
	}
}

func ArgocdNodePortgetter() (string,error){
	con,err:=Connect()
	if err!=nil{
		return "",err
	}
	opts:=metav1.ListOptions{LabelSelector:"app.kubernetes.io/component=server"}
	svcs,err:=con.CoreV1().Services("argocd").List(opts)
	if err!=nil {
		return "",err
	}
	for _,svc := range svcs.Items{
		for _,port := range svc.Spec.Ports{
			if port.Name == "https"{
				portnum := fmt.Sprintf("%d",port.NodePort)
				return portnum,nil
			}
		}

	}
	return "",errors.New("There is no argocd-server service in k8s cluster")
}

func ClusterConfig() (*rest.Config,error){
	var config, err = rest.InClusterConfig()
	if err == nil {
		return config, nil
	} else {
		if kubeconfig == nil {
			if home := homeDir(); home != "" {
				kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
			flag.Parse()
		}

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
			return nil, err
		}
	}
	return config,nil
}

func K8sclusterIp() (string,error){
	config,err:=ClusterConfig()
	strs:=strings.Split(config.Host,":")
	result:=strings.TrimPrefix(strs[1],"//")
	url:=fmt.Sprintf("%s",result)
	return url,err
}
