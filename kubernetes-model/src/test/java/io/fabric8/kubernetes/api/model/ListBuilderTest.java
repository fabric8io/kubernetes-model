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



import com.fasterxml.jackson.databind.ObjectMapper;

import org.junit.Assert;
import org.junit.Test;

import java.io.IOException;

import io.fabric8.kubernetes.api.model.extensions.Deployment;
import io.fabric8.kubernetes.api.model.extensions.DeploymentList;
import io.fabric8.kubernetes.api.model.extensions.DeploymentListBuilder;
import io.fabric8.kubernetes.api.model.extensions.ReplicaSet;
import io.fabric8.kubernetes.api.model.extensions.ReplicaSetList;
import io.fabric8.kubernetes.api.model.extensions.ReplicaSetListBuilder;

public class ListBuilderTest {

    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    @Test
    public void podListShouldHaveNotNullItems() throws IOException {
        PodList list = new PodListBuilder().build();
        Assert.assertNotNull(list.getItems());

        list = new PodListBuilder().withItems((Pod[]) null).build();
        Assert.assertNotNull(list.getItems());

        list = OBJECT_MAPPER.readValue(getClass().getResourceAsStream("/null-pod-list.json"), PodList.class);
        Assert.assertNotNull(list.getItems());
    }

    @Test
    public void serviceListShouldHaveNotNullItems() throws IOException {
        ServiceList list = new ServiceListBuilder().build();
        Assert.assertNotNull(list.getItems());

        list = new ServiceListBuilder().withItems((Service[]) null).build();
        Assert.assertNotNull(list.getItems());

        list = OBJECT_MAPPER.readValue(getClass().getResourceAsStream("/null-service-list.json"), ServiceList.class);
        Assert.assertNotNull(list.getItems());

    }

    @Test
    public void replicationControllerListShouldHaveNotNullItems() throws IOException {
        ReplicationControllerList list = new ReplicationControllerListBuilder().build();
        Assert.assertNotNull(list.getItems());

        list = new ReplicationControllerListBuilder().withItems((ReplicationController[]) null).build();
        Assert.assertNotNull(list.getItems());

        list = OBJECT_MAPPER.readValue(getClass().getResourceAsStream("/null-rc-list.json"), ReplicationControllerList.class);
        Assert.assertNotNull(list.getItems());
    }



    @Test
    public void replicaSetListShouldHaveNotNullItems() throws IOException {
        ReplicaSetList list = new ReplicaSetListBuilder().build();
        Assert.assertNotNull(list.getItems());

        list = new ReplicaSetListBuilder().withItems((ReplicaSet[]) null).build();
        Assert.assertNotNull(list.getItems());

        list = OBJECT_MAPPER.readValue(getClass().getResourceAsStream("/null-rs-list.json"), ReplicaSetList.class);
        Assert.assertNotNull(list.getItems());
    }

    @Test
    public void deploymentListShouldHaveNotNullItems() throws IOException {
        DeploymentList list = new DeploymentListBuilder().build();
        Assert.assertNotNull(list.getItems());

        list = new DeploymentListBuilder().withItems((Deployment[]) null).build();
        Assert.assertNotNull(list.getItems());

        list = OBJECT_MAPPER.readValue(getClass().getResourceAsStream("/null-deployment-list.json"), DeploymentList.class);
        Assert.assertNotNull(list.getItems());
    }
}
