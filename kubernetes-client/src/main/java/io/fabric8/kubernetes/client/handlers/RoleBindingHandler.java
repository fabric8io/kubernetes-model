package io.fabric8.kubernetes.client.handlers;

import io.fabric8.kubernetes.api.model.RoleBinding;
import io.fabric8.kubernetes.api.model.RoleBindingBuilder;
import io.fabric8.kubernetes.client.Config;
import io.fabric8.kubernetes.client.ResourceHandler;
import io.fabric8.kubernetes.client.Watch;
import io.fabric8.kubernetes.client.Watcher;
import io.fabric8.kubernetes.client.dsl.internal.RoleBindingOperationsImpl;
import okhttp3.OkHttpClient;

import java.util.TreeMap;
import java.util.concurrent.TimeUnit;

public class RoleBindingHandler implements ResourceHandler<RoleBinding, RoleBindingBuilder> {
  @Override
  public String getKind() {
    return RoleBinding.class.getSimpleName();
  }

  @Override
  public RoleBinding create(OkHttpClient client, Config config, String namespace, RoleBinding item) {
    return new RoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).create();
  }

  @Override
  public RoleBinding replace(OkHttpClient client, Config config, String namespace, RoleBinding item) {
    return new RoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, true, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).replace(item);
  }

  @Override
  public RoleBinding reload(OkHttpClient client, Config config, String namespace, RoleBinding item) {
    return new RoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).fromServer().get();
  }

  @Override
  public RoleBindingBuilder edit(RoleBinding item) {
    return new RoleBindingBuilder(item);
  }

  @Override
  public Boolean delete(OkHttpClient client, Config config, String namespace, RoleBinding item) {
    return new RoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).delete(item);
  }

  @Override
  public Watch watch(OkHttpClient client, Config config, String namespace, RoleBinding item, Watcher<RoleBinding> watcher) {
    return new RoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).watch(watcher);
  }

  @Override
  public Watch watch(OkHttpClient client, Config config, String namespace, RoleBinding item, String resourceVersion, Watcher<RoleBinding> watcher) {
    return new RoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).watch(resourceVersion, watcher);
  }

  @Override
  public RoleBinding waitUntilReady(OkHttpClient client, Config config, String namespace, RoleBinding item, long amount, TimeUnit timeUnit) throws InterruptedException {
    return new RoleBindingOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).waitUntilReady(amount, timeUnit);
  }
}

