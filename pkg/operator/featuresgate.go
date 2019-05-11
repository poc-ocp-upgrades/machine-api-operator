package operator

import (
	"fmt"
	osev1 "github.com/openshift/api/config/v1"
)

const (
	MachineAPIFeatureGateName		= "machine-api"
	FeatureGateMachineHealthCheck	= "MachineHealthCheck"
)

var MachineAPIOperatorFeatureSets = map[osev1.FeatureSet]*osev1.FeatureGateEnabledDisabled{osev1.Default: {Disabled: []string{FeatureGateMachineHealthCheck}}, osev1.TechPreviewNoUpgrade: {Enabled: []string{FeatureGateMachineHealthCheck}}}

func generateFeatureMap(featureSet osev1.FeatureSet) (map[string]bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rv := map[string]bool{}
	set, ok := MachineAPIOperatorFeatureSets[featureSet]
	if !ok {
		return nil, fmt.Errorf("enabled FeatureSet %v does not have a corresponding config", featureSet)
	}
	for _, featEnabled := range set.Enabled {
		rv[featEnabled] = true
	}
	for _, featDisabled := range set.Disabled {
		rv[featDisabled] = false
	}
	return rv, nil
}
