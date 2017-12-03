package io.ucosty.kubernetes.client.handlers;

import io.ucosty.kubernetes.api.model.Role;
import io.ucosty.kubernetes.api.model.RoleBuilder;
import io.ucosty.kubernetes.client.Config;
import io.ucosty.kubernetes.client.ResourceHandler;
import io.ucosty.kubernetes.client.Watch;
import io.ucosty.kubernetes.client.Watcher;
import io.ucosty.kubernetes.client.dsl.internal.RoleOperationsImpl;
import okhttp3.OkHttpClient;

import java.util.TreeMap;
import java.util.concurrent.TimeUnit;

public class RoleHandler implements ResourceHandler<Role, RoleBuilder> {
  @Override
  public String getKind() {
    return Role.class.getSimpleName();
  }

  @Override
  public Role create(OkHttpClient client, Config config, String namespace, Role item) {
    return new RoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).create();
  }

  @Override
  public Role replace(OkHttpClient client, Config config, String namespace, Role item) {
    return new RoleOperationsImpl(client, config, null, namespace, null, true, item, null, true, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).replace(item);
  }

  @Override
  public Role reload(OkHttpClient client, Config config, String namespace, Role item) {
    return new RoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).fromServer().get();
  }

  @Override
  public RoleBuilder edit(Role item) {
    return new RoleBuilder(item);
  }

  @Override
  public Boolean delete(OkHttpClient client, Config config, String namespace, Role item) {
    return new RoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).delete(item);
  }

  @Override
  public Watch watch(OkHttpClient client, Config config, String namespace, Role item, Watcher<Role> watcher) {
    return new RoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).watch(watcher);
  }

  @Override
  public Watch watch(OkHttpClient client, Config config, String namespace, Role item, String resourceVersion, Watcher<Role> watcher) {
    return new RoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).watch(resourceVersion, watcher);
  }

  @Override
  public Role waitUntilReady(OkHttpClient client, Config config, String namespace, Role item, long amount, TimeUnit timeUnit) throws InterruptedException {
    return new RoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).waitUntilReady(amount, timeUnit);
  }
}
