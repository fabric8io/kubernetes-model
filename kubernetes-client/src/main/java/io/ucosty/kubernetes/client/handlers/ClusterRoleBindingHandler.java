package io.ucosty.kubernetes.client.handlers;

import io.ucosty.kubernetes.api.model.ClusterRoleBinding;
import io.ucosty.kubernetes.api.model.ClusterRoleBindingBuilder;
import io.ucosty.kubernetes.client.Config;
import io.ucosty.kubernetes.client.ResourceHandler;
import io.ucosty.kubernetes.client.Watch;
import io.ucosty.kubernetes.client.Watcher;
import io.ucosty.kubernetes.client.dsl.internal.ClusterRoleBindingOperationsImpl;
import okhttp3.OkHttpClient;

import java.util.TreeMap;
import java.util.concurrent.TimeUnit;

public class ClusterRoleBindingHandler implements ResourceHandler<ClusterRoleBinding, ClusterRoleBindingBuilder> {
  @Override
  public String getKind() {
    return ClusterRoleBinding.class.getSimpleName();
  }

  @Override
  public ClusterRoleBinding create(OkHttpClient client, Config config, String namespace, ClusterRoleBinding item) {
    return new ClusterRoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).create();
  }

  @Override
  public ClusterRoleBinding replace(OkHttpClient client, Config config, String namespace, ClusterRoleBinding item) {
    return new ClusterRoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, true, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).replace(item);
  }

  @Override
  public ClusterRoleBinding reload(OkHttpClient client, Config config, String namespace, ClusterRoleBinding item) {
    return new ClusterRoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).fromServer().get();
  }

  @Override
  public ClusterRoleBindingBuilder edit(ClusterRoleBinding item) {
    return new ClusterRoleBindingBuilder(item);
  }

  @Override
  public Boolean delete(OkHttpClient client, Config config, String namespace, ClusterRoleBinding item) {
    return new ClusterRoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).delete(item);
  }

  @Override
  public Watch watch(OkHttpClient client, Config config, String namespace, ClusterRoleBinding item, Watcher<ClusterRoleBinding> watcher) {
    return new ClusterRoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).watch(watcher);
  }

  @Override
  public Watch watch(OkHttpClient client, Config config, String namespace, ClusterRoleBinding item, String resourceVersion, Watcher<ClusterRoleBinding> watcher) {
    return new ClusterRoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).watch(resourceVersion, watcher);
  }

  @Override
  public ClusterRoleBinding waitUntilReady(OkHttpClient client, Config config, String namespace, ClusterRoleBinding item, long amount, TimeUnit timeUnit) throws InterruptedException {
    return new ClusterRoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).waitUntilReady(amount, timeUnit);
  }
}
