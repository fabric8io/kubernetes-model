package io.fabric8.kubernetes.client.dsl.internal;

import io.fabric8.kubernetes.api.model.RoleBinding;
import io.fabric8.kubernetes.api.model.RoleBindingList;
import io.fabric8.kubernetes.api.model.DoneableRoleBinding;
import io.fabric8.kubernetes.client.Config;
import io.fabric8.kubernetes.client.dsl.Resource;
import io.fabric8.kubernetes.client.dsl.base.HasMetadataOperation;
import okhttp3.OkHttpClient;

import java.util.Map;
import java.util.TreeMap;

public class RoleBindingOperationsImpl extends HasMetadataOperation<RoleBinding, RoleBindingList, DoneableRoleBinding, Resource<RoleBinding, DoneableRoleBinding>> {
  public RoleBindingOperationsImpl(OkHttpClient client, Config config, String namespace) {
    this(client, config, "v1", namespace, null, true, null, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>());
  }

  public RoleBindingOperationsImpl(OkHttpClient client, Config config, String apiVersion, String namespace, String name, Boolean cascading, RoleBinding item, String resourceVersion, Boolean reloadingFromServer, long gracePeriodSeconds, Map<String, String> labels, Map<String, String> labelsNot, Map<String, String[]> labelsIn, Map<String, String[]> labelsNotIn, Map<String, String> fields) {
    super(client, config, "rbac.authorization.k8s.io", "v1beta1", "rolebindings", namespace, name, cascading, item, resourceVersion, reloadingFromServer, gracePeriodSeconds, labels, labelsNot, labelsIn, labelsNotIn, fields);
  }
}
