/*
Copyright 2014 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package app implements a server that runs a set of active
// components.  This includes replication controllers, service endpoints and
// nodes.
//
// CAUTION: If you update code in this file, you may need to also update code
//          in contrib/mesos/pkg/controllermanager/controllermanager.go
package app

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	"k8s.io/kubernetes/pkg/client"
	"k8s.io/kubernetes/pkg/client/clientcmd"
	clientcmdapi "k8s.io/kubernetes/pkg/client/clientcmd/api"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller/endpoint"
	"k8s.io/kubernetes/pkg/controller/namespace"
	"k8s.io/kubernetes/pkg/controller/node"
	"k8s.io/kubernetes/pkg/controller/persistentvolume"
	replicationControllerPkg "k8s.io/kubernetes/pkg/controller/replication"
	"k8s.io/kubernetes/pkg/controller/resourcequota"
	"k8s.io/kubernetes/pkg/controller/route"
	"k8s.io/kubernetes/pkg/controller/service"
	"k8s.io/kubernetes/pkg/controller/serviceaccount"
	"k8s.io/kubernetes/pkg/healthz"
	"k8s.io/kubernetes/pkg/master/ports"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/volume"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/pflag"
)

// CMServer is the main context object for the controller manager.
type CMServer struct {
	Port                    int
	Address                 net.IP
	CloudProvider           string
	CloudConfigFile         string
	ConcurrentEndpointSyncs int
	ConcurrentRCSyncs       int
	ServiceSyncPeriod       time.Duration
	NodeSyncPeriod          time.Duration
	ResourceQuotaSyncPeriod time.Duration
	NamespaceSyncPeriod     time.Duration
	PVClaimBinderSyncPeriod time.Duration
	RegisterRetryCount      int
	NodeMonitorGracePeriod  time.Duration
	NodeStartupGracePeriod  time.Duration
	NodeMonitorPeriod       time.Duration
	NodeStatusUpdateRetry   int
	PodEvictionTimeout      time.Duration
	DeletingPodsQps         float32
	DeletingPodsBurst       int
	ServiceAccountKeyFile   string
	RootCAFile              string

	// volumeConfig
	PersistentVolumeRecyclerDefaultScrubPod          string
	PersistentVolumeRecyclerMinTimeoutNfs            int
	PersistentVolumeRecyclerTimeoutIncrementNfs      int
	PersistentVolumeRecyclerMinTimeoutHostPath       int
	PersistentVolumeRecyclerTimeoutIncrementHostPath int

	ClusterName       string
	ClusterCIDR       net.IPNet
	AllocateNodeCIDRs bool
	EnableProfiling   bool

	Master     string
	Kubeconfig string
}

// NewCMServer creates a new CMServer with a default config.
func NewCMServer() *CMServer {
	s := CMServer{
		Port:                             ports.ControllerManagerPort,
		Address:                          net.ParseIP("127.0.0.1"),
		ConcurrentEndpointSyncs:          5,
		ConcurrentRCSyncs:                5,
		ServiceSyncPeriod:                5 * time.Minute,
		NodeSyncPeriod:                   10 * time.Second,
		ResourceQuotaSyncPeriod:          10 * time.Second,
		NamespaceSyncPeriod:              5 * time.Minute,
		PVClaimBinderSyncPeriod:          10 * time.Second,
		PersistentVolumeRecyclerMinTimeoutNfs:            300,
		PersistentVolumeRecyclerTimeoutIncrementNfs:      30,
		PersistentVolumeRecyclerMinTimeoutHostPath:       60,
		PersistentVolumeRecyclerTimeoutIncrementHostPath: 30,
		RegisterRetryCount:               10,
		PodEvictionTimeout:               5 * time.Minute,
		ClusterName:                      "kubernetes",
	}
	return &s
}

// AddFlags adds flags for a specific CMServer to the specified FlagSet
func (s *CMServer) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&s.Port, "port", s.Port, "The port that the controller-manager's http service runs on")
	fs.IPVar(&s.Address, "address", s.Address, "The IP address to serve on (set to 0.0.0.0 for all interfaces)")
	fs.StringVar(&s.CloudProvider, "cloud-provider", s.CloudProvider, "The provider for cloud services.  Empty string for no provider.")
	fs.StringVar(&s.CloudConfigFile, "cloud-config", s.CloudConfigFile, "The path to the cloud provider configuration file.  Empty string for no configuration file.")
	fs.IntVar(&s.ConcurrentEndpointSyncs, "concurrent-endpoint-syncs", s.ConcurrentEndpointSyncs, "The number of endpoint syncing operations that will be done concurrently. Larger number = faster endpoint updating, but more CPU (and network) load")
	fs.IntVar(&s.ConcurrentRCSyncs, "concurrent_rc_syncs", s.ConcurrentRCSyncs, "The number of replication controllers that are allowed to sync concurrently. Larger number = more reponsive replica management, but more CPU (and network) load")
	fs.DurationVar(&s.ServiceSyncPeriod, "service-sync-period", s.ServiceSyncPeriod, "The period for syncing services with their external load balancers")
	fs.DurationVar(&s.NodeSyncPeriod, "node-sync-period", s.NodeSyncPeriod, ""+
		"The period for syncing nodes from cloudprovider. Longer periods will result in "+
		"fewer calls to cloud provider, but may delay addition of new nodes to cluster.")
	fs.DurationVar(&s.ResourceQuotaSyncPeriod, "resource-quota-sync-period", s.ResourceQuotaSyncPeriod, "The period for syncing quota usage status in the system")
	fs.DurationVar(&s.NamespaceSyncPeriod, "namespace-sync-period", s.NamespaceSyncPeriod, "The period for syncing namespace life-cycle updates")
	fs.DurationVar(&s.PVClaimBinderSyncPeriod, "pvclaimbinder-sync-period", s.PVClaimBinderSyncPeriod, "The period for syncing persistent volumes and persistent volume claims")
	fs.StringVar(&s.PersistentVolumeRecyclerDefaultScrubPod, "pv-recycler-default-scrub-pod", s.PersistentVolumeRecyclerDefaultScrubPod, "The file path to a pod definition used as a template for persistent volume recycling")
	fs.IntVar(&s.PersistentVolumeRecyclerMinTimeoutNfs, "pv-recycler-min-timeout-nfs", s.PersistentVolumeRecyclerMinTimeoutNfs, "The minimum ActiveDeadlineSeconds to use for an NFS Recycler pod")
	fs.IntVar(&s.PersistentVolumeRecyclerTimeoutIncrementNfs, "pv-recycler-timeout-increment-nfs", s.PersistentVolumeRecyclerTimeoutIncrementNfs, "the increment of time added per Gi to ActiveDeadlineSeconds for an NFS scrubber pod")
	fs.IntVar(&s.PersistentVolumeRecyclerMinTimeoutHostPath, "pv-recycler-min-timeout-hostpath", s.PersistentVolumeRecyclerMinTimeoutHostPath, "The minimum ActiveDeadlineSeconds to use for a HostPath Recycler pod")
	fs.IntVar(&s.PersistentVolumeRecyclerTimeoutIncrementHostPath, "pv-recycler-timeout-increment-hostpath", s.PersistentVolumeRecyclerTimeoutIncrementHostPath, "the increment of time added per Gi to ActiveDeadlineSeconds for a HostPath scrubber pod")
	fs.DurationVar(&s.PodEvictionTimeout, "pod-eviction-timeout", s.PodEvictionTimeout, "The grace period for deleting pods on failed nodes.")
	fs.Float32Var(&s.DeletingPodsQps, "deleting-pods-qps", 0.1, "Number of nodes per second on which pods are deleted in case of node failure.")
	fs.IntVar(&s.DeletingPodsBurst, "deleting-pods-burst", 10, "Number of nodes on which pods are bursty deleted in case of node failure. For more details look into RateLimiter.")
	fs.IntVar(&s.RegisterRetryCount, "register-retry-count", s.RegisterRetryCount, ""+
		"The number of retries for initial node registration.  Retry interval equals node-sync-period.")
	fs.MarkDeprecated("register-retry-count", "This flag is currently no-op and will be deleted.")
	fs.DurationVar(&s.NodeMonitorGracePeriod, "node-monitor-grace-period", 40*time.Second,
		"Amount of time which we allow running Node to be unresponsive before marking it unhealty. "+
			"Must be N times more than kubelet's nodeStatusUpdateFrequency, "+
			"where N means number of retries allowed for kubelet to post node status.")
	fs.DurationVar(&s.NodeStartupGracePeriod, "node-startup-grace-period", 60*time.Second,
		"Amount of time which we allow starting Node to be unresponsive before marking it unhealty.")
	fs.DurationVar(&s.NodeMonitorPeriod, "node-monitor-period", 5*time.Second,
		"The period for syncing NodeStatus in NodeController.")
	fs.StringVar(&s.ServiceAccountKeyFile, "service-account-private-key-file", s.ServiceAccountKeyFile, "Filename containing a PEM-encoded private RSA key used to sign service account tokens.")
	fs.BoolVar(&s.EnableProfiling, "profiling", true, "Enable profiling via web interface host:port/debug/pprof/")
	fs.StringVar(&s.ClusterName, "cluster-name", s.ClusterName, "The instance prefix for the cluster")
	fs.IPNetVar(&s.ClusterCIDR, "cluster-cidr", s.ClusterCIDR, "CIDR Range for Pods in cluster.")
	fs.BoolVar(&s.AllocateNodeCIDRs, "allocate-node-cidrs", false, "Should CIDRs for Pods be allocated and set on the cloud provider.")
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	fs.StringVar(&s.RootCAFile, "root-ca-file", s.RootCAFile, "If set, this root certificate authority will be included in service account's token secret. This must be a valid PEM-encoded CA bundle.")
}

// Run runs the CMServer.  This should never exit.
func (s *CMServer) Run(_ []string) error {
	if s.Kubeconfig == "" && s.Master == "" {
		glog.Warningf("Neither --kubeconfig nor --master was specified.  Using default API client.  This might not work.")
	}

	// This creates a client, first loading any specified kubeconfig
	// file, and then overriding the Master flag, if non-empty.
	kubeconfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: s.Kubeconfig},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: s.Master}}).ClientConfig()
	if err != nil {
		return err
	}

	kubeconfig.QPS = 20.0
	kubeconfig.Burst = 30

	kubeClient, err := client.New(kubeconfig)
	if err != nil {
		glog.Fatalf("Invalid API configuration: %v", err)
	}

	go func() {
		mux := http.NewServeMux()
		healthz.InstallHandler(mux)
		if s.EnableProfiling {
			mux.HandleFunc("/debug/pprof/", pprof.Index)
			mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		}
		mux.Handle("/metrics", prometheus.Handler())

		server := &http.Server{
			Addr:    net.JoinHostPort(s.Address.String(), strconv.Itoa(s.Port)),
			Handler: mux,
		}
		glog.Fatal(server.ListenAndServe())
	}()

	endpoints := endpointcontroller.NewEndpointController(kubeClient)
	go endpoints.Run(s.ConcurrentEndpointSyncs, util.NeverStop)

	controllerManager := replicationControllerPkg.NewReplicationManager(kubeClient, replicationControllerPkg.BurstReplicas)
	go controllerManager.Run(s.ConcurrentRCSyncs, util.NeverStop)

	cloud, err := cloudprovider.InitCloudProvider(s.CloudProvider, s.CloudConfigFile)
	if err != nil {
		glog.Fatalf("Cloud provider could not be initialized: %v", err)
	}

	nodeController := nodecontroller.NewNodeController(cloud, kubeClient,
		s.PodEvictionTimeout, nodecontroller.NewPodEvictor(util.NewTokenBucketRateLimiter(s.DeletingPodsQps, s.DeletingPodsBurst)),
		s.NodeMonitorGracePeriod, s.NodeStartupGracePeriod, s.NodeMonitorPeriod, &s.ClusterCIDR, s.AllocateNodeCIDRs)
	nodeController.Run(s.NodeSyncPeriod)

	serviceController := servicecontroller.New(cloud, kubeClient, s.ClusterName)
	if err := serviceController.Run(s.ServiceSyncPeriod, s.NodeSyncPeriod); err != nil {
		glog.Errorf("Failed to start service controller: %v", err)
	}

	if s.AllocateNodeCIDRs {
		if cloud == nil {
			glog.Warning("allocate-node-cidrs is set, but no cloud provider specified. Will not manage routes.")
		} else if routes, ok := cloud.Routes(); !ok {
			glog.Warning("allocate-node-cidrs is set, but cloud provider does not support routes. Will not manage routes.")
		} else {
			routeController := routecontroller.New(routes, kubeClient, s.ClusterName, &s.ClusterCIDR)
			routeController.Run(s.NodeSyncPeriod)
		}
	}

	resourceQuotaController := resourcequotacontroller.NewResourceQuotaController(kubeClient)
	resourceQuotaController.Run(s.ResourceQuotaSyncPeriod)

	namespaceController := namespacecontroller.NewNamespaceController(kubeClient, s.NamespaceSyncPeriod)
	namespaceController.Run()

	pvclaimBinder := volumeclaimbinder.NewPersistentVolumeClaimBinder(kubeClient, s.PVClaimBinderSyncPeriod)
	pvclaimBinder.Run()

	volumeConfig := volume.NewVolumeConfig()
	volumeConfig.PersistentVolumeRecyclerMinTimeoutHostPath = int64(s.PersistentVolumeRecyclerMinTimeoutHostPath)
	volumeConfig.PersistentVolumeRecyclerTimeoutIncrementHostPath = int64(s.PersistentVolumeRecyclerTimeoutIncrementHostPath)
	volumeConfig.PersistentVolumeRecyclerMinTimeoutNfs = int64(s.PersistentVolumeRecyclerMinTimeoutNfs)
	volumeConfig.PersistentVolumeRecyclerTimeoutIncrementNfs = int64(s.PersistentVolumeRecyclerTimeoutIncrementNfs)
	if s.PersistentVolumeRecyclerDefaultScrubPod != "" {
		scrubPod, err := volume.InitScrubPod(s.PersistentVolumeRecyclerDefaultScrubPod)
		if err != nil {
			glog.Fatalf("Override of default PersistentVolume scrub pod failed: %+v", err)
		}
		volumeConfig.PersistentVolumeRecyclerDefaultScrubPod = scrubPod
	}

	pvRecycler, err := volumeclaimbinder.NewPersistentVolumeRecycler(kubeClient, s.PVClaimBinderSyncPeriod, ProbeRecyclableVolumePlugins(volumeConfig))
	if err != nil {
		glog.Fatalf("Failed to start persistent volume recycler: %+v", err)
	}
	pvRecycler.Run()

	var rootCA []byte

	if s.RootCAFile != "" {
		rootCA, err = ioutil.ReadFile(s.RootCAFile)
		if err != nil {
			return fmt.Errorf("error reading root-ca-file at %s: %v", s.RootCAFile, err)
		}
		if _, err := util.CertsFromPEM(rootCA); err != nil {
			return fmt.Errorf("error parsing root-ca-file at %s: %v", s.RootCAFile, err)
		}
	} else {
		rootCA = kubeconfig.CAData
	}

	if len(s.ServiceAccountKeyFile) > 0 {
		privateKey, err := serviceaccount.ReadPrivateKey(s.ServiceAccountKeyFile)
		if err != nil {
			glog.Errorf("Error reading key for service account token controller: %v", err)
		} else {
			serviceaccount.NewTokensController(
				kubeClient,
				serviceaccount.TokensControllerOptions{
					TokenGenerator: serviceaccount.JWTTokenGenerator(privateKey),
					RootCA:         rootCA,
				},
			).Run()
		}
	}

	serviceaccount.NewServiceAccountsController(
		kubeClient,
		serviceaccount.DefaultServiceAccountsControllerOptions(),
	).Run()

	select {}
	return nil
}
