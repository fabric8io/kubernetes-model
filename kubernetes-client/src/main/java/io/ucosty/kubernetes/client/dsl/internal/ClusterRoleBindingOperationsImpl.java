package io.ucosty.kubernetes.client.dsl.internal;

import io.ucosty.kubernetes.api.model.ClusterRoleBinding;
import io.ucosty.kubernetes.api.model.ClusterRoleBindingList;
import io.ucosty.kubernetes.api.model.DoneableClusterRoleBinding;
import io.ucosty.kubernetes.client.Config;
import io.ucosty.kubernetes.client.dsl.Resource;
import io.ucosty.kubernetes.client.dsl.base.HasMetadataOperation;
import okhttp3.OkHttpClient;

import java.util.Map;
import java.util.TreeMap;

public class ClusterRoleBindingOperationsImpl extends HasMetadataOperation<ClusterRoleBinding, ClusterRoleBindingList, DoneableClusterRoleBinding, Resource<ClusterRoleBinding, DoneableClusterRoleBinding>> {
  public ClusterRoleBindingOperationsImpl(OkHttpClient client, Config config, String namespace) {
    this(client, config, "v1", namespace, null, true, null, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>());
  }

  public ClusterRoleBindingOperationsImpl(OkHttpClient client, Config config, String apiVersion, String namespace, String name, Boolean cascading, ClusterRoleBinding item, String resourceVersion, Boolean reloadingFromServer, long gracePeriodSeconds, Map<String, String> labels, Map<String, String> labelsNot, Map<String, String[]> labelsIn, Map<String, String[]> labelsNotIn, Map<String, String> fields) {
    super(client, config, "rbac.authorization.k8s.io", "v1beta1", "clusterrolebindings", namespace, name, cascading, item, resourceVersion, reloadingFromServer, gracePeriodSeconds, labels, labelsNot, labelsIn, labelsNotIn, fields);
  }

  @Override
  public boolean isResourceNamespaced() {
    return false;
  }
}
