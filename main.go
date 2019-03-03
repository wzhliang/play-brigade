package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Azure/brigade/pkg/brigade"
	"github.com/Azure/brigade/pkg/storage/kube"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Test struct {
	Name string
	Type string
}

func main() {
	var ns string
	flag.StringVar(&ns, "namespace", "", "namespace")

	// Bootstrap k8s configuration from local 	Kubernetes config file
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config.backup")
	log.Println("Using kubeconfig file: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create an rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	tst := &Test{
		Name: "mytest",
		Type: "shell",
	}

	payload, _ := json.Marshal(tst)
	js, err := ioutil.ReadFile("brigade.js")
	if err != nil {
		log.Fatal(err)
	}
	b := &brigade.Build{
		ProjectID: "brigade-68d2c7440da7da85970d5abf22c2fd2eea6239e67cfca22a9766c1",
		Type:      "wisecloud/test",
		Provider:  "wliang",
		Revision:  &brigade.Revision{Ref: "refs/heads/master"},
		Payload:   payload,
		Script:    js,
	}

	store := kube.New(clientset, "default")
	err = store.CreateBuild(b)
	if err != nil {
		log.Fatal(err)
	}
}
