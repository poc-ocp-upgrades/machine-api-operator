package operator

import (
	"fmt"
	"os"
	"time"
	"github.com/golang/glog"
	osconfigv1 "github.com/openshift/api/config/v1"
	osclientset "github.com/openshift/client-go/config/clientset/versioned"
	configinformersv1 "github.com/openshift/client-go/config/informers/externalversions/config/v1"
	configlistersv1 "github.com/openshift/client-go/config/listers/config/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformersv1 "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	coreclientsetv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisterv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const (
	maxRetries		= 15
	ownedManifestsDir	= "owned-manifests"
)

type Operator struct {
	namespace, name		string
	imagesFile		string
	config			string
	ownedManifestsDir	string
	kubeClient		kubernetes.Interface
	osClient		osclientset.Interface
	eventRecorder		record.EventRecorder
	syncHandler		func(ic string) error
	deployLister		appslisterv1.DeploymentLister
	deployListerSynced	cache.InformerSynced
	featureGateLister	configlistersv1.FeatureGateLister
	featureGateCacheSynced	cache.InformerSynced
	queue			workqueue.RateLimitingInterface
	operandVersions		[]osconfigv1.OperandVersion
}

func New(namespace, name string, imagesFile string, config string, deployInformer appsinformersv1.DeploymentInformer, featureGateInformer configinformersv1.FeatureGateInformer, kubeClient kubernetes.Interface, osClient osclientset.Interface) *Operator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&coreclientsetv1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	operandVersions := []osconfigv1.OperandVersion{}
	if releaseVersion := os.Getenv("RELEASE_VERSION"); len(releaseVersion) > 0 {
		operandVersions = append(operandVersions, osconfigv1.OperandVersion{Name: "operator", Version: releaseVersion})
	}
	eventRecorderScheme := runtime.NewScheme()
	osconfigv1.Install(eventRecorderScheme)
	optr := &Operator{namespace: namespace, name: name, imagesFile: imagesFile, ownedManifestsDir: ownedManifestsDir, kubeClient: kubeClient, osClient: osClient, eventRecorder: eventBroadcaster.NewRecorder(eventRecorderScheme, v1.EventSource{Component: "machineapioperator"}), queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "machineapioperator"), operandVersions: operandVersions}
	deployInformer.Informer().AddEventHandler(optr.eventHandler())
	featureGateInformer.Informer().AddEventHandler(optr.eventHandler())
	optr.config = config
	optr.syncHandler = optr.sync
	optr.deployLister = deployInformer.Lister()
	optr.deployListerSynced = deployInformer.Informer().HasSynced
	optr.featureGateLister = featureGateInformer.Lister()
	optr.featureGateCacheSynced = featureGateInformer.Informer().HasSynced
	return optr
}
func (optr *Operator) Run(workers int, stopCh <-chan struct{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer utilruntime.HandleCrash()
	defer optr.queue.ShutDown()
	glog.Info("Starting Machine API Operator")
	defer glog.Info("Shutting down Machine API Operator")
	if !cache.WaitForCacheSync(stopCh, optr.deployListerSynced, optr.featureGateCacheSynced) {
		glog.Error("Failed to sync caches")
		return
	}
	glog.Info("Synced up caches")
	for i := 0; i < workers; i++ {
		go wait.Until(optr.worker, time.Second, stopCh)
	}
	<-stopCh
}
func (optr *Operator) eventHandler() cache.ResourceEventHandler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	workQueueKey := fmt.Sprintf("%s/%s", optr.namespace, optr.name)
	return cache.ResourceEventHandlerFuncs{AddFunc: func(obj interface{}) {
		optr.queue.Add(workQueueKey)
	}, UpdateFunc: func(old, new interface{}) {
		optr.queue.Add(workQueueKey)
	}, DeleteFunc: func(obj interface{}) {
		optr.queue.Add(workQueueKey)
	}}
}
func (optr *Operator) worker() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for optr.processNextWorkItem() {
	}
}
func (optr *Operator) processNextWorkItem() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	key, quit := optr.queue.Get()
	if quit {
		return false
	}
	defer optr.queue.Done(key)
	glog.V(4).Infof("Processing key %s", key)
	err := optr.syncHandler(key.(string))
	optr.handleErr(err, key)
	return true
}
func (optr *Operator) handleErr(err error, key interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err == nil {
		optr.queue.Forget(key)
		return
	}
	if optr.queue.NumRequeues(key) < maxRetries {
		glog.V(1).Infof("Error syncing operator %v: %v", key, err)
		optr.queue.AddRateLimited(key)
		return
	}
	utilruntime.HandleError(err)
	glog.V(1).Infof("Dropping operator %q out of the queue: %v", key, err)
	optr.queue.Forget(key)
}
func (optr *Operator) sync(key string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	startTime := time.Now()
	glog.V(4).Infof("Started syncing operator %q (%v)", key, startTime)
	defer func() {
		glog.V(4).Infof("Finished syncing operator %q (%v)", key, time.Since(startTime))
	}()
	operatorConfig, err := optr.maoConfigFromInfrastructure()
	if err != nil {
		glog.Errorf("Failed getting operator config: %v", err)
		return err
	}
	return optr.syncAll(*operatorConfig)
}
func (optr *Operator) maoConfigFromInfrastructure() (*OperatorConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	infra, err := optr.osClient.ConfigV1().Infrastructures().Get("cluster", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	provider, err := getProviderFromInfrastructure(infra)
	if err != nil {
		return nil, err
	}
	images, err := getImagesFromJSONFile(optr.imagesFile)
	if err != nil {
		return nil, err
	}
	providerControllerImage, err := getProviderControllerFromImages(provider, *images)
	if err != nil {
		return nil, err
	}
	machineAPIOperatorImage, err := getMachineAPIOperatorFromImages(*images)
	if err != nil {
		return nil, err
	}
	return &OperatorConfig{TargetNamespace: optr.namespace, Controllers: Controllers{Provider: providerControllerImage, NodeLink: machineAPIOperatorImage, MachineHealthCheck: machineAPIOperatorImage}}, nil
}
