package operator

import (
	"fmt"
	"reflect"
	"strings"
	v1 "k8s.io/api/core/v1"
	"github.com/golang/glog"
	osconfigv1 "github.com/openshift/api/config/v1"
	cvoresourcemerge "github.com/openshift/cluster-version-operator/lib/resourcemerge"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StatusReason string

const (
	ReasonEmpty		StatusReason	= ""
	ReasonSyncing		StatusReason	= "SyncingResources"
	ReasonSyncFailed	StatusReason	= "SyncingFailed"
)
const (
	clusterOperatorName = "machine-api"
)

func (optr *Operator) statusProgressing() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	desiredVersions := optr.operandVersions
	currentVersions, err := optr.getCurrentVersions()
	if err != nil {
		glog.Errorf("Error getting operator current versions: %v", err)
		return err
	}
	var isProgressing osconfigv1.ConditionStatus
	co, err := optr.getOrCreateClusterOperator()
	if err != nil {
		glog.Errorf("Failed to get or create Cluster Operator: %v", err)
		return err
	}
	var message string
	if !reflect.DeepEqual(desiredVersions, currentVersions) {
		glog.V(2).Info("Syncing status: progressing")
		message = fmt.Sprintf("Progressing towards %s", optr.printOperandVersions())
		optr.eventRecorder.Eventf(co, v1.EventTypeNormal, "Status upgrade", message)
		isProgressing = osconfigv1.ConditionTrue
	} else {
		glog.V(2).Info("Syncing status: re-syncing")
		message = fmt.Sprintf("Running resync for %s", optr.printOperandVersions())
		isProgressing = osconfigv1.ConditionFalse
	}
	conds := []osconfigv1.ClusterOperatorStatusCondition{newClusterOperatorStatusCondition(osconfigv1.OperatorProgressing, isProgressing, string(ReasonSyncing), message), newClusterOperatorStatusCondition(osconfigv1.OperatorAvailable, osconfigv1.ConditionTrue, "", ""), newClusterOperatorStatusCondition(osconfigv1.OperatorFailing, osconfigv1.ConditionFalse, "", "")}
	return optr.syncStatus(co, conds)
}
func (optr *Operator) statusAvailable() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	conds := []osconfigv1.ClusterOperatorStatusCondition{newClusterOperatorStatusCondition(osconfigv1.OperatorAvailable, osconfigv1.ConditionTrue, string(ReasonEmpty), fmt.Sprintf("Cluster Machine API Operator is available at %s", optr.printOperandVersions())), newClusterOperatorStatusCondition(osconfigv1.OperatorProgressing, osconfigv1.ConditionFalse, "", ""), newClusterOperatorStatusCondition(osconfigv1.OperatorFailing, osconfigv1.ConditionFalse, "", "")}
	co, err := optr.getOrCreateClusterOperator()
	if err != nil {
		return err
	}
	co.Status.Versions = optr.operandVersions
	glog.V(2).Info("Syncing status: available")
	return optr.syncStatus(co, conds)
}
func (optr *Operator) statusFailing(error string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	desiredVersions := optr.operandVersions
	currentVersions, err := optr.getCurrentVersions()
	if err != nil {
		glog.Errorf("Error getting current versions: %v", err)
		return err
	}
	var message string
	if !reflect.DeepEqual(desiredVersions, currentVersions) {
		message = fmt.Sprintf("Failed when progressing towards %s because %s", optr.printOperandVersions(), error)
	} else {
		message = fmt.Sprintf("Failed to resync for %s because %s", optr.printOperandVersions(), error)
	}
	conds := []osconfigv1.ClusterOperatorStatusCondition{newClusterOperatorStatusCondition(osconfigv1.OperatorFailing, osconfigv1.ConditionTrue, string(ReasonSyncFailed), message), newClusterOperatorStatusCondition(osconfigv1.OperatorProgressing, osconfigv1.ConditionFalse, "", ""), newClusterOperatorStatusCondition(osconfigv1.OperatorAvailable, osconfigv1.ConditionTrue, "", "")}
	co, err := optr.getOrCreateClusterOperator()
	if err != nil {
		return err
	}
	optr.eventRecorder.Eventf(co, v1.EventTypeWarning, "Status failing", error)
	glog.V(2).Info("Syncing status: failing")
	return optr.syncStatus(co, conds)
}
func newClusterOperatorStatusCondition(conditionType osconfigv1.ClusterStatusConditionType, conditionStatus osconfigv1.ConditionStatus, reason string, message string) osconfigv1.ClusterOperatorStatusCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return osconfigv1.ClusterOperatorStatusCondition{Type: conditionType, Status: conditionStatus, LastTransitionTime: metav1.Now(), Reason: reason, Message: message}
}
func (optr *Operator) syncStatus(co *osconfigv1.ClusterOperator, conds []osconfigv1.ClusterOperatorStatusCondition) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, c := range conds {
		cvoresourcemerge.SetOperatorStatusCondition(&co.Status.Conditions, c)
	}
	_, err := optr.osClient.ConfigV1().ClusterOperators().UpdateStatus(co)
	return err
}
func (optr *Operator) getOrCreateClusterOperator() (*osconfigv1.ClusterOperator, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	co, err := optr.osClient.ConfigV1().ClusterOperators().Get(clusterOperatorName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		co = &osconfigv1.ClusterOperator{ObjectMeta: metav1.ObjectMeta{Name: clusterOperatorName}, Status: osconfigv1.ClusterOperatorStatus{Versions: optr.operandVersions, RelatedObjects: []osconfigv1.ObjectReference{{Group: "", Resource: "namespaces", Name: optr.namespace}}}}
		glog.Infof("%s clusterOperator status does not exist, creating %v", clusterOperatorName, co)
		co, err := optr.osClient.ConfigV1().ClusterOperators().Create(co)
		if err != nil {
			return nil, err
		}
		return co, nil
	}
	if err != nil {
		return nil, err
	}
	return co, nil
}
func (optr *Operator) getCurrentVersions() ([]osconfigv1.OperandVersion, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	co, err := optr.getOrCreateClusterOperator()
	if err != nil {
		return nil, err
	}
	return co.Status.Versions, nil
}
func (optr *Operator) printOperandVersions() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	versionsOutput := []string{}
	for _, operand := range optr.operandVersions {
		versionsOutput = append(versionsOutput, fmt.Sprintf("%s: %s", operand.Name, operand.Version))
	}
	return strings.Join(versionsOutput, ", ")
}
