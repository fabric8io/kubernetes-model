package io.fabric8.config;

import io.sundr.builder.annotations.ExternalBuildables;

@ExternalBuildables(validationEnabled = true, builderPackage = "io.fabric8.common", value = {
        "io.fabric8.kubernetes.api.model.base.Capabilities",
        "io.fabric8.kubernetes.api.model.base.Container",
        "io.fabric8.kubernetes.api.model.base.ContainerPort",
        "io.fabric8.kubernetes.api.model.base.EmptyDirVolumeSource",
        "io.fabric8.kubernetes.api.model.base.EnvVar",
        "io.fabric8.kubernetes.api.model.base.ExecAction",
        "io.fabric8.kubernetes.api.model.base.GCEPersistentDiskVolumeSource",
        "io.fabric8.kubernetes.api.model.base.GitRepoVolumeSource",
        "io.fabric8.kubernetes.api.model.base.GlusterfsVolumeSource",
        "io.fabric8.kubernetes.api.model.base.HTTPGetAction",
        "io.fabric8.kubernetes.api.model.base.Handler",
        "io.fabric8.kubernetes.api.model.base.HostPathVolumeSource",
        "io.fabric8.kubernetes.api.model.base.ISCSIVolumeSource",
        "io.fabric8.kubernetes.api.model.base.Lifecycle",
        "io.fabric8.kubernetes.api.model.base.NFSVolumeSource",
        "io.fabric8.kubernetes.api.model.base.ObjectReference",
        "io.fabric8.kubernetes.api.model.base.PodSpec",
        "io.fabric8.kubernetes.api.model.base.PodTemplateSpec",
        "io.fabric8.kubernetes.api.model.base.Probe",
        "io.fabric8.kubernetes.api.model.base.ReplicationControllerSpec",
        "io.fabric8.kubernetes.api.model.base.ResourceRequirements",
        "io.fabric8.kubernetes.api.model.base.SecretVolumeSource",
        "io.fabric8.kubernetes.api.model.base.Status",
        "io.fabric8.kubernetes.api.model.base.StatusCause",
        "io.fabric8.kubernetes.api.model.base.StatusDetails",
        "io.fabric8.kubernetes.api.model.base.TCPSocketAction",
        "io.fabric8.kubernetes.api.model.base.Volume",
        "io.fabric8.kubernetes.api.model.base.VolumeMount"
})
public class KubernetesBaseConfig {
}
