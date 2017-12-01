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
package io.fabric8.kubernetes.api.model;

import java.util.HashMap;
import java.util.Map;


@Deprecated
public enum KubernetesKind {

    List(KubernetesList.class),
    ServiceAccount(ServiceAccount.class),
    ServiceAccountList(ServiceAccountList.class),
    Service(Service.class),
    ServiceList(ServiceList.class),
    Pod(Pod.class),
    PodList(PodList.class),
    Namespace(Namespace.class),
    NamespaceList(NamespaceList.class),
    Secret(Secret.class),
    SecretList(SecretList.class),
    Endpoints(Endpoints.class),
    EndpointsList(EndpointsList.class),
    Node(Node.class),
    NodeList(NodeList.class),
    Role(Role.class),
    RoleList(RoleList.class),
    RoleBinding(RoleBinding.class),
    RoleBindingList(RoleBindingList.class),
    ClusterRoleBinding(ClusterRoleBinding.class),
    ClusterRoleBindingList(ClusterRoleBindingList.class),
    PersistentVolume(PersistentVolume.class),
    PersistentVolumeList(PersistentVolumeList.class),
    PersistentVolumeClaim(PersistentVolumeClaim.class),
    PersistentVolumeClaimList(PersistentVolumeClaimList.class);

    private static final Map<String, Class<? extends KubernetesResource>> map = new HashMap<>();

    static {
        for (KubernetesKind kind : KubernetesKind.values()) {
            map.put(kind.name(), kind.type);
        }
    }

    private final Class<? extends KubernetesResource> type;

    KubernetesKind(Class type) {
        this.type = type;
    }

    public Class getType() {
        return type;
    }

    public static Class<? extends KubernetesResource> getTypeForName(String name) {
        return map.get(name);
    }
}
