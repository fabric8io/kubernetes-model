/**
 * Copyright (C) 2015 Red Hat, Inc.
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
package io.ucosty.kubernetes.client.handlers;

import io.ucosty.kubernetes.client.Watch;
import io.ucosty.kubernetes.client.Watcher;
import okhttp3.OkHttpClient;
import io.ucosty.kubernetes.api.model.PersistentVolumeClaim;
import io.ucosty.kubernetes.api.model.PersistentVolumeClaimBuilder;
import io.ucosty.kubernetes.client.Config;
import io.ucosty.kubernetes.client.ResourceHandler;
import io.ucosty.kubernetes.client.dsl.internal.PersistentVolumeClaimOperationsImpl;
import org.apache.felix.scr.annotations.Component;
import org.apache.felix.scr.annotations.Service;

import java.util.TreeMap;
import java.util.concurrent.TimeUnit;

@Component
@Service
public class PersistentVolumeClaimHandler implements ResourceHandler<PersistentVolumeClaim, PersistentVolumeClaimBuilder> {
  @Override
  public String getKind() {
    return PersistentVolumeClaim.class.getSimpleName();
  }

  @Override
  public PersistentVolumeClaim create(OkHttpClient client, Config config, String namespace, PersistentVolumeClaim item) {
    return new PersistentVolumeClaimOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).create();
  }

  @Override
  public PersistentVolumeClaim replace(OkHttpClient client, Config config, String namespace, PersistentVolumeClaim item) {
    return new PersistentVolumeClaimOperationsImpl(client, config, null, namespace, null, true, item, null, true, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).replace(item);
  }

  @Override
  public PersistentVolumeClaim reload(OkHttpClient client, Config config, String namespace, PersistentVolumeClaim item) {
    return new PersistentVolumeClaimOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).fromServer().get();
  }

  @Override
  public PersistentVolumeClaimBuilder edit(PersistentVolumeClaim item) {
    return new PersistentVolumeClaimBuilder(item);
  }

  @Override
  public Boolean delete(OkHttpClient client, Config config, String namespace, PersistentVolumeClaim item) {
    return new PersistentVolumeClaimOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).delete(item);
  }

  @Override
  public Watch watch(OkHttpClient client, Config config, String namespace, PersistentVolumeClaim item, Watcher<PersistentVolumeClaim> watcher) {
    return new PersistentVolumeClaimOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).watch(watcher);
  }

  @Override
  public Watch watch(OkHttpClient client, Config config, String namespace, PersistentVolumeClaim item, String resourceVersion, Watcher<PersistentVolumeClaim> watcher) {
    return new PersistentVolumeClaimOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).watch(resourceVersion, watcher);
  }

  @Override
  public PersistentVolumeClaim waitUntilReady(OkHttpClient client, Config config, String namespace, PersistentVolumeClaim item, long amount, TimeUnit timeUnit) throws InterruptedException {
    return new PersistentVolumeClaimOperationsImpl(client, config, null, namespace, null, true, item, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>()).waitUntilReady(amount, timeUnit);
  }
}
