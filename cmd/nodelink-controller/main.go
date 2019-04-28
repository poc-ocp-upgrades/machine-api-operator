package main

import (
	"flag"
	"fmt"
	"reflect"
	"sync"
	"time"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	kubeclientset "k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	corelister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	mapiv1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset"
	mapiclient "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset"
	mapiinformersfactory "github.com/openshift/cluster-api/pkg/client/informers_generated/externalversions"
	mapiinformers "github.com/openshift/cluster-api/pkg/client/informers_generated/externalversions/machine/v1beta1"
	lister "github.com/openshift/cluster-api/pkg/client/listers_generated/machine/v1beta1"
	kubeinformers "k8s.io/client-go/informers"
)

const (
	maxRetries		= 15
	controllerName		= "nodelink"
	machineAnnotationKey	= "machine.openshift.io/machine"
)

func NewController(nodeInformer coreinformers.NodeInformer, machineInformer mapiinformers.MachineInformer, kubeClient kubeclientset.Interface, capiClient mapiclient.Interface) *Controller {
	_logClusterCodePath()
	defer _logClusterCodePath()
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubeClient.CoreV1().RESTClient()).Events("")})
	c := &Controller{capiClient: capiClient, kubeClient: kubeClient, queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "nodelink")}
	nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{AddFunc: c.addNode, UpdateFunc: c.updateNode, DeleteFunc: c.deleteNode})
	machineInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{AddFunc: c.addMachine, UpdateFunc: c.updateMachine, DeleteFunc: c.deleteMachine})
	c.nodeLister = nodeInformer.Lister()
	c.nodesSynced = nodeInformer.Informer().HasSynced
	c.machinesLister = machineInformer.Lister()
	c.machinesSynced = machineInformer.Informer().HasSynced
	c.syncHandler = c.syncNode
	c.enqueueNode = c.enqueue
	c.machineAddress = make(map[string]*mapiv1.Machine)
	return c
}

type Controller struct {
	capiClient		mapiclient.Interface
	kubeClient		kubeclientset.Interface
	syncHandler		func(hKey string) error
	enqueueNode		func(node *corev1.Node)
	nodeLister		corelister.NodeLister
	nodesSynced		cache.InformerSynced
	machinesLister		lister.MachineLister
	machinesSynced		cache.InformerSynced
	queue			workqueue.RateLimitingInterface
	machineAddress		map[string]*mapiv1.Machine
	machineAddressMux	sync.Mutex
}

func (c *Controller) addNode(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	node := obj.(*corev1.Node)
	glog.V(3).Infof("Adding node: %q", node.Name)
	c.enqueueNode(node)
}
func (c *Controller) updateNode(old, cur interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	curNode := cur.(*corev1.Node)
	glog.V(3).Infof("Updating node: %q", curNode.Name)
	c.enqueueNode(curNode)
}
func (c *Controller) deleteNode(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	node, ok := obj.(*corev1.Node)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		node, ok = tombstone.Obj.(*corev1.Node)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Node %#v", obj))
			return
		}
	}
	glog.V(3).Infof("Deleting node")
	c.enqueueNode(node)
}
func (c *Controller) addMachine(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	machine := obj.(*mapiv1.Machine)
	c.machineAddressMux.Lock()
	defer c.machineAddressMux.Unlock()
	for _, a := range machine.Status.Addresses {
		if a.Type == corev1.NodeInternalIP {
			glog.V(3).Infof("Adding machine %q into machineAddress list for %q", machine.Name, a.Address)
			c.machineAddress[a.Address] = machine
			break
		}
	}
}
func (c *Controller) updateMachine(old, cur interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	machine := cur.(*mapiv1.Machine)
	c.machineAddressMux.Lock()
	defer c.machineAddressMux.Unlock()
	for _, a := range machine.Status.Addresses {
		if a.Type == corev1.NodeInternalIP {
			c.machineAddress[a.Address] = machine
			glog.V(3).Infof("Updating machine addresses list. Machine: %q, address: %q", machine.Name, a.Address)
			break
		}
	}
}
func (c *Controller) deleteMachine(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	machine := obj.(*mapiv1.Machine)
	c.machineAddressMux.Lock()
	defer c.machineAddressMux.Unlock()
	for _, a := range machine.Status.Addresses {
		if a.Type == corev1.NodeInternalIP {
			delete(c.machineAddress, a.Address)
			break
		}
	}
	glog.V(3).Infof("Delete obsolete machines from machine addresses list")
}
func WaitForCacheSync(controllerName string, stopCh <-chan struct{}, cacheSyncs ...cache.InformerSynced) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Infof("Waiting for caches to sync for %s controller", controllerName)
	if !cache.WaitForCacheSync(stopCh, cacheSyncs...) {
		utilruntime.HandleError(fmt.Errorf("unable to sync caches for %s controller", controllerName))
		return false
	}
	glog.Infof("Caches are synced for %s controller", controllerName)
	return true
}
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()
	glog.Infof("Starting nodelink controller")
	defer glog.Infof("Shutting down nodelink controller")
	if !WaitForCacheSync("machine", stopCh, c.machinesSynced) {
		return
	}
	if !WaitForCacheSync("node", stopCh, c.nodesSynced) {
		return
	}
	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}
	<-stopCh
}
func (c *Controller) enqueue(node *corev1.Node) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(node)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", node, err))
		return
	}
	c.queue.Add(key)
}
func (c *Controller) worker() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for c.processNextWorkItem() {
	}
}
func (c *Controller) processNextWorkItem() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)
	err := c.syncHandler(key.(string))
	c.handleErr(err, key)
	return true
}
func (c *Controller) handleErr(err error, key interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err == nil {
		c.queue.Forget(key)
		return
	}
	if c.queue.NumRequeues(key) < maxRetries {
		glog.Infof("Error syncing node %v: %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}
	utilruntime.HandleError(err)
	glog.Infof("Dropping node %q out of the queue: %v", key, err)
	c.queue.Forget(key)
}
func (c *Controller) syncNode(key string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	startTime := time.Now()
	glog.V(3).Infof("Syncing node")
	defer func() {
		glog.V(3).Infof("Finished syncing node, duration: %s", time.Now().Sub(startTime))
	}()
	_, _, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	node, err := c.nodeLister.Get(key)
	if errors.IsNotFound(err) {
		glog.Infof("Error syncing, Node %s has been deleted", key)
		return nil
	}
	if err != nil {
		return err
	}
	return c.processNode(node)
}
func (c *Controller) processNode(node *corev1.Node) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	machineKey, ok := node.Annotations[machineAnnotationKey]
	var matchingMachine *mapiv1.Machine
	if ok {
		var err error
		namespace, machineName, err := cache.SplitMetaNamespaceKey(machineKey)
		if err != nil {
			glog.Infof("Error processing node %q. Machine annotation format is incorrect %q: %v", node.Name, machineKey, err)
			return err
		}
		matchingMachine, err = c.machinesLister.Machines(namespace).Get(machineName)
		if err != nil {
			if errors.IsNotFound(err) {
				glog.Warningf("Machine %q associated to node %q has been deleted, will attempt to find new machine by IP", machineKey, node.Name)
			} else {
				return err
			}
		}
	}
	if matchingMachine == nil {
		var nodeInternalIP string
		for _, a := range node.Status.Addresses {
			if a.Type == corev1.NodeInternalIP {
				nodeInternalIP = a.Address
				break
			}
		}
		if nodeInternalIP == "" {
			glog.Warningf("Unable to find InternalIP for node %q", node.Name)
			return fmt.Errorf("unable to find InternalIP for node: %q", node.Name)
		}
		glog.V(4).Infof("Searching machine cache for IP match for node %q", node.Name)
		c.machineAddressMux.Lock()
		machine, found := c.machineAddress[nodeInternalIP]
		c.machineAddressMux.Unlock()
		if found {
			matchingMachine = machine
		}
	}
	if matchingMachine == nil {
		return fmt.Errorf("no machine was found for node: %q", node.Name)
	}
	glog.V(3).Infof("Found machine %s for node %s", machineKey, node.Name)
	modNode := node.DeepCopy()
	if modNode.Annotations == nil {
		modNode.Annotations = map[string]string{}
	}
	modNode.Annotations[machineAnnotationKey] = fmt.Sprintf("%s/%s", matchingMachine.Namespace, matchingMachine.Name)
	if modNode.Labels == nil {
		modNode.Labels = map[string]string{}
	}
	for k, v := range matchingMachine.Spec.Labels {
		glog.V(3).Infof("Copying label %s = %s", k, v)
		modNode.Labels[k] = v
	}
	addTaintsToNode(modNode, matchingMachine)
	if !reflect.DeepEqual(node, modNode) {
		glog.V(3).Infof("Node %q has changed, updating", modNode.Name)
		_, err := c.kubeClient.CoreV1().Nodes().Update(modNode)
		if err != nil {
			glog.Errorf("Error updating node: %v", err)
			return err
		}
	}
	return nil
}
func addTaintsToNode(node *corev1.Node, machine *mapiv1.Machine) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, mTaint := range machine.Spec.Taints {
		glog.V(3).Infof("Adding taint %v from machine %q to node %q", mTaint, machine.Name, node.Name)
		alreadyPresent := false
		for _, nTaint := range node.Spec.Taints {
			if nTaint.Key == mTaint.Key && nTaint.Effect == mTaint.Effect {
				glog.V(3).Infof("Skipping to add machine taint, %v, to the node. Node already has a taint with same key and effect", mTaint)
				alreadyPresent = true
				break
			}
		}
		if !alreadyPresent {
			node.Spec.Taints = append(node.Spec.Taints, mTaint)
		}
	}
}

var (
	logLevel string
)

func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	config, err := config.GetConfig()
	if err != nil {
		glog.Fatalf("Could not create Config for talking to the apiserver: %v", err)
	}
	client, err := clientset.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Could not create client for talking to the apiserver: %v", err)
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Could not create kubernetes client to talk to the apiserver: %v", err)
	}
	kubeSharedInformers := kubeinformers.NewSharedInformerFactory(kubeClient, 5*time.Second)
	mapiInformers := mapiinformersfactory.NewSharedInformerFactory(client, 5*time.Second)
	ctrl := NewController(kubeSharedInformers.Core().V1().Nodes(), mapiInformers.Machine().V1beta1().Machines(), kubeClient, client)
	go ctrl.Run(1, wait.NeverStop)
	mapiInformers.Start(wait.NeverStop)
	kubeSharedInformers.Start(wait.NeverStop)
	select {}
}
