package io.fabric8.kubernetes.api.model;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.Test;

import static net.javacrumbs.jsonunit.core.Option.IGNORING_ARRAY_ORDER;
import static net.javacrumbs.jsonunit.core.Option.IGNORING_EXTRA_FIELDS;
import static net.javacrumbs.jsonunit.core.Option.TREATING_NULL_AS_ABSENT;
import static net.javacrumbs.jsonunit.fluent.JsonFluentAssert.assertThatJson;

public class SecretTest {
    private final ObjectMapper mapper = new ObjectMapper();

    @Test
    public void SecretTest() throws Exception {
        // given
        final String originalJson = Helper.loadJson("/valid-secret.json");

        // when
        final Secret secret = mapper.readValue(originalJson, Secret.class);
        final String serializedJson = mapper.writeValueAsString(secret);

        // then
        assertThatJson(serializedJson).when(IGNORING_ARRAY_ORDER, TREATING_NULL_AS_ABSENT, IGNORING_EXTRA_FIELDS)
                .isEqualTo(originalJson);
    }
}
