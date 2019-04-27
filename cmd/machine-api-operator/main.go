package main

import (
	"flag"
	"os"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var componentNamespace = "openshift-machine-api"

const (
	componentName = "machine-api-operator"
)

var (
	rootCmd	= &cobra.Command{Use: componentName, Short: "Run Cluster API Controller", Long: ""}
	config	string
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}
func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if namespace, ok := os.LookupEnv("COMPONENT_NAMESPACE"); ok {
		componentNamespace = namespace
	}
	if err := rootCmd.Execute(); err != nil {
		glog.Exitf("Error executing mao: %v", err)
	}
}
