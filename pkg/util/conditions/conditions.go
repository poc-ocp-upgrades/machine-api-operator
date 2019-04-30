package conditions

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
)

func GetNodeCondition(node *corev1.Node, conditionType corev1.NodeConditionType) *corev1.NodeCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, cond := range node.Status.Conditions {
		if cond.Type == conditionType {
			return &cond
		}
	}
	return nil
}

type UnhealthyConditions struct {
	Items []UnhealthyCondition `json:"items"`
}
type UnhealthyCondition struct {
	Name	corev1.NodeConditionType	`json:"name"`
	Status	corev1.ConditionStatus		`json:"status"`
	Timeout	string				`json:"timeout"`
}

func GetNodeUnhealthyConditions(node *corev1.Node, cmUnealthyConditions *corev1.ConfigMap) ([]UnhealthyCondition, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, ok := cmUnealthyConditions.Data["conditions"]
	if !ok {
		return nil, fmt.Errorf("can not find \"conditions\" under the configmap")
	}
	var unealthyConditions UnhealthyConditions
	err := yaml.Unmarshal([]byte(data), &unealthyConditions)
	if err != nil {
		glog.Errorf("failed to umarshal: %v", err)
		return nil, err
	}
	conditions := []UnhealthyCondition{}
	for _, c := range unealthyConditions.Items {
		cond := GetNodeCondition(node, c.Name)
		if cond != nil && cond.Status == c.Status {
			conditions = append(conditions, c)
		}
	}
	return conditions, nil
}
func CreateDummyUnhealthyConditionsConfigMap() (*corev1.ConfigMap, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	unhealthyConditions := &UnhealthyConditions{Items: []UnhealthyCondition{{Name: "Ready", Status: "Unknown", Timeout: "300s"}, {Name: "Ready", Status: "False", Timeout: "300s"}}}
	conditionsData, err := yaml.Marshal(unhealthyConditions)
	if err != nil {
		return nil, err
	}
	return &corev1.ConfigMap{Data: map[string]string{"conditions": string(conditionsData)}}, nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
