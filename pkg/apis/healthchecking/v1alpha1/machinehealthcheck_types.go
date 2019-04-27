package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ConfigMapNodeUnhealthyConditions = "node-unhealthy-conditions"

type MachineHealthCheck struct {
	metav1.TypeMeta		`json:",inline"`
	metav1.ObjectMeta	`json:"metadata,omitempty"`
	Spec			MachineHealthCheckSpec		`json:"spec,omitempty"`
	Status			MachineHealthCheckStatus	`json:"status,omitempty"`
}
type MachineHealthCheckList struct {
	metav1.TypeMeta	`json:",inline"`
	metav1.ListMeta	`json:"metadata,omitempty"`
	Items		[]MachineHealthCheck	`json:"items"`
}
type MachineHealthCheckSpec struct {
	Selector metav1.LabelSelector `json:"selector"`
}
type MachineHealthCheckStatus struct{}
