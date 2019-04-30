package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func (in *MachineHealthCheck) DeepCopyInto(out *MachineHealthCheck) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}
func (in *MachineHealthCheck) DeepCopy() *MachineHealthCheck {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(MachineHealthCheck)
	in.DeepCopyInto(out)
	return out
}
func (in *MachineHealthCheck) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *MachineHealthCheckList) DeepCopyInto(out *MachineHealthCheckList) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MachineHealthCheck, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}
func (in *MachineHealthCheckList) DeepCopy() *MachineHealthCheckList {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(MachineHealthCheckList)
	in.DeepCopyInto(out)
	return out
}
func (in *MachineHealthCheckList) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *MachineHealthCheckSpec) DeepCopyInto(out *MachineHealthCheckSpec) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	in.Selector.DeepCopyInto(&out.Selector)
	return
}
func (in *MachineHealthCheckSpec) DeepCopy() *MachineHealthCheckSpec {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(MachineHealthCheckSpec)
	in.DeepCopyInto(out)
	return out
}
func (in *MachineHealthCheckStatus) DeepCopyInto(out *MachineHealthCheckStatus) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *MachineHealthCheckStatus) DeepCopy() *MachineHealthCheckStatus {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(MachineHealthCheckStatus)
	in.DeepCopyInto(out)
	return out
}
