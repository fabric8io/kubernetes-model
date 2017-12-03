package io.fabric8.kubernetes.client.handlers;

import io.fabric8.kubernetes.api.model.ClusterRole;
import io.fabric8.kubernetes.api.model.ClusterRoleBuilder;
import io.fabric8.kubernetes.client.Config;
import io.fabric8.kubernetes.client.ResourceHandler;
import io.fabric8.kubernetes.client.Watch;
import io.fabric8.kubernetes.client.Watcher;
import io.fabric8.kubernetes.client.dsl.internal.ClusterRoleOperationsImpl;
import okhttp3.OkHttpClient;

import java.util.TreeMap;
import java.util.concurrent.TimeUnit;

public class ClusterRoleHandler implements ResourceHandler<ClusterRole, ClusterRoleBuilder> {
    @Override
    public String getKind() {
      return ClusterRole.class.getSimpleName();
    }
  
    @Override
    public ClusterRole create(OkHttpClient client, Config config, String namespace, ClusterRole item) {
        return new ClusterRoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).create();
    }
  
    @Override
    public ClusterRole replace(OkHttpClient client, Config config, String namespace, ClusterRole item) {
        return new ClusterRoleOperationsImpl(client, config, null, namespace, null, true, item, null, true, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).replace(item);
    }
  
    @Override
    public ClusterRole reload(OkHttpClient client, Config config, String namespace, ClusterRole item) {
        return new ClusterRoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).fromServer().get();
    }
  
    @Override
    public ClusterRoleBuilder edit(ClusterRole item) {
      return new ClusterRoleBuilder(item);
    }
  
    @Override
    public Boolean delete(OkHttpClient client, Config config, String namespace, ClusterRole item) {
        return new ClusterRoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).delete(item);
    }
  
    @Override
    public Watch watch(OkHttpClient client, Config config, String namespace, ClusterRole item, Watcher<ClusterRole> watcher) {
        return new ClusterRoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).watch(watcher);
    }
  
    @Override
    public Watch watch(OkHttpClient client, Config config, String namespace, ClusterRole item, String resourceVersion, Watcher<ClusterRole> watcher) {
        return new ClusterRoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).watch(resourceVersion, watcher);
    }
  
    @Override
    public ClusterRole waitUntilReady(OkHttpClient client, Config config, String namespace, ClusterRole item, long amount, TimeUnit timeUnit) throws InterruptedException {
        return new ClusterRoleOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).waitUntilReady(amount, timeUnit);
    }
}
