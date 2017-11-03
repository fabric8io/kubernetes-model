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
package io.fabric8.kubernetes.internal;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonDeserializer;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.node.ObjectNode;
import io.fabric8.kubernetes.api.model.KubernetesList;
import io.fabric8.kubernetes.api.model.KubernetesResource;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

public class KubernetesDeserializer extends JsonDeserializer<KubernetesResource> {

    private static final String KIND = "kind";

    private static final String KUBERNETES_PACKAGE_PREFIX = "io.fabric8.kubernetes.api.model.";
    private static final String KUBERNETES_EXTENSIONS_PACKAGE_PREFIX = "io.fabric8.kubernetes.api.model.extensions.";
    private static final String KUBERNETES_APIEXTENSIONS_PACKAGE_PREFIX = "io.fabric8.kubernetes.api.model.apiextensions.";
    private static final String OPENSHIFT_PACKAGE_PREFIX = "io.fabric8.openshift.api.model.";

    private static final Map<String, Class<? extends KubernetesResource>> MAP = new HashMap<>();


    static {
        // Exceptions (not just package prefix + class name) can be added here.
        MAP.put("List", KubernetesList.class);
    }

    @Override
    public KubernetesResource deserialize(JsonParser jp, DeserializationContext ctxt) throws IOException, JsonProcessingException {
        ObjectNode node = jp.readValueAsTree();
        JsonNode kind = node.get(KIND);
        if (kind != null) {
            String value = kind.textValue();
            Class<? extends KubernetesResource> resourceType = getTypeForName(value);
            if (resourceType == null) {
                throw ctxt.mappingException("No resource type found for kind:" + value);
            } else {
                return jp.getCodec().treeToValue(node, resourceType);
            }
        }
        return null;
    }

    /**
     * Registers a Custom Resource Definition Kind
     */
    public static void registerCustomKind(String kind, Class<? extends KubernetesResource> clazz) {
        MAP.put(kind, clazz);
    }

    private static Class getTypeForName(String name) {
        Class result = MAP.get(name);
        if (result == null) {
            result = loadClassIfExists(KUBERNETES_PACKAGE_PREFIX + name);
            if (result == null) {
                result = loadClassIfExists(KUBERNETES_EXTENSIONS_PACKAGE_PREFIX + name);
                if (result == null) {
                    result = loadClassIfExists(OPENSHIFT_PACKAGE_PREFIX + name);
                    if (result == null) {
                        result = loadClassIfExists(KUBERNETES_APIEXTENSIONS_PACKAGE_PREFIX + name);
                    }
                }
            }
        }

        if (result != null) {
            MAP.put(name, result);
        }
        return result;
    }

    private static Class loadClassIfExists(String className) {
        try {
            return KubernetesDeserializer.class.getClassLoader().loadClass(className);
        } catch (Throwable t) {
            return null;
        }
    }
}
