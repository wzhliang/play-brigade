package main

import (
	"encoding/json"
	"flag"
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
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
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
	b := &brigade.Build{
		ProjectID: "abcdefg",
		Type:      "idon'tknow",
		Provider:  "wliang",
		Revision:  &brigade.Revision{Ref: "refs/heads/master"},
		Payload:   payload,
	}

	store := kube.New(clientset, "default")
	err = store.CreateBuild(b)
	if err != nil {
		log.Fatal(err)
	}
}
