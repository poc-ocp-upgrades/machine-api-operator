package v1alpha1

import (
	"time"
	v1alpha1 "github.com/openshift/machine-api-operator/pkg/apis/healthchecking/v1alpha1"
	scheme "github.com/openshift/machine-api-operator/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

type MachineHealthChecksGetter interface {
	MachineHealthChecks(namespace string) MachineHealthCheckInterface
}
type MachineHealthCheckInterface interface {
	Create(*v1alpha1.MachineHealthCheck) (*v1alpha1.MachineHealthCheck, error)
	Update(*v1alpha1.MachineHealthCheck) (*v1alpha1.MachineHealthCheck, error)
	UpdateStatus(*v1alpha1.MachineHealthCheck) (*v1alpha1.MachineHealthCheck, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.MachineHealthCheck, error)
	List(opts v1.ListOptions) (*v1alpha1.MachineHealthCheckList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MachineHealthCheck, err error)
	MachineHealthCheckExpansion
}
type machineHealthChecks struct {
	client	rest.Interface
	ns	string
}

func newMachineHealthChecks(c *HealthcheckingV1alpha1Client, namespace string) *machineHealthChecks {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &machineHealthChecks{client: c.RESTClient(), ns: namespace}
}
func (c *machineHealthChecks) Get(name string, options v1.GetOptions) (result *v1alpha1.MachineHealthCheck, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result = &v1alpha1.MachineHealthCheck{}
	err = c.client.Get().Namespace(c.ns).Resource("machinehealthchecks").Name(name).VersionedParams(&options, scheme.ParameterCodec).Do().Into(result)
	return
}
func (c *machineHealthChecks) List(opts v1.ListOptions) (result *v1alpha1.MachineHealthCheckList, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.MachineHealthCheckList{}
	err = c.client.Get().Namespace(c.ns).Resource("machinehealthchecks").VersionedParams(&opts, scheme.ParameterCodec).Timeout(timeout).Do().Into(result)
	return
}
func (c *machineHealthChecks) Watch(opts v1.ListOptions) (watch.Interface, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().Namespace(c.ns).Resource("machinehealthchecks").VersionedParams(&opts, scheme.ParameterCodec).Timeout(timeout).Watch()
}
func (c *machineHealthChecks) Create(machineHealthCheck *v1alpha1.MachineHealthCheck) (result *v1alpha1.MachineHealthCheck, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result = &v1alpha1.MachineHealthCheck{}
	err = c.client.Post().Namespace(c.ns).Resource("machinehealthchecks").Body(machineHealthCheck).Do().Into(result)
	return
}
func (c *machineHealthChecks) Update(machineHealthCheck *v1alpha1.MachineHealthCheck) (result *v1alpha1.MachineHealthCheck, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result = &v1alpha1.MachineHealthCheck{}
	err = c.client.Put().Namespace(c.ns).Resource("machinehealthchecks").Name(machineHealthCheck.Name).Body(machineHealthCheck).Do().Into(result)
	return
}
func (c *machineHealthChecks) UpdateStatus(machineHealthCheck *v1alpha1.MachineHealthCheck) (result *v1alpha1.MachineHealthCheck, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result = &v1alpha1.MachineHealthCheck{}
	err = c.client.Put().Namespace(c.ns).Resource("machinehealthchecks").Name(machineHealthCheck.Name).SubResource("status").Body(machineHealthCheck).Do().Into(result)
	return
}
func (c *machineHealthChecks) Delete(name string, options *v1.DeleteOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.client.Delete().Namespace(c.ns).Resource("machinehealthchecks").Name(name).Body(options).Do().Error()
}
func (c *machineHealthChecks) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().Namespace(c.ns).Resource("machinehealthchecks").VersionedParams(&listOptions, scheme.ParameterCodec).Timeout(timeout).Body(options).Do().Error()
}
func (c *machineHealthChecks) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MachineHealthCheck, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result = &v1alpha1.MachineHealthCheck{}
	err = c.client.Patch(pt).Namespace(c.ns).Resource("machinehealthchecks").SubResource(subresources...).Name(name).Body(data).Do().Into(result)
	return
}
