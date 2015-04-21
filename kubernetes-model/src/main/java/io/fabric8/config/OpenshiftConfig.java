package io.fabric8.config;

import io.sundr.builder.annotations.ExternalBuildables;

@ExternalBuildables(builderPackage = "io.fabric8.common", value = {
        "io.fabric8.openshift.api.model.Build",
        "io.fabric8.openshift.api.model.BuildConfig",
        "io.fabric8.openshift.api.model.BuildConfigList",
        "io.fabric8.openshift.api.model.BuildList",
        "io.fabric8.openshift.api.model.BuildOutput",
        "io.fabric8.openshift.api.model.BuildParameters",
        "io.fabric8.openshift.api.model.BuildSource",
        "io.fabric8.openshift.api.model.BuildStrategy",
        "io.fabric8.openshift.api.model.BuildTriggerPolicy",
        "io.fabric8.openshift.api.model.CustomBuildStrategy",
        "io.fabric8.openshift.api.model.CustomDeploymentStrategyParams",
        "io.fabric8.openshift.api.model.Deployment",
        "io.fabric8.openshift.api.model.DeploymentCause",
        "io.fabric8.openshift.api.model.DeploymentCauseImageTrigger",
        "io.fabric8.openshift.api.model.DeploymentConfig",
        "io.fabric8.openshift.api.model.DeploymentConfigList",
        "io.fabric8.openshift.api.model.DeploymentDetails",
        "io.fabric8.openshift.api.model.DeploymentList",
        "io.fabric8.openshift.api.model.DeploymentStrategy",
        "io.fabric8.openshift.api.model.DeploymentTemplate",
        "io.fabric8.openshift.api.model.DeploymentTriggerImageChangeParams",
        "io.fabric8.openshift.api.model.DeploymentTriggerPolicy",
        "io.fabric8.openshift.api.model.DockerBuildStrategy",
        "io.fabric8.openshift.api.model.GitBuildSource",
        "io.fabric8.openshift.api.model.GitSourceRevision",
        "io.fabric8.openshift.api.model.Image",
        "io.fabric8.openshift.api.model.ImageList",
        "io.fabric8.openshift.api.model.ImageRepository",
        "io.fabric8.openshift.api.model.ImageRepositoryList",
        "io.fabric8.openshift.api.model.ImageRepositoryStatus",
        "io.fabric8.openshift.api.model.Route",
        "io.fabric8.openshift.api.model.RouteList",
        "io.fabric8.openshift.api.model.STIBuildStrategy",
        "io.fabric8.openshift.api.model.SourceControlUser",
        "io.fabric8.openshift.api.model.SourceRevision",
        "io.fabric8.openshift.api.model.WebHookTrigger",
        "io.fabric8.openshift.api.model.config.Config",
        "io.fabric8.openshift.api.model.template.Template",
        "io.fabric8.openshift.api.model.template.Parameter"
})
public class OpenshiftConfig {
}
