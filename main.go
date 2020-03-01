package main

import (
	"flag"
	"fmt"

	yaml "gopkg.in/yaml.v2"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type NetworkEntry struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
	Type   string `yaml:"type`
}

type StorageEntry struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type AffinityEntry struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
	Policy string `yaml:"policy"`
}

type Mapping struct {
	NetworkMapping  []NetworkEntry  `yaml:"networkMapping"`
	StorageMapping  []StorageEntry  `yaml:"storageMapping"`
	AffinityMapping []AffinityEntry `yaml:"affinityMapping"`
}

const (
	kubeconfigPath = "/home/masayag/.crc/cache/crc_libvirt_4.2.13/kubeconfig"
	ovirtMapping   = "ovirt-mapping-example"
	namespace      = "default"
)

func main() {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", kubeconfigPath, "absolute path to the kubeconfig file")

	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var confmap *k8sv1.ConfigMap
	confmap, err = clientset.CoreV1().ConfigMaps(namespace).Get(ovirtMapping, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("ConfigMap %s in namespace %s not found\n", ovirtMapping, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting configmap %s in namespace %s: %v\n",
			ovirtMapping, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found configmap %s in namespace %s\n", ovirtMapping, namespace)
	}

	configyaml := confmap.Data["mappings"]
	out := Mapping{}
	if err := yaml.Unmarshal([]byte(configyaml), &out); err != nil {
		panic(err)
	}

	fmt.Println("Network Mappings:")
	for _, element := range out.NetworkMapping {
		fmt.Printf("source: %s, target: %s, type: %s\n", element.Source, element.Target, element.Type)
	}

	fmt.Println("Storage Mapping:")
	for _, element := range out.StorageMapping {
		fmt.Printf("source: %s, target: %s\n", element.Source, element.Target)
	}

	fmt.Println("Affinity Mappings:")
	for _, element := range out.AffinityMapping {
		fmt.Printf("source: %s, target: %s, policy: %s\n", element.Source, element.Target, element.Policy)
	}
}
