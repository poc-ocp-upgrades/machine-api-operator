package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"github.com/openshift/machine-api-operator/pkg/apis/healthchecking"
)

var SchemeGroupVersion = schema.GroupVersion{Group: healthchecking.GroupName, Version: "v1alpha1"}

func Kind(kind string) schema.GroupKind {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}
func Resource(resource string) schema.GroupResource {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder	= runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme		= SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	scheme.AddKnownTypes(SchemeGroupVersion, &MachineHealthCheck{}, &MachineHealthCheckList{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
