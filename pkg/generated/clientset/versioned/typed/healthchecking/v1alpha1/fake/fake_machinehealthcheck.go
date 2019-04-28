package fake

import (
	v1alpha1 "github.com/openshift/machine-api-operator/pkg/apis/healthchecking/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

type FakeMachineHealthChecks struct {
	Fake	*FakeHealthcheckingV1alpha1
	ns	string
}

var machinehealthchecksResource = schema.GroupVersionResource{Group: "healthchecking.openshift.io", Version: "v1alpha1", Resource: "machinehealthchecks"}
var machinehealthchecksKind = schema.GroupVersionKind{Group: "healthchecking.openshift.io", Version: "v1alpha1", Kind: "MachineHealthCheck"}

func (c *FakeMachineHealthChecks) Get(name string, options v1.GetOptions) (result *v1alpha1.MachineHealthCheck, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewGetAction(machinehealthchecksResource, c.ns, name), &v1alpha1.MachineHealthCheck{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineHealthCheck), err
}
func (c *FakeMachineHealthChecks) List(opts v1.ListOptions) (result *v1alpha1.MachineHealthCheckList, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewListAction(machinehealthchecksResource, machinehealthchecksKind, c.ns, opts), &v1alpha1.MachineHealthCheckList{})
	if obj == nil {
		return nil, err
	}
	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MachineHealthCheckList{ListMeta: obj.(*v1alpha1.MachineHealthCheckList).ListMeta}
	for _, item := range obj.(*v1alpha1.MachineHealthCheckList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}
func (c *FakeMachineHealthChecks) Watch(opts v1.ListOptions) (watch.Interface, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.Fake.InvokesWatch(testing.NewWatchAction(machinehealthchecksResource, c.ns, opts))
}
func (c *FakeMachineHealthChecks) Create(machineHealthCheck *v1alpha1.MachineHealthCheck) (result *v1alpha1.MachineHealthCheck, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewCreateAction(machinehealthchecksResource, c.ns, machineHealthCheck), &v1alpha1.MachineHealthCheck{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineHealthCheck), err
}
func (c *FakeMachineHealthChecks) Update(machineHealthCheck *v1alpha1.MachineHealthCheck) (result *v1alpha1.MachineHealthCheck, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewUpdateAction(machinehealthchecksResource, c.ns, machineHealthCheck), &v1alpha1.MachineHealthCheck{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineHealthCheck), err
}
func (c *FakeMachineHealthChecks) UpdateStatus(machineHealthCheck *v1alpha1.MachineHealthCheck) (*v1alpha1.MachineHealthCheck, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewUpdateSubresourceAction(machinehealthchecksResource, "status", c.ns, machineHealthCheck), &v1alpha1.MachineHealthCheck{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineHealthCheck), err
}
func (c *FakeMachineHealthChecks) Delete(name string, options *v1.DeleteOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := c.Fake.Invokes(testing.NewDeleteAction(machinehealthchecksResource, c.ns, name), &v1alpha1.MachineHealthCheck{})
	return err
}
func (c *FakeMachineHealthChecks) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	action := testing.NewDeleteCollectionAction(machinehealthchecksResource, c.ns, listOptions)
	_, err := c.Fake.Invokes(action, &v1alpha1.MachineHealthCheckList{})
	return err
}
func (c *FakeMachineHealthChecks) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MachineHealthCheck, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewPatchSubresourceAction(machinehealthchecksResource, c.ns, name, pt, data, subresources...), &v1alpha1.MachineHealthCheck{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineHealthCheck), err
}
