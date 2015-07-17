package io.fabric8.config;

import io.fabric8.kubernetes.api.model.Doneable;
import io.sundr.builder.annotations.ExternalBuildables;
import io.sundr.builder.annotations.Inline;

@ExternalBuildables(editableEnabled=false, validationEnabled = true, builderPackage = "io.fabric8.kubernetes.api.builder",
        inline = {
                @Inline(type = Doneable.class, prefix = "Doneable", value = "done")
        },
        value = {
        "io.fabric8.kubernetes.api.model.base.ListMeta",
        "io.fabric8.kubernetes.api.model.base.ObjectMeta",
        "io.fabric8.kubernetes.api.model.base.ObjectReference",
        "io.fabric8.kubernetes.api.model.base.Status",
        "io.fabric8.kubernetes.api.model.base.StatusCause",
        "io.fabric8.kubernetes.api.model.base.StatusDetails"
})
public class KubernetesBaseConfig {
}
