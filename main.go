package main

import (
	"encoding/json"
	"fmt"
	"os"

	"io/ioutil"
	"path/filepath"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Source struct {
	KubeConfig string `json:"kubeconfig"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
}

type Version struct {
	ResourceVersion string `json:"resourceVersion"`
}

type Params struct {
	FileName string `json:"fileName"`
}

type Input struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
	Params  Params  `json:"params"`
}

type Output struct {
	Version  Version             `json:"version"`
	Metadata []map[string]string `json:"metadata"`
}

func getKubeClient(kubeconfig string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func cmdCheck(input Input) {
	client, err := getKubeClient(input.Source.KubeConfig)
	if err != nil {
		panic(err)
	}

	configMap, err := client.CoreV1().ConfigMaps(input.Source.Namespace).Get(input.Source.Name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	var versions []Version

	if input.Version.ResourceVersion == "" || input.Version.ResourceVersion != configMap.ObjectMeta.ResourceVersion {
		versions = append(versions, Version{ResourceVersion: configMap.ObjectMeta.ResourceVersion})
	}

	ret, err := json.Marshal(versions)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(ret))
}

func cmdIn(input Input, dir string) {
	client, err := getKubeClient(input.Source.KubeConfig)
	if err != nil {
		panic(err)
	}

	configMap, err := client.CoreV1().ConfigMaps(input.Source.Namespace).Get(input.Source.Name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	// write configMap json into file
	if input.Version.ResourceVersion == configMap.ObjectMeta.ResourceVersion {
		cfg, err := json.Marshal(configMap)
		if err != nil {
			panic(err)
		}

		filePath := filepath.Join(dir, input.Params.FileName)
		if err = os.MkdirAll(dir, 0777); err != nil {
			panic(err)
		}

		if err = ioutil.WriteFile(filePath, cfg, 0777); err != nil {
			panic(err)
		}
	} else {
		panic(fmt.Sprintf("resourceVersion %s != %s", input.Version.ResourceVersion, configMap.ObjectMeta.ResourceVersion))
	}

	output := Output{
		Version: Version{
			ResourceVersion: input.Version.ResourceVersion,
		},
		Metadata: []map[string]string {},
	}

	ret, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(ret))
}

func cmdOut(input Input, dir string) {
	client, err := getKubeClient(input.Source.KubeConfig)
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(dir, input.Params.FileName)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	configMap := &core.ConfigMap{}
	if err := json.Unmarshal(data, configMap); err != nil {
		panic(err)
	}

	_, err = client.CoreV1().ConfigMaps(input.Source.Namespace).Update(configMap)
	if err != nil {
		panic(err)
	}

	output := Output{
		Version: Version{
			ResourceVersion: input.Version.ResourceVersion,
		},
		Metadata: []map[string]string {},
	}

	ret, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(ret))
}

func main() {
	// takes json input from stdin
	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		panic(err)
	}

	// fmt.Println(input, os.Args)

	switch os.Args[1] {
	case "check":
		cmdCheck(input)
	case "in":
		cmdIn(input, os.Args[2])
	case "out":
		cmdOut(input, os.Args[2])
	}
}
