package io.fabric8.config;

import io.sundr.builder.annotations.ExternalBuildables;

@ExternalBuildables(validationEnabled = true, builderPackage = "io.fabric8.common", value = {
        "io.fabri8c.kubernetes.api.model.base.ObjectReference",
        "io.fabri8c.kubernetes.api.model.base.Status",
        "io.fabri8c.kubernetes.api.model.base.StatusCause",
        "io.fabri8c.kubernetes.api.model.base.StatusDetails"
})
public class KubernetesBaseConfig {
}
