package argocd

import (
	"flag"
	"k8s.io/client-go/rest"
	"os"
	"log"
	"path/filepath"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"errors"
)

var kubeconfig *string



func main() {
	url,err := ArgocdCallurl()
	if err!=nil {
		panic(err)
	}
	fmt.Println(url)
}

func ArgocdCallurl() (string,error){
	argocdport,err := ArgocdServerPortgetter()
	if err!=nil {
		log.Fatalln(err)
		return "",err
	}
	clusterip,err :=K8sclusterIp()
	if err!=nil {
		log.Fatalln(err)
		return "",err
	}
	return fmt.Sprintf("%s:%s",clusterip,argocdport),nil
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

		log.Println("Running out of Kubernetes cluster")
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

func ArgocdServerPortgetter() (string,error){
	con,err:=Connect()
	if err!=nil{
		return "",err
	}
	opts:=metav1.ListOptions{LabelSelector:"name=argocd-server"}
	svcs,err:=con.CoreV1().Services("argocd").List(opts)
	if err!=nil {
		return "",err
	}
	for _,svc := range svcs.Items{
		for _,port := range svc.Spec.Ports{
			if port.Name == "https"{
				return string(port.NodePort),nil
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

		log.Println("Running out of Kubernetes cluster")
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
	return config.Host,err
}
