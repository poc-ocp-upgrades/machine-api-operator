package fake

import (
	v1alpha1 "github.com/openshift/machine-api-operator/pkg/generated/clientset/versioned/typed/healthchecking/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeHealthcheckingV1alpha1 struct{ *testing.Fake }

func (c *FakeHealthcheckingV1alpha1) MachineHealthChecks(namespace string) v1alpha1.MachineHealthCheckInterface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &FakeMachineHealthChecks{c, namespace}
}
func (c *FakeHealthcheckingV1alpha1) RESTClient() rest.Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var ret *rest.RESTClient
	return ret
}
