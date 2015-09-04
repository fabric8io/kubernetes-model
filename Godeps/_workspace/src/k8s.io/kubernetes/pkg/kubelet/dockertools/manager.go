/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package dockertools

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
	"github.com/golang/groupcache/lru"
	cadvisorApi "github.com/google/cadvisor/info/v1"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/latest"
	"k8s.io/kubernetes/pkg/client/record"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
	"k8s.io/kubernetes/pkg/kubelet/network"
	"k8s.io/kubernetes/pkg/kubelet/prober"
	"k8s.io/kubernetes/pkg/kubelet/qos"
	kubeletTypes "k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/kubernetes/pkg/probe"
	"k8s.io/kubernetes/pkg/securitycontext"
	"k8s.io/kubernetes/pkg/types"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/util/oom"
	"k8s.io/kubernetes/pkg/util/procfs"
)

const (
	maxReasonCacheEntries = 200

	kubernetesPodLabel       = "io.kubernetes.pod.data"
	kubernetesContainerLabel = "io.kubernetes.container.name"
	// ndots specifies the minimum number of dots that a domain name must contain for the resolver to consider it as FQDN (fully-qualified)
	// we want to able to consider SRV lookup names like _dns._udp.kube-dns.default.svc to be considered relative.
	// hence, setting ndots to be 5.
	ndotsDNSOption = "options ndots:5\n"
)

// DockerManager implements the Runtime interface.
var _ kubecontainer.Runtime = &DockerManager{}

// TODO: make this a TTL based pull (if image older than X policy, pull)
var podInfraContainerImagePullPolicy = api.PullIfNotPresent

type DockerManager struct {
	client              DockerInterface
	recorder            record.EventRecorder
	readinessManager    *kubecontainer.ReadinessManager
	containerRefManager *kubecontainer.RefManager
	os                  kubecontainer.OSInterface
	machineInfo         *cadvisorApi.MachineInfo

	// The image name of the pod infra container.
	podInfraContainerImage string
	// reasonCache stores the failure reason of the last container creation
	// and/or start in a string, keyed by <pod_UID>_<container_name>. The goal
	// is to propagate this reason to the container status. This endeavor is
	// "best-effort" for two reasons:
	//   1. The cache is not persisted.
	//   2. We use an LRU cache to avoid extra garbage collection work. This
	//      means that some entries may be recycled before a pod has been
	//      deleted.
	reasonCache stringCache
	// TODO(yifan): Record the pull failure so we can  eliminate the image checking
	// in GetPodStatus()?
	puller DockerPuller

	// Root of the Docker runtime.
	dockerRoot string

	// Directory of container logs.
	containerLogsDir string

	// Network plugin.
	networkPlugin network.NetworkPlugin

	// Health check prober.
	prober prober.Prober

	// Generator of runtime container options.
	generator kubecontainer.RunContainerOptionsGenerator

	// Runner of lifecycle events.
	runner kubecontainer.HandlerRunner

	// Hooks injected into the container runtime.
	runtimeHooks kubecontainer.RuntimeHooks

	// Handler used to execute commands in containers.
	execHandler ExecHandler

	// Used to set OOM scores of processes.
	oomAdjuster *oom.OomAdjuster

	// Get information from /proc mount.
	procFs procfs.ProcFsInterface
}

func NewDockerManager(
	client DockerInterface,
	recorder record.EventRecorder,
	readinessManager *kubecontainer.ReadinessManager,
	containerRefManager *kubecontainer.RefManager,
	machineInfo *cadvisorApi.MachineInfo,
	podInfraContainerImage string,
	qps float32,
	burst int,
	containerLogsDir string,
	osInterface kubecontainer.OSInterface,
	networkPlugin network.NetworkPlugin,
	generator kubecontainer.RunContainerOptionsGenerator,
	httpClient kubeletTypes.HttpGetter,
	runtimeHooks kubecontainer.RuntimeHooks,
	execHandler ExecHandler,
	oomAdjuster *oom.OomAdjuster,
	procFs procfs.ProcFsInterface) *DockerManager {
	// Work out the location of the Docker runtime, defaulting to /var/lib/docker
	// if there are any problems.
	dockerRoot := "/var/lib/docker"
	dockerInfo, err := client.Info()
	if err != nil {
		glog.Errorf("Failed to execute Info() call to the Docker client: %v", err)
		glog.Warningf("Using fallback default of /var/lib/docker for location of Docker runtime")
	} else {
		driverStatus := dockerInfo.Get("DriverStatus")
		// The DriverStatus is a*string* which represents a list of list of strings (pairs) e.g.
		// DriverStatus=[["Root Dir","/var/lib/docker/aufs"],["Backing Filesystem","extfs"],["Dirs","279"]]
		// Strip out the square brakcets and quotes.
		s := strings.Replace(driverStatus, "[", "", -1)
		s = strings.Replace(s, "]", "", -1)
		s = strings.Replace(s, `"`, "", -1)
		// Separate by commas.
		ss := strings.Split(s, ",")
		// Search for the Root Dir string
		for i, k := range ss {
			if k == "Root Dir" && i+1 < len(ss) {
				// Discard the /aufs suffix.
				dockerRoot, _ = path.Split(ss[i+1])
				// Trim the last slash.
				dockerRoot = strings.TrimSuffix(dockerRoot, "/")
				glog.Infof("Setting dockerRoot to %s", dockerRoot)
			}

		}
	}

	reasonCache := stringCache{cache: lru.New(maxReasonCacheEntries)}

	dm := &DockerManager{
		client:                 client,
		recorder:               recorder,
		readinessManager:       readinessManager,
		containerRefManager:    containerRefManager,
		os:                     osInterface,
		machineInfo:            machineInfo,
		podInfraContainerImage: podInfraContainerImage,
		reasonCache:            reasonCache,
		puller:                 newDockerPuller(client, qps, burst),
		dockerRoot:             dockerRoot,
		containerLogsDir:       containerLogsDir,
		networkPlugin:          networkPlugin,
		prober:                 nil,
		generator:              generator,
		runtimeHooks:           runtimeHooks,
		execHandler:            execHandler,
		oomAdjuster:            oomAdjuster,
		procFs:                 procFs,
	}
	dm.runner = lifecycle.NewHandlerRunner(httpClient, dm, dm)
	dm.prober = prober.New(dm, readinessManager, containerRefManager, recorder)

	return dm
}

// A cache which stores strings keyed by <pod_UID>_<container_name>.
type stringCache struct {
	lock  sync.RWMutex
	cache *lru.Cache
}

func (sc *stringCache) composeKey(uid types.UID, name string) string {
	return fmt.Sprintf("%s_%s", uid, name)
}

func (sc *stringCache) Add(uid types.UID, name string, value string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	sc.cache.Add(sc.composeKey(uid, name), value)
}

func (sc *stringCache) Remove(uid types.UID, name string) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	sc.cache.Remove(sc.composeKey(uid, name))
}

func (sc *stringCache) Get(uid types.UID, name string) (string, bool) {
	sc.lock.RLock()
	defer sc.lock.RUnlock()
	value, ok := sc.cache.Get(sc.composeKey(uid, name))
	if ok {
		return value.(string), ok
	} else {
		return "", ok
	}
}

// GetContainerLogs returns logs of a specific container. By
// default, it returns a snapshot of the container log. Set 'follow' to true to
// stream the log. Set 'follow' to false and specify the number of lines (e.g.
// "100" or "all") to tail the log.
// TODO: Make 'RawTerminal' option  flagable.
func (dm *DockerManager) GetContainerLogs(pod *api.Pod, containerID, tail string, follow bool, stdout, stderr io.Writer) (err error) {
	opts := docker.LogsOptions{
		Container:    containerID,
		Stdout:       true,
		Stderr:       true,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Timestamps:   false,
		RawTerminal:  false,
		Follow:       follow,
	}

	if !follow {
		opts.Tail = tail
	}

	err = dm.client.Logs(opts)
	return
}

var (
	// ErrNoContainersInPod is returned when there are no containers for a given pod
	ErrNoContainersInPod = errors.New("no containers exist for this pod")

	// ErrNoPodInfraContainerInPod is returned when there is no pod infra container for a given pod
	ErrNoPodInfraContainerInPod = errors.New("No pod infra container exists for this pod")

	// ErrContainerCannotRun is returned when a container is created, but cannot run properly
	ErrContainerCannotRun = errors.New("Container cannot run")
)

// Internal information kept for containers from inspection
type containerStatusResult struct {
	status api.ContainerStatus
	ip     string
	err    error
}

func (dm *DockerManager) inspectContainer(dockerID, containerName, tPath string, pod *api.Pod) *containerStatusResult {
	result := containerStatusResult{api.ContainerStatus{}, "", nil}

	inspectResult, err := dm.client.InspectContainer(dockerID)

	if err != nil {
		result.err = err
		return &result
	}
	if inspectResult == nil {
		// Why did we not get an error?
		return &result
	}

	glog.V(3).Infof("Container inspect result: %+v", *inspectResult)
	result.status = api.ContainerStatus{
		Name:        containerName,
		Image:       inspectResult.Config.Image,
		ImageID:     DockerPrefix + inspectResult.Image,
		ContainerID: DockerPrefix + dockerID,
	}

	if inspectResult.State.Running {
		result.status.State.Running = &api.ContainerStateRunning{
			StartedAt: util.NewTime(inspectResult.State.StartedAt),
		}
		if containerName == PodInfraContainerName {
			if inspectResult.NetworkSettings != nil {
				result.ip = inspectResult.NetworkSettings.IPAddress
			}
			// override the above if a network plugin exists
			if dm.networkPlugin.Name() != network.DefaultPluginName {
				netStatus, err := dm.networkPlugin.Status(pod.Namespace, pod.Name, kubeletTypes.DockerID(dockerID))
				if err != nil {
					glog.Errorf("NetworkPlugin %s failed on the status hook for pod '%s' - %v", dm.networkPlugin.Name(), pod.Name, err)
				} else if netStatus != nil {
					result.ip = netStatus.IP.String()
				}
			}
		}
	} else if !inspectResult.State.FinishedAt.IsZero() {
		reason := ""
		// Note: An application might handle OOMKilled gracefully.
		// In that case, the container is oom killed, but the exit
		// code could be 0.
		if inspectResult.State.OOMKilled {
			reason = "OOM Killed"
		} else {
			reason = inspectResult.State.Error
		}
		result.status.State.Terminated = &api.ContainerStateTerminated{
			ExitCode:    inspectResult.State.ExitCode,
			Reason:      reason,
			StartedAt:   util.NewTime(inspectResult.State.StartedAt),
			FinishedAt:  util.NewTime(inspectResult.State.FinishedAt),
			ContainerID: DockerPrefix + dockerID,
		}
		if tPath != "" {
			path, found := inspectResult.Volumes[tPath]
			if found {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					result.status.State.Terminated.Message = fmt.Sprintf("Error on reading termination-log %s: %v", path, err)
				} else {
					result.status.State.Terminated.Message = string(data)
				}
			}
		}
	} else {
		// TODO(dchen1107): Separate issue docker/docker#8294 was filed
		result.status.State.Waiting = &api.ContainerStateWaiting{
			Reason: ErrContainerCannotRun.Error(),
		}
	}

	return &result
}

// GetPodStatus returns docker related status for all containers in the pod as
// well as the infrastructure container.
func (dm *DockerManager) GetPodStatus(pod *api.Pod) (*api.PodStatus, error) {
	podFullName := kubecontainer.GetPodFullName(pod)
	uid := pod.UID
	manifest := pod.Spec

	oldStatuses := make(map[string]api.ContainerStatus, len(pod.Spec.Containers))
	lastObservedTime := make(map[string]util.Time, len(pod.Spec.Containers))
	for _, status := range pod.Status.ContainerStatuses {
		oldStatuses[status.Name] = status
		if status.LastTerminationState.Terminated != nil {
			lastObservedTime[status.Name] = status.LastTerminationState.Terminated.FinishedAt
		}
	}

	var podStatus api.PodStatus
	statuses := make(map[string]*api.ContainerStatus, len(pod.Spec.Containers))

	expectedContainers := make(map[string]api.Container)
	for _, container := range manifest.Containers {
		expectedContainers[container.Name] = container
	}
	expectedContainers[PodInfraContainerName] = api.Container{}

	containers, err := dm.client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return nil, err
	}

	containerDone := util.NewStringSet()
	// Loop through list of running and exited docker containers to construct
	// the statuses. We assume docker returns a list of containers sorted in
	// reverse by time.
	for _, value := range containers {
		if len(value.Names) == 0 {
			continue
		}
		dockerName, _, err := ParseDockerName(value.Names[0])
		if err != nil {
			continue
		}
		if dockerName.PodFullName != podFullName {
			continue
		}
		if uid != "" && dockerName.PodUID != uid {
			continue
		}
		dockerContainerName := dockerName.ContainerName
		c, found := expectedContainers[dockerContainerName]
		if !found {
			continue
		}
		terminationMessagePath := c.TerminationMessagePath
		if containerDone.Has(dockerContainerName) {
			continue
		}

		var terminationState *api.ContainerState = nil
		// Inspect the container.
		result := dm.inspectContainer(value.ID, dockerContainerName, terminationMessagePath, pod)
		if result.err != nil {
			return nil, result.err
		} else if result.status.State.Terminated != nil {
			terminationState = &result.status.State
		}

		if containerStatus, found := statuses[dockerContainerName]; found {
			if containerStatus.LastTerminationState.Terminated == nil && terminationState != nil {
				// Populate the last termination state.
				containerStatus.LastTerminationState = *terminationState
			}
			count := true
			// Only count dead containers terminated after last time we observed,
			if lastObservedTime, ok := lastObservedTime[dockerContainerName]; ok {
				if terminationState != nil && terminationState.Terminated.FinishedAt.After(lastObservedTime.Time) {
					count = false
				} else {
					// The container finished before the last observation. No
					// need to examine/count the older containers. Mark the
					// container name as done.
					containerDone.Insert(dockerContainerName)
				}
			}
			if count {
				containerStatus.RestartCount += 1
			}
			continue
		}

		if dockerContainerName == PodInfraContainerName {
			// Found network container
			if result.status.State.Running != nil {
				podStatus.PodIP = result.ip
			}
		} else {
			// Add user container information.
			if oldStatus, found := oldStatuses[dockerContainerName]; found {
				// Use the last observed restart count if it's available.
				result.status.RestartCount = oldStatus.RestartCount
			}
			statuses[dockerContainerName] = &result.status
		}
	}

	// Handle the containers for which we cannot find any associated active or
	// dead docker containers.
	for _, container := range manifest.Containers {
		if _, found := statuses[container.Name]; found {
			continue
		}
		var containerStatus api.ContainerStatus
		containerStatus.Name = container.Name
		containerStatus.Image = container.Image
		if oldStatus, found := oldStatuses[container.Name]; found {
			// Some states may be lost due to GC; apply the last observed
			// values if possible.
			containerStatus.RestartCount = oldStatus.RestartCount
			containerStatus.LastTerminationState = oldStatus.LastTerminationState
		}
		//Check image is ready on the node or not.
		image := container.Image
		// TODO(dchen1107): docker/docker/issues/8365 to figure out if the image exists
		_, err := dm.client.InspectImage(image)
		if err == nil {
			containerStatus.State.Waiting = &api.ContainerStateWaiting{
				Reason: fmt.Sprintf("Image: %s is ready, container is creating", image),
			}
		} else if err == docker.ErrNoSuchImage {
			containerStatus.State.Waiting = &api.ContainerStateWaiting{
				Reason: fmt.Sprintf("Image: %s is not ready on the node", image),
			}
		}
		statuses[container.Name] = &containerStatus
	}

	podStatus.ContainerStatuses = make([]api.ContainerStatus, 0)
	for containerName, status := range statuses {
		if status.State.Waiting != nil {
			// For containers in the waiting state, fill in a specific reason if it is recorded.
			if reason, ok := dm.reasonCache.Get(uid, containerName); ok {
				status.State.Waiting.Reason = reason
			}
		}
		podStatus.ContainerStatuses = append(podStatus.ContainerStatuses, *status)
	}
	// Sort the container statuses since clients of this interface expect the list
	// of containers in a pod to behave like the output of `docker list`, which has a
	// deterministic order.
	sort.Sort(kubeletTypes.SortedContainerStatuses(podStatus.ContainerStatuses))
	return &podStatus, nil
}

// makeEnvList converts EnvVar list to a list of strings, in the form of
// '<key>=<value>', which can be understood by docker.
func makeEnvList(envs []kubecontainer.EnvVar) (result []string) {
	for _, env := range envs {
		result = append(result, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}
	return
}

// makeMountBindings converts the mount list to a list of strings that
// can be understood by docker.
// Each element in the string is in the form of:
// '<HostPath>:<ContainerPath>', or
// '<HostPath>:<ContainerPath>:ro', if the path is read only.
func makeMountBindings(mounts []kubecontainer.Mount) (result []string) {
	for _, m := range mounts {
		bind := fmt.Sprintf("%s:%s", m.HostPath, m.ContainerPath)
		if m.ReadOnly {
			bind += ":ro"
		}
		result = append(result, bind)
	}
	return
}

func makePortsAndBindings(portMappings []kubecontainer.PortMapping) (map[docker.Port]struct{}, map[docker.Port][]docker.PortBinding) {
	exposedPorts := map[docker.Port]struct{}{}
	portBindings := map[docker.Port][]docker.PortBinding{}
	for _, port := range portMappings {
		exteriorPort := port.HostPort
		if exteriorPort == 0 {
			// No need to do port binding when HostPort is not specified
			continue
		}
		interiorPort := port.ContainerPort
		// Some of this port stuff is under-documented voodoo.
		// See http://stackoverflow.com/questions/20428302/binding-a-port-to-a-host-interface-using-the-rest-api
		var protocol string
		switch strings.ToUpper(string(port.Protocol)) {
		case "UDP":
			protocol = "/udp"
		case "TCP":
			protocol = "/tcp"
		default:
			glog.Warningf("Unknown protocol %q: defaulting to TCP", port.Protocol)
			protocol = "/tcp"
		}
		dockerPort := docker.Port(strconv.Itoa(interiorPort) + protocol)
		exposedPorts[dockerPort] = struct{}{}
		portBindings[dockerPort] = []docker.PortBinding{
			{
				HostPort: strconv.Itoa(exteriorPort),
				HostIP:   port.HostIP,
			},
		}
	}
	return exposedPorts, portBindings
}

func (dm *DockerManager) runContainer(
	pod *api.Pod,
	container *api.Container,
	opts *kubecontainer.RunContainerOptions,
	ref *api.ObjectReference,
	netMode string,
	ipcMode string) (string, error) {

	dockerName := KubeletContainerName{
		PodFullName:   kubecontainer.GetPodFullName(pod),
		PodUID:        pod.UID,
		ContainerName: container.Name,
	}
	exposedPorts, portBindings := makePortsAndBindings(opts.PortMappings)

	// TODO(vmarmol): Handle better.
	// Cap hostname at 63 chars (specification is 64bytes which is 63 chars and the null terminating char).
	const hostnameMaxLen = 63
	containerHostname := pod.Name
	if len(containerHostname) > hostnameMaxLen {
		containerHostname = containerHostname[:hostnameMaxLen]
	}
	namespacedName := types.NamespacedName{pod.Namespace, pod.Name}
	labels := map[string]string{
		"io.kubernetes.pod.name": namespacedName.String(),
	}
	if container.Lifecycle != nil && container.Lifecycle.PreStop != nil {
		glog.V(1).Infof("Setting preStop hook")
		// TODO: This is kind of hacky, we should really just encode the bits we need.
		data, err := latest.Codec.Encode(pod)
		if err != nil {
			glog.Errorf("Failed to encode pod: %s for prestop hook", pod.Name)
		} else {
			labels[kubernetesPodLabel] = string(data)
			labels[kubernetesContainerLabel] = container.Name
		}
	}
	memoryLimit := container.Resources.Limits.Memory().Value()
	cpuRequest := container.Resources.Requests.Cpu()
	cpuLimit := container.Resources.Limits.Cpu()
	var cpuShares int64
	// If request is not specified, but limit is, we want request to default to limit.
	// API server does this for new containers, but we repeat this logic in Kubelet
	// for containers running on existing Kubernetes clusters.
	if cpuRequest.Amount == nil && cpuLimit.Amount != nil {
		cpuShares = milliCPUToShares(cpuLimit.MilliValue())
	} else {
		// if cpuRequest.Amount is nil, then milliCPUToShares will return the minimal number
		// of CPU shares.
		cpuShares = milliCPUToShares(cpuRequest.MilliValue())
	}
	dockerOpts := docker.CreateContainerOptions{
		Name: BuildDockerName(dockerName, container),
		Config: &docker.Config{
			Env:          makeEnvList(opts.Envs),
			ExposedPorts: exposedPorts,
			Hostname:     containerHostname,
			Image:        container.Image,
			// Memory and CPU are set here for older versions of Docker (pre-1.6).
			Memory:     memoryLimit,
			MemorySwap: -1,
			CPUShares:  cpuShares,
			WorkingDir: container.WorkingDir,
			Labels:     labels,
			// Interactive containers:
			OpenStdin: container.Stdin,
			Tty:       container.TTY,
		},
	}

	setEntrypointAndCommand(container, opts, &dockerOpts)

	glog.V(3).Infof("Container %v/%v/%v: setting entrypoint \"%v\" and command \"%v\"", pod.Namespace, pod.Name, container.Name, dockerOpts.Config.Entrypoint, dockerOpts.Config.Cmd)

	securityContextProvider := securitycontext.NewSimpleSecurityContextProvider()
	securityContextProvider.ModifyContainerConfig(pod, container, dockerOpts.Config)
	dockerContainer, err := dm.client.CreateContainer(dockerOpts)
	if err != nil {
		if ref != nil {
			dm.recorder.Eventf(ref, "failed", "Failed to create docker container with error: %v", err)
		}
		return "", err
	}

	if ref != nil {
		dm.recorder.Eventf(ref, "created", "Created with docker id %v", util.ShortenString(dockerContainer.ID, 12))
	}

	binds := makeMountBindings(opts.Mounts)

	// The reason we create and mount the log file in here (not in kubelet) is because
	// the file's location depends on the ID of the container, and we need to create and
	// mount the file before actually starting the container.
	// TODO(yifan): Consider to pull this logic out since we might need to reuse it in
	// other container runtime.
	if opts.PodContainerDir != "" && len(container.TerminationMessagePath) != 0 {
		containerLogPath := path.Join(opts.PodContainerDir, dockerContainer.ID)
		fs, err := os.Create(containerLogPath)
		if err != nil {
			// TODO: Clean up the previouly created dir? return the error?
			glog.Errorf("Error on creating termination-log file %q: %v", containerLogPath, err)
		} else {
			fs.Close() // Close immediately; we're just doing a `touch` here
			b := fmt.Sprintf("%s:%s", containerLogPath, container.TerminationMessagePath)
			binds = append(binds, b)
		}
	}

	hc := &docker.HostConfig{
		PortBindings: portBindings,
		Binds:        binds,
		NetworkMode:  netMode,
		IpcMode:      ipcMode,
		// Memory and CPU are set here for newer versions of Docker (1.6+).
		Memory:     memoryLimit,
		MemorySwap: -1,
		CPUShares:  cpuShares,
	}
	if len(opts.DNS) > 0 {
		hc.DNS = opts.DNS
	}
	if len(opts.DNSSearch) > 0 {
		hc.DNSSearch = opts.DNSSearch
	}
	if len(opts.CgroupParent) > 0 {
		hc.CgroupParent = opts.CgroupParent
	}
	securityContextProvider.ModifyHostConfig(pod, container, hc)

	if err = dm.client.StartContainer(dockerContainer.ID, hc); err != nil {
		if ref != nil {
			dm.recorder.Eventf(ref, "failed",
				"Failed to start with docker id %v with error: %v", util.ShortenString(dockerContainer.ID, 12), err)
		}
		return "", err
	}
	if ref != nil {
		dm.recorder.Eventf(ref, "started", "Started with docker id %v", util.ShortenString(dockerContainer.ID, 12))
	}
	return dockerContainer.ID, nil
}

func setEntrypointAndCommand(container *api.Container, opts *kubecontainer.RunContainerOptions, dockerOpts *docker.CreateContainerOptions) {
	command, args := kubecontainer.ExpandContainerCommandAndArgs(container, opts.Envs)

	dockerOpts.Config.Entrypoint = command
	dockerOpts.Config.Cmd = args
}

// A helper function to get the KubeletContainerName and hash from a docker
// container.
func getDockerContainerNameInfo(c *docker.APIContainers) (*KubeletContainerName, uint64, error) {
	if len(c.Names) == 0 {
		return nil, 0, fmt.Errorf("cannot parse empty docker container name: %#v", c.Names)
	}
	dockerName, hash, err := ParseDockerName(c.Names[0])
	if err != nil {
		return nil, 0, fmt.Errorf("parse docker container name %q error: %v", c.Names[0], err)
	}
	return dockerName, hash, nil
}

// Get pod UID, name, and namespace by examining the container names.
func getPodInfoFromContainer(c *docker.APIContainers) (types.UID, string, string, error) {
	dockerName, _, err := getDockerContainerNameInfo(c)
	if err != nil {
		return types.UID(""), "", "", err
	}
	name, namespace, err := kubecontainer.ParsePodFullName(dockerName.PodFullName)
	if err != nil {
		return types.UID(""), "", "", fmt.Errorf("parse pod full name %q error: %v", dockerName.PodFullName, err)
	}
	return dockerName.PodUID, name, namespace, nil
}

// GetContainers returns a list of running containers if |all| is false;
// otherwise, it returns all containers.
func (dm *DockerManager) GetContainers(all bool) ([]*kubecontainer.Container, error) {
	containers, err := GetKubeletDockerContainers(dm.client, all)
	if err != nil {
		return nil, err
	}
	// Convert DockerContainers to []*kubecontainer.Container
	result := make([]*kubecontainer.Container, 0, len(containers))
	for _, c := range containers {
		converted, err := toRuntimeContainer(c)
		if err != nil {
			glog.Errorf("Error examining the container: %v", err)
			continue
		}
		result = append(result, converted)
	}
	return result, nil
}

func (dm *DockerManager) GetPods(all bool) ([]*kubecontainer.Pod, error) {
	start := time.Now()
	defer func() {
		metrics.ContainerManagerLatency.WithLabelValues("GetPods").Observe(metrics.SinceInMicroseconds(start))
	}()
	pods := make(map[types.UID]*kubecontainer.Pod)
	var result []*kubecontainer.Pod

	containers, err := GetKubeletDockerContainers(dm.client, all)
	if err != nil {
		return nil, err
	}

	// Group containers by pod.
	for _, c := range containers {
		converted, err := toRuntimeContainer(c)
		if err != nil {
			glog.Errorf("Error examining the container: %v", err)
			continue
		}

		podUID, podName, podNamespace, err := getPodInfoFromContainer(c)
		if err != nil {
			glog.Errorf("Error examining the container: %v", err)
			continue
		}

		pod, found := pods[podUID]
		if !found {
			pod = &kubecontainer.Pod{
				ID:        podUID,
				Name:      podName,
				Namespace: podNamespace,
			}
			pods[podUID] = pod
		}
		pod.Containers = append(pod.Containers, converted)
	}

	// Convert map to list.
	for _, c := range pods {
		result = append(result, c)
	}
	return result, nil
}

// List all images in the local storage.
func (dm *DockerManager) ListImages() ([]kubecontainer.Image, error) {
	var images []kubecontainer.Image

	dockerImages, err := dm.client.ListImages(docker.ListImagesOptions{})
	if err != nil {
		return images, err
	}

	for _, di := range dockerImages {
		image, err := toRuntimeImage(&di)
		if err != nil {
			continue
		}
		images = append(images, *image)
	}
	return images, nil
}

// TODO(vmarmol): Consider unexporting.
// PullImage pulls an image from network to local storage.
func (dm *DockerManager) PullImage(image kubecontainer.ImageSpec, secrets []api.Secret) error {
	return dm.puller.Pull(image.Image, secrets)
}

// IsImagePresent checks whether the container image is already in the local storage.
func (dm *DockerManager) IsImagePresent(image kubecontainer.ImageSpec) (bool, error) {
	return dm.puller.IsImagePresent(image.Image)
}

// Removes the specified image.
func (dm *DockerManager) RemoveImage(image kubecontainer.ImageSpec) error {
	return dm.client.RemoveImage(image.Image)
}

// podInfraContainerChanged returns true if the pod infra container has changed.
func (dm *DockerManager) podInfraContainerChanged(pod *api.Pod, podInfraContainer *kubecontainer.Container) (bool, error) {
	networkMode := ""
	var ports []api.ContainerPort

	dockerPodInfraContainer, err := dm.client.InspectContainer(string(podInfraContainer.ID))
	if err != nil {
		return false, err
	}

	// Check network mode.
	if dockerPodInfraContainer.HostConfig != nil {
		networkMode = dockerPodInfraContainer.HostConfig.NetworkMode
	}
	if pod.Spec.HostNetwork {
		if networkMode != "host" {
			glog.V(4).Infof("host: %v, %v", pod.Spec.HostNetwork, networkMode)
			return true, nil
		}
	} else {
		// Docker only exports ports from the pod infra container. Let's
		// collect all of the relevant ports and export them.
		for _, container := range pod.Spec.Containers {
			ports = append(ports, container.Ports...)
		}
	}
	expectedPodInfraContainer := &api.Container{
		Name:            PodInfraContainerName,
		Image:           dm.podInfraContainerImage,
		Ports:           ports,
		ImagePullPolicy: podInfraContainerImagePullPolicy,
	}
	return podInfraContainer.Hash != kubecontainer.HashContainer(expectedPodInfraContainer), nil
}

type dockerVersion docker.APIVersion

func NewVersion(input string) (dockerVersion, error) {
	version, err := docker.NewAPIVersion(input)
	return dockerVersion(version), err
}

func (dv dockerVersion) String() string {
	return docker.APIVersion(dv).String()
}

func (dv dockerVersion) Compare(other string) (int, error) {
	a := docker.APIVersion(dv)
	b, err := docker.NewAPIVersion(other)
	if err != nil {
		return 0, err
	}
	if a.LessThan(b) {
		return -1, nil
	}
	if a.GreaterThan(b) {
		return 1, nil
	}
	return 0, nil
}

func (dm *DockerManager) Version() (kubecontainer.Version, error) {
	env, err := dm.client.Version()
	if err != nil {
		return nil, fmt.Errorf("docker: failed to get docker version: %v", err)
	}

	apiVersion := env.Get("ApiVersion")
	version, err := docker.NewAPIVersion(apiVersion)
	if err != nil {
		glog.Errorf("docker: failed to parse docker server version %q: %v", apiVersion, err)
		return nil, fmt.Errorf("docker: failed to parse docker server version %q: %v", apiVersion, err)
	}
	return dockerVersion(version), nil
}

// The first version of docker that supports exec natively is 1.3.0 == API 1.15
var dockerAPIVersionWithExec = "1.15"

func (dm *DockerManager) nativeExecSupportExists() (bool, error) {
	version, err := dm.Version()
	if err != nil {
		return false, err
	}
	result, err := version.Compare(dockerAPIVersionWithExec)
	if result >= 0 {
		return true, err
	}
	return false, err
}

func (dm *DockerManager) getRunInContainerCommand(containerID string, cmd []string) (*exec.Cmd, error) {
	args := append([]string{"exec"}, cmd...)
	command := exec.Command("/usr/sbin/nsinit", args...)
	command.Dir = fmt.Sprintf("/var/lib/docker/execdriver/native/%s", containerID)
	return command, nil
}

func (dm *DockerManager) runInContainerUsingNsinit(containerID string, cmd []string) ([]byte, error) {
	c, err := dm.getRunInContainerCommand(containerID, cmd)
	if err != nil {
		return nil, err
	}
	return c.CombinedOutput()
}

// RunInContainer uses nsinit to run the command inside the container identified by containerID
// TODO(yifan): Use strong type for containerID.
func (dm *DockerManager) RunInContainer(containerID string, cmd []string) ([]byte, error) {
	// If native exec support does not exist in the local docker daemon use nsinit.
	useNativeExec, err := dm.nativeExecSupportExists()
	if err != nil {
		return nil, err
	}
	if !useNativeExec {
		glog.V(2).Infof("Using nsinit to run the command %+v inside container %s", cmd, containerID)
		return dm.runInContainerUsingNsinit(containerID, cmd)
	}
	glog.V(2).Infof("Using docker native exec to run cmd %+v inside container %s", cmd, containerID)
	createOpts := docker.CreateExecOptions{
		Container:    containerID,
		Cmd:          cmd,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
	}
	execObj, err := dm.client.CreateExec(createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to run in container - Exec setup failed - %v", err)
	}
	var buf bytes.Buffer
	startOpts := docker.StartExecOptions{
		Detach:       false,
		Tty:          false,
		OutputStream: &buf,
		ErrorStream:  &buf,
		RawTerminal:  false,
	}
	err = dm.client.StartExec(execObj.ID, startOpts)
	if err != nil {
		glog.V(2).Infof("StartExec With error: %v", err)
		return nil, err
	}
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		inspect, err2 := dm.client.InspectExec(execObj.ID)
		if err2 != nil {
			glog.V(2).Infof("InspectExec %s failed with error: %+v", execObj.ID, err2)
			return buf.Bytes(), err2
		}
		if !inspect.Running {
			if inspect.ExitCode != 0 {
				glog.V(2).Infof("InspectExec %s exit with result %+v", execObj.ID, inspect)
				err = &dockerExitError{inspect}
			}
			break
		}
		<-ticker.C
	}

	return buf.Bytes(), err
}

type dockerExitError struct {
	Inspect *docker.ExecInspect
}

func (d *dockerExitError) String() string {
	return d.Error()
}

func (d *dockerExitError) Error() string {
	return fmt.Sprintf("Error executing in Docker Container: %d", d.Inspect.ExitCode)
}

func (d *dockerExitError) Exited() bool {
	return !d.Inspect.Running
}

func (d *dockerExitError) ExitStatus() int {
	return d.Inspect.ExitCode
}

// ExecInContainer runs the command inside the container identified by containerID.
//
// TODO:
//  - use strong type for containerId
func (dm *DockerManager) ExecInContainer(containerId string, cmd []string, stdin io.Reader, stdout, stderr io.WriteCloser, tty bool) error {
	if dm.execHandler == nil {
		return errors.New("unable to exec without an exec handler")
	}

	container, err := dm.client.InspectContainer(containerId)
	if err != nil {
		return err
	}
	if !container.State.Running {
		return fmt.Errorf("container not running (%s)", container)
	}

	return dm.execHandler.ExecInContainer(dm.client, container, cmd, stdin, stdout, stderr, tty)
}

func (dm *DockerManager) AttachContainer(containerId string, stdin io.Reader, stdout, stderr io.WriteCloser, tty bool) error {
	opts := docker.AttachToContainerOptions{
		Container:    containerId,
		InputStream:  stdin,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Stream:       true,
		Logs:         true,
		Stdin:        stdin != nil,
		Stdout:       stdout != nil,
		Stderr:       stderr != nil,
		RawTerminal:  tty,
	}
	return dm.client.AttachToContainer(opts)
}

func noPodInfraContainerError(podName, podNamespace string) error {
	return fmt.Errorf("cannot find pod infra container in pod %q", kubecontainer.BuildPodFullName(podName, podNamespace))
}

// PortForward executes socat in the pod's network namespace and copies
// data between stream (representing the user's local connection on their
// computer) and the specified port in the container.
//
// TODO:
//  - match cgroups of container
//  - should we support nsenter + socat on the host? (current impl)
//  - should we support nsenter + socat in a container, running with elevated privs and --pid=host?
func (dm *DockerManager) PortForward(pod *kubecontainer.Pod, port uint16, stream io.ReadWriteCloser) error {
	podInfraContainer := pod.FindContainerByName(PodInfraContainerName)
	if podInfraContainer == nil {
		return noPodInfraContainerError(pod.Name, pod.Namespace)
	}
	container, err := dm.client.InspectContainer(string(podInfraContainer.ID))
	if err != nil {
		return err
	}

	if !container.State.Running {
		return fmt.Errorf("container not running (%s)", container)
	}

	containerPid := container.State.Pid
	// TODO what if the host doesn't have it???
	_, lookupErr := exec.LookPath("socat")
	if lookupErr != nil {
		return fmt.Errorf("Unable to do port forwarding: socat not found.")
	}
	args := []string{"-t", fmt.Sprintf("%d", containerPid), "-n", "socat", "-", fmt.Sprintf("TCP4:localhost:%d", port)}
	// TODO use exec.LookPath
	command := exec.Command("nsenter", args...)
	command.Stdin = stream
	command.Stdout = stream
	return command.Run()
}

// Kills all containers in the specified pod
func (dm *DockerManager) KillPod(pod kubecontainer.Pod) error {
	// Send the kills in parallel since they may take a long time. Len + 1 since there
	// can be Len errors + the networkPlugin teardown error.
	errs := make(chan error, len(pod.Containers)+1)
	wg := sync.WaitGroup{}
	var networkID types.UID
	for _, container := range pod.Containers {
		wg.Add(1)
		go func(container *kubecontainer.Container) {
			defer util.HandleCrash()
			defer wg.Done()

			// TODO: Handle this without signaling the pod infra container to
			// adapt to the generic container runtime.
			if container.Name == PodInfraContainerName {
				// Store the container runtime for later deletion.
				// We do this so that PreStop handlers can run in the network namespace.
				networkID = container.ID
				return
			}
			if err := dm.killContainer(container.ID); err != nil {
				glog.Errorf("Failed to delete container: %v; Skipping pod %q", err, pod.ID)
				errs <- err
			}
		}(container)
	}
	wg.Wait()
	if len(networkID) > 0 {
		if err := dm.networkPlugin.TearDownPod(pod.Namespace, pod.Name, kubeletTypes.DockerID(networkID)); err != nil {
			glog.Errorf("Failed tearing down the infra container: %v", err)
			errs <- err
		}
		if err := dm.killContainer(networkID); err != nil {
			glog.Errorf("Failed to delete container: %v; Skipping pod %q", err, pod.ID)
			errs <- err
		}
	}
	close(errs)
	if len(errs) > 0 {
		errList := []error{}
		for err := range errs {
			errList = append(errList, err)
		}
		return fmt.Errorf("failed to delete containers (%v)", errList)
	}
	return nil
}

// KillContainerInPod kills a container in the pod.
func (dm *DockerManager) KillContainerInPod(container api.Container, pod *api.Pod) error {
	// Locate the container.
	pods, err := dm.GetPods(false)
	if err != nil {
		return err
	}
	targetPod := kubecontainer.Pods(pods).FindPod(kubecontainer.GetPodFullName(pod), pod.UID)
	targetContainer := targetPod.FindContainerByName(container.Name)
	if targetContainer == nil {
		return fmt.Errorf("unable to find container %q in pod %q", container.Name, targetPod.Name)
	}
	return dm.killContainer(targetContainer.ID)
}

// TODO(vmarmol): Unexport this as it is no longer used externally.
// KillContainer kills a container identified by containerID.
// Internally, it invokes docker's StopContainer API with a timeout of 10s.
// TODO: Deprecate this function in favor of KillContainerInPod.
func (dm *DockerManager) KillContainer(containerID types.UID) error {
	return dm.killContainer(containerID)
}

func (dm *DockerManager) killContainer(containerID types.UID) error {
	ID := string(containerID)
	glog.V(2).Infof("Killing container with id %q", ID)
	inspect, err := dm.client.InspectContainer(ID)
	if err != nil {
		return err
	}
	var found bool
	var preStop string
	if inspect != nil && inspect.Config != nil && inspect.Config.Labels != nil {
		preStop, found = inspect.Config.Labels[kubernetesPodLabel]
	}
	if found {
		var pod api.Pod
		err := latest.Codec.DecodeInto([]byte(preStop), &pod)
		if err != nil {
			glog.Errorf("Failed to decode prestop: %s, %s", preStop, ID)
		} else {
			name := inspect.Config.Labels[kubernetesContainerLabel]
			var container *api.Container
			for ix := range pod.Spec.Containers {
				if pod.Spec.Containers[ix].Name == name {
					container = &pod.Spec.Containers[ix]
					break
				}
			}
			if container != nil {
				glog.V(1).Infof("Running preStop hook")
				if err := dm.runner.Run(ID, &pod, container, container.Lifecycle.PreStop); err != nil {
					glog.Errorf("failed to run preStop hook: %v", err)
				}
			} else {
				glog.Errorf("unable to find container %v, %s", pod, name)
			}
		}
	}
	dm.readinessManager.RemoveReadiness(ID)
	err = dm.client.StopContainer(ID, 10)
	ref, ok := dm.containerRefManager.GetRef(ID)
	if !ok {
		glog.Warningf("No ref for pod '%v'", ID)
	} else {
		// TODO: pass reason down here, and state, or move this call up the stack.
		dm.recorder.Eventf(ref, "killing", "Killing with docker id %v", util.ShortenString(ID, 12))
	}
	return err
}

// Run a single container from a pod. Returns the docker container ID
func (dm *DockerManager) runContainerInPod(pod *api.Pod, container *api.Container, netMode, ipcMode string) (kubeletTypes.DockerID, error) {
	start := time.Now()
	defer func() {
		metrics.ContainerManagerLatency.WithLabelValues("runContainerInPod").Observe(metrics.SinceInMicroseconds(start))
	}()

	ref, err := kubecontainer.GenerateContainerRef(pod, container)
	if err != nil {
		glog.Errorf("Couldn't make a ref to pod %v, container %v: '%v'", pod.Name, container.Name, err)
	}

	opts, err := dm.generator.GenerateRunContainerOptions(pod, container)
	if err != nil {
		return "", err
	}

	id, err := dm.runContainer(pod, container, opts, ref, netMode, ipcMode)
	if err != nil {
		return "", err
	}

	// Remember this reference so we can report events about this container
	if ref != nil {
		dm.containerRefManager.SetRef(id, ref)
	}

	if container.Lifecycle != nil && container.Lifecycle.PostStart != nil {
		handlerErr := dm.runner.Run(id, pod, container, container.Lifecycle.PostStart)
		if handlerErr != nil {
			dm.killContainer(types.UID(id))
			return kubeletTypes.DockerID(""), fmt.Errorf("failed to call event handler: %v", handlerErr)
		}
	}

	// Create a symbolic link to the Docker container log file using a name which captures the
	// full pod name, the container name and the Docker container ID. Cluster level logging will
	// capture these symbolic filenames which can be used for search terms in Elasticsearch or for
	// labels for Cloud Logging.
	podFullName := kubecontainer.GetPodFullName(pod)
	containerLogFile := path.Join(dm.dockerRoot, "containers", id, fmt.Sprintf("%s-json.log", id))
	symlinkFile := LogSymlink(dm.containerLogsDir, podFullName, container.Name, id)
	if err = dm.os.Symlink(containerLogFile, symlinkFile); err != nil {
		glog.Errorf("Failed to create symbolic link to the log file of pod %q container %q: %v", podFullName, container.Name, err)
	}

	// Container information is used in adjusting OOM scores and adding ndots.
	containerInfo, err := dm.client.InspectContainer(string(id))
	if err != nil {
		return "", err
	}
	// Ensure the PID actually exists, else we'll move ourselves.
	if containerInfo.State.Pid == 0 {
		return "", fmt.Errorf("failed to get init PID for Docker container %q", string(id))
	}

	// Set OOM score of the container based on the priority of the container.
	// Processes in lower-priority pods should be killed first if the system runs out of memory.
	// The main pod infrastructure container is considered high priority, since if it is killed the
	// whole pod will die.
	var oomScoreAdj int
	if container.Name == PodInfraContainerName {
		oomScoreAdj = qos.PodInfraOomAdj
	} else {
		oomScoreAdj = qos.GetContainerOomScoreAdjust(container, dm.machineInfo.MemoryCapacity)
	}
	cgroupName, err := dm.procFs.GetFullContainerName(containerInfo.State.Pid)
	if err != nil {
		return "", err
	}
	if err = dm.oomAdjuster.ApplyOomScoreAdjContainer(cgroupName, oomScoreAdj, 5); err != nil {
		return "", err
	}

	// currently, Docker does not have a flag by which the ndots option can be passed.
	// (A separate issue has been filed with Docker to add a ndots flag)
	// The addNDotsOption call appends the ndots option to the resolv.conf file generated by docker.
	// This resolv.conf file is shared by all containers of the same pod, and needs to be modified only once per pod.
	// we modify it when the pause container is created since it is the first container created in the pod since it holds
	// the networking namespace.
	if container.Name == PodInfraContainerName {
		err = addNDotsOption(containerInfo.ResolvConfPath)
	}

	return kubeletTypes.DockerID(id), err
}

func addNDotsOption(resolvFilePath string) error {
	if len(resolvFilePath) == 0 {
		glog.Errorf("DNS ResolvConfPath is empty.")
		return nil
	}

	if _, err := os.Stat(resolvFilePath); os.IsNotExist(err) {
		return fmt.Errorf("DNS ResolvConfPath specified but does not exist. It could not be updated: %s", resolvFilePath)
	}

	glog.V(4).Infof("DNS ResolvConfPath exists: %s. Will attempt to add ndots option: %s", resolvFilePath, ndotsDNSOption)

	if err := appendToFile(resolvFilePath, ndotsDNSOption); err != nil {
		glog.Errorf("resolv.conf could not be updated. err:%v", err)
		return err
	}
	return nil
}

func appendToFile(filePath, stringToAppend string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(stringToAppend)
	return err
}

// createPodInfraContainer starts the pod infra container for a pod. Returns the docker container ID of the newly created container.
func (dm *DockerManager) createPodInfraContainer(pod *api.Pod) (kubeletTypes.DockerID, error) {
	start := time.Now()
	defer func() {
		metrics.ContainerManagerLatency.WithLabelValues("createPodInfraContainer").Observe(metrics.SinceInMicroseconds(start))
	}()
	// Use host networking if specified.
	netNamespace := ""
	var ports []api.ContainerPort

	if pod.Spec.HostNetwork {
		netNamespace = "host"
	} else {
		// Docker only exports ports from the pod infra container.  Let's
		// collect all of the relevant ports and export them.
		for _, container := range pod.Spec.Containers {
			ports = append(ports, container.Ports...)
		}
	}

	container := &api.Container{
		Name:            PodInfraContainerName,
		Image:           dm.podInfraContainerImage,
		Ports:           ports,
		ImagePullPolicy: podInfraContainerImagePullPolicy,
	}

	// No pod secrets for the infra container.
	if err := dm.pullImage(pod, container, nil); err != nil {
		return "", err
	}

	id, err := dm.runContainerInPod(pod, container, netNamespace, "")
	if err != nil {
		return "", err
	}

	return id, nil
}

// TODO(vmarmol): This will soon be made non-public when its only use is internal.
// Structure keeping information on changes that need to happen for a pod. The semantics is as follows:
// - startInfraContainer is true if new Infra Containers have to be started and old one (if running) killed.
//   Additionally if it is true then containersToKeep have to be empty
// - infraContainerId have to be set iff startInfraContainer is false. It stores dockerID of running Infra Container
// - containersToStart keeps indices of Specs of containers that have to be started.
// - containersToKeep stores mapping from dockerIDs of running containers to indices of their Specs for containers that
//   should be kept running. If startInfraContainer is false then it contains an entry for infraContainerId (mapped to -1).
//   It shouldn't be the case where containersToStart is empty and containersToKeep contains only infraContainerId. In such case
//   Infra Container should be killed, hence it's removed from this map.
// - all running containers which are NOT contained in containersToKeep should be killed.
type empty struct{}
type PodContainerChangesSpec struct {
	StartInfraContainer bool
	InfraContainerId    kubeletTypes.DockerID
	ContainersToStart   map[int]empty
	ContainersToKeep    map[kubeletTypes.DockerID]int
}

func (dm *DockerManager) computePodContainerChanges(pod *api.Pod, runningPod kubecontainer.Pod, podStatus api.PodStatus) (PodContainerChangesSpec, error) {
	start := time.Now()
	defer func() {
		metrics.ContainerManagerLatency.WithLabelValues("computePodContainerChanges").Observe(metrics.SinceInMicroseconds(start))
	}()

	podFullName := kubecontainer.GetPodFullName(pod)
	uid := pod.UID
	glog.V(4).Infof("Syncing Pod %+v, podFullName: %q, uid: %q", pod, podFullName, uid)

	containersToStart := make(map[int]empty)
	containersToKeep := make(map[kubeletTypes.DockerID]int)
	createPodInfraContainer := false

	var err error
	var podInfraContainerID kubeletTypes.DockerID
	var changed bool
	podInfraContainer := runningPod.FindContainerByName(PodInfraContainerName)
	if podInfraContainer != nil {
		glog.V(4).Infof("Found pod infra container for %q", podFullName)
		changed, err = dm.podInfraContainerChanged(pod, podInfraContainer)
		if err != nil {
			return PodContainerChangesSpec{}, err
		}
	}

	createPodInfraContainer = true
	if podInfraContainer == nil {
		glog.V(2).Infof("Need to restart pod infra container for %q because it is not found", podFullName)
	} else if changed {
		glog.V(2).Infof("Need to restart pod infra container for %q because it is changed", podFullName)
	} else {
		glog.V(4).Infof("Pod infra container looks good, keep it %q", podFullName)
		createPodInfraContainer = false
		podInfraContainerID = kubeletTypes.DockerID(podInfraContainer.ID)
		containersToKeep[podInfraContainerID] = -1
	}

	for index, container := range pod.Spec.Containers {
		expectedHash := kubecontainer.HashContainer(&container)

		c := runningPod.FindContainerByName(container.Name)
		if c == nil {
			if kubecontainer.ShouldContainerBeRestarted(&container, pod, &podStatus, dm.readinessManager) {
				// If we are here it means that the container is dead and should be restarted, or never existed and should
				// be created. We may be inserting this ID again if the container has changed and it has
				// RestartPolicy::Always, but it's not a big deal.
				glog.V(3).Infof("Container %+v is dead, but RestartPolicy says that we should restart it.", container)
				containersToStart[index] = empty{}
			}
			continue
		}

		containerID := kubeletTypes.DockerID(c.ID)
		hash := c.Hash
		glog.V(3).Infof("pod %q container %q exists as %v", podFullName, container.Name, containerID)

		if createPodInfraContainer {
			// createPodInfraContainer == true and Container exists
			// If we're creating infra containere everything will be killed anyway
			// If RestartPolicy is Always or OnFailure we restart containers that were running before we
			// killed them when restarting Infra Container.
			if pod.Spec.RestartPolicy != api.RestartPolicyNever {
				glog.V(1).Infof("Infra Container is being recreated. %q will be restarted.", container.Name)
				containersToStart[index] = empty{}
			}
			continue
		}

		// At this point, the container is running and pod infra container is good.
		// We will look for changes and check healthiness for the container.
		containerChanged := hash != 0 && hash != expectedHash
		if containerChanged {
			glog.Infof("pod %q container %q hash changed (%d vs %d), it will be killed and re-created.", podFullName, container.Name, hash, expectedHash)
			containersToStart[index] = empty{}
			continue
		}

		result, err := dm.prober.Probe(pod, podStatus, container, string(c.ID), c.Created)
		if err != nil {
			// TODO(vmarmol): examine this logic.
			glog.V(2).Infof("probe no-error: %q", container.Name)
			containersToKeep[containerID] = index
			continue
		}
		if result == probe.Success {
			glog.V(4).Infof("probe success: %q", container.Name)
			containersToKeep[containerID] = index
			continue
		}
		glog.Infof("pod %q container %q is unhealthy (probe result: %v), it will be killed and re-created.", podFullName, container.Name, result)
		containersToStart[index] = empty{}
	}

	// After the loop one of the following should be true:
	// - createPodInfraContainer is true and containersToKeep is empty.
	// (In fact, when createPodInfraContainer is false, containersToKeep will not be touched).
	// - createPodInfraContainer is false and containersToKeep contains at least ID of Infra Container

	// If Infra container is the last running one, we don't want to keep it.
	if !createPodInfraContainer && len(containersToStart) == 0 && len(containersToKeep) == 1 {
		containersToKeep = make(map[kubeletTypes.DockerID]int)
	}

	return PodContainerChangesSpec{
		StartInfraContainer: createPodInfraContainer,
		InfraContainerId:    podInfraContainerID,
		ContainersToStart:   containersToStart,
		ContainersToKeep:    containersToKeep,
	}, nil
}

// updateReasonCache updates the failure reason based on the latest error.
func (dm *DockerManager) updateReasonCache(pod *api.Pod, container *api.Container, err error) {
	if err == nil {
		return
	}
	errString := err.Error()
	dm.reasonCache.Add(pod.UID, container.Name, errString)
}

// clearReasonCache removes the entry in the reason cache.
func (dm *DockerManager) clearReasonCache(pod *api.Pod, container *api.Container) {
	dm.reasonCache.Remove(pod.UID, container.Name)
}

// Pull the image for the specified pod and container.
func (dm *DockerManager) pullImage(pod *api.Pod, container *api.Container, pullSecrets []api.Secret) error {
	ref, err := kubecontainer.GenerateContainerRef(pod, container)
	if err != nil {
		glog.Errorf("Couldn't make a ref to pod %v, container %v: '%v'", pod.Name, container.Name, err)
	}
	spec := kubecontainer.ImageSpec{container.Image}
	present, err := dm.IsImagePresent(spec)
	if err != nil {
		if ref != nil {
			dm.recorder.Eventf(ref, "failed", "Failed to inspect image %q: %v", container.Image, err)
		}
		return fmt.Errorf("failed to inspect image %q: %v", container.Image, err)
	}
	if !dm.runtimeHooks.ShouldPullImage(pod, container, present) {
		if present && ref != nil {
			dm.recorder.Eventf(ref, "pulled", "Container image %q already present on machine", container.Image)
		}
		return nil
	}

	dm.runtimeHooks.ReportImagePulling(pod, container)
	err = dm.PullImage(spec, pullSecrets)
	dm.runtimeHooks.ReportImagePulled(pod, container, err)
	return err
}

// Sync the running pod to match the specified desired pod.
func (dm *DockerManager) SyncPod(pod *api.Pod, runningPod kubecontainer.Pod, podStatus api.PodStatus, pullSecrets []api.Secret) error {
	start := time.Now()
	defer func() {
		metrics.ContainerManagerLatency.WithLabelValues("SyncPod").Observe(metrics.SinceInMicroseconds(start))
	}()

	podFullName := kubecontainer.GetPodFullName(pod)
	containerChanges, err := dm.computePodContainerChanges(pod, runningPod, podStatus)
	glog.V(3).Infof("Got container changes for pod %q: %+v", podFullName, containerChanges)
	if err != nil {
		return err
	}

	if containerChanges.StartInfraContainer || (len(containerChanges.ContainersToKeep) == 0 && len(containerChanges.ContainersToStart) == 0) {
		if len(containerChanges.ContainersToKeep) == 0 && len(containerChanges.ContainersToStart) == 0 {
			glog.V(4).Infof("Killing Infra Container for %q because all other containers are dead.", podFullName)
		} else {
			glog.V(4).Infof("Killing Infra Container for %q, will start new one", podFullName)
		}

		// Killing phase: if we want to start new infra container, or nothing is running kill everything (including infra container)
		err = dm.KillPod(runningPod)
		if err != nil {
			return err
		}
	} else {
		// Otherwise kill any containers in this pod which are not specified as ones to keep.
		for _, container := range runningPod.Containers {
			_, keep := containerChanges.ContainersToKeep[kubeletTypes.DockerID(container.ID)]
			if !keep {
				glog.V(3).Infof("Killing unwanted container %+v", container)
				err = dm.KillContainer(container.ID)
				if err != nil {
					glog.Errorf("Error killing container: %v", err)
				}
			}
		}
	}

	// If we should create infra container then we do it first.
	podInfraContainerID := containerChanges.InfraContainerId
	if containerChanges.StartInfraContainer && (len(containerChanges.ContainersToStart) > 0) {
		glog.V(4).Infof("Creating pod infra container for %q", podFullName)
		podInfraContainerID, err = dm.createPodInfraContainer(pod)

		// Call the networking plugin
		if err == nil {
			err = dm.networkPlugin.SetUpPod(pod.Namespace, pod.Name, podInfraContainerID)
		}
		if err != nil {
			glog.Errorf("Failed to create pod infra container: %v; Skipping pod %q", err, podFullName)
			return err
		}
	}

	// Start everything
	for idx := range containerChanges.ContainersToStart {
		container := &pod.Spec.Containers[idx]
		glog.V(4).Infof("Creating container %+v in pod %v", container, podFullName)
		err := dm.pullImage(pod, container, pullSecrets)
		dm.updateReasonCache(pod, container, err)
		if err != nil {
			glog.Warningf("Failed to pull image %q from pod %q and container %q: %v", container.Image, kubecontainer.GetPodFullName(pod), container.Name, err)
			continue
		}

		if container.SecurityContext != nil && container.SecurityContext.RunAsNonRoot {
			err := dm.verifyNonRoot(container)
			dm.updateReasonCache(pod, container, err)
			if err != nil {
				glog.Errorf("Error running pod %q container %q: %v", kubecontainer.GetPodFullName(pod), container.Name, err)
				continue
			}
		}

		// TODO(dawnchen): Check RestartPolicy.DelaySeconds before restart a container
		namespaceMode := fmt.Sprintf("container:%v", podInfraContainerID)
		_, err = dm.runContainerInPod(pod, container, namespaceMode, namespaceMode)
		dm.updateReasonCache(pod, container, err)
		if err != nil {
			// TODO(bburns) : Perhaps blacklist a container after N failures?
			glog.Errorf("Error running pod %q container %q: %v", kubecontainer.GetPodFullName(pod), container.Name, err)
			continue
		}
		// Successfully started the container; clear the entry in the failure
		// reason cache.
		dm.clearReasonCache(pod, container)
	}

	return nil
}

// verifyNonRoot returns an error if the container or image will run as the root user.
func (dm *DockerManager) verifyNonRoot(container *api.Container) error {
	if securitycontext.HasRunAsUser(container) {
		if securitycontext.HasRootRunAsUser(container) {
			return fmt.Errorf("container's runAsUser breaks non-root policy")
		}
		return nil
	}

	imgRoot, err := dm.isImageRoot(container.Image)
	if err != nil {
		return err
	}
	if imgRoot {
		return fmt.Errorf("container has no runAsUser and image will run as root")
	}

	return nil
}

// isImageRoot returns true if the user directive is not set on the image, the user is set to 0
// or the user is set to root.  If there is an error inspecting the image this method will return
// false and return the error.
func (dm *DockerManager) isImageRoot(image string) (bool, error) {
	img, err := dm.client.InspectImage(image)
	if err != nil {
		return false, err
	}
	if img == nil || img.Config == nil {
		return false, fmt.Errorf("unable to inspect image %s, nil Config", image)
	}

	user := getUidFromUser(img.Config.User)
	// if no user is defined container will run as root
	if user == "" {
		return true, nil
	}
	// do not allow non-numeric user directives
	uid, err := strconv.Atoi(user)
	if err != nil {
		return false, fmt.Errorf("unable to validate image is non-root, non-numeric user (%s) is not allowed", user)
	}
	// user is numeric, check for 0
	return uid == 0, nil
}

// getUidFromUser splits the uid out of a uid:gid string.
func getUidFromUser(id string) string {
	if id == "" {
		return id
	}
	// split instances where the id may contain uid:gid
	if strings.Contains(id, ":") {
		return strings.Split(id, ":")[0]
	}
	// no gid, just return the id
	return id
}
