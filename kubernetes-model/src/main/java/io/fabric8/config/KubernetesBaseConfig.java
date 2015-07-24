/**
 * Copyright (C) 2011 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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
