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
import io.fabric8.kubernetes.api.model.extensions.Deployment;
import io.fabric8.openshift.api.model.DeploymentConfig;
import org.junit.Test;

import static net.javacrumbs.jsonunit.core.Option.IGNORING_ARRAY_ORDER;
import static net.javacrumbs.jsonunit.core.Option.IGNORING_EXTRA_FIELDS;
import static net.javacrumbs.jsonunit.core.Option.TREATING_NULL_AS_ABSENT;
import static net.javacrumbs.jsonunit.fluent.JsonFluentAssert.assertThatJson;

public class DeploymentTest {

    private final ObjectMapper mapper = new ObjectMapper();

    @Test
    public void kubernetesDeploymentTest() throws Exception {
        // given
        final String originalJson = Helper.loadJson("/valid-deployment.json");

        // when
        final Deployment deployment = mapper.readValue(originalJson, Deployment.class);
        final String serializedJson = mapper.writeValueAsString(deployment);

        // then
        assertThatJson(serializedJson).when(IGNORING_ARRAY_ORDER, TREATING_NULL_AS_ABSENT, IGNORING_EXTRA_FIELDS)
                .isEqualTo(originalJson);
    }

    @Test
    public void openshiftDeploymentConfigTest() throws Exception {
        // given
        final String originalJson = Helper.loadJson("/valid-deploymentConfig.json");

        // when
        final DeploymentConfig deploymentConfig = mapper.readValue(originalJson, DeploymentConfig.class);
        final String serializedJson = mapper.writeValueAsString(deploymentConfig);

        // then
        assertThatJson(serializedJson).when(IGNORING_ARRAY_ORDER, TREATING_NULL_AS_ABSENT, IGNORING_EXTRA_FIELDS)
                .isEqualTo(originalJson);
    }

}
