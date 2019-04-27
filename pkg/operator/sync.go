package operator

import (
	"fmt"
	"time"
	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"path/filepath"
	osev1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-version-operator/lib/resourceapply"
	"github.com/openshift/cluster-version-operator/lib/resourceread"
)

const (
	deploymentRolloutPollInterval	= time.Second
	deploymentRolloutTimeout	= 5 * time.Minute
)

func (optr *Operator) syncAll(config OperatorConfig) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := optr.statusProgressing(); err != nil {
		glog.Errorf("Error syncing ClusterOperatorStatus: %v", err)
		return fmt.Errorf("error syncing ClusterOperatorStatus: %v", err)
	}
	if config.Controllers.Provider != clusterAPIControllerNoOp {
		if err := optr.syncClusterAPIController(config); err != nil {
			if err := optr.statusFailing(err.Error()); err != nil {
				glog.Errorf("Error syncing ClusterOperatorStatus: %v", err)
			}
			glog.Errorf("Error syncing cluster api controller: %v", err)
			return err
		}
		glog.V(3).Info("Synced up all components")
	}
	if err := optr.statusAvailable(); err != nil {
		glog.Errorf("Error syncing ClusterOperatorStatus: %v", err)
		return fmt.Errorf("error syncing ClusterOperatorStatus: %v", err)
	}
	return nil
}
func (optr *Operator) syncClusterAPIController(config OperatorConfig) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	featureGate, err := optr.featureGateLister.Get(MachineAPIFeatureGateName)
	var featureSet osev1.FeatureSet
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		glog.V(2).Infof("Failed to find feature gate %q, will use default feature set", MachineAPIFeatureGateName)
		featureSet = osev1.Default
	} else {
		featureSet = featureGate.Spec.FeatureSet
	}
	features, err := generateFeatureMap(featureSet)
	if err != nil {
		return err
	}
	if enabled, ok := features[FeatureGateMachineHealthCheck]; ok && enabled {
		glog.V(2).Infof("Feature %q is enabled", FeatureGateMachineHealthCheck)
		config.Controllers.MachineHealthCheckEnabled = true
	}
	controllerBytes, err := PopulateTemplate(&config, filepath.Join(optr.ownedManifestsDir, "machine-api-controllers.yaml"))
	if err != nil {
		return err
	}
	controller := resourceread.ReadDeploymentV1OrDie(controllerBytes)
	_, updated, err := resourceapply.ApplyDeployment(optr.kubeClient.AppsV1(), controller)
	if err != nil {
		return err
	}
	if updated {
		return optr.waitForDeploymentRollout(controller)
	}
	return nil
}
func (optr *Operator) waitForDeploymentRollout(resource *appsv1.Deployment) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return wait.Poll(deploymentRolloutPollInterval, deploymentRolloutTimeout, func() (bool, error) {
		d, err := optr.deployLister.Deployments(resource.Namespace).Get(resource.Name)
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		if err != nil {
			glog.Errorf("Error getting Deployment %q during rollout: %v", resource.Name, err)
			return false, nil
		}
		if d.DeletionTimestamp != nil {
			return false, fmt.Errorf("deployment %q is being deleted", resource.Name)
		}
		if d.Generation <= d.Status.ObservedGeneration && d.Status.UpdatedReplicas == d.Status.Replicas && d.Status.UnavailableReplicas == 0 {
			return true, nil
		}
		glog.V(4).Infof("Deployment %q is not ready. status: (replicas: %d, updated: %d, ready: %d, unavailable: %d)", d.Name, d.Status.Replicas, d.Status.UpdatedReplicas, d.Status.ReadyReplicas, d.Status.UnavailableReplicas)
		return false, nil
	})
}
