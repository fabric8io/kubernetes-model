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
    public void KubernetesDeploymentTest() throws Exception {
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
    public void OpenshiftDeploymentConfigTest() throws Exception {
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
