package io.fabric8.kubernetes.api.model;

import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.IOException;
import java.io.InputStream;
import java.util.Scanner;
import org.junit.Test;

import static net.javacrumbs.jsonunit.core.Option.IGNORING_ARRAY_ORDER;
import static net.javacrumbs.jsonunit.core.Option.IGNORING_EXTRA_FIELDS;
import static net.javacrumbs.jsonunit.core.Option.TREATING_NULL_AS_ABSENT;
import static net.javacrumbs.jsonunit.fluent.JsonFluentAssert.assertThatJson;

public class JsonConsistencyTest {

    private final ObjectMapper mapper = new ObjectMapper();

    @Test
    public void should_produce_same_json_from_unmarshalled_one() throws Exception {
        // given
        final String originalPodJson = loadJson("/valid-pod.json");

        // when
        final Pod pod = mapper.readValue(originalPodJson, Pod.class);
        final String serializedPodAsJson = mapper.writeValueAsString(pod);

        // then
        assertThatJson(serializedPodAsJson).when(IGNORING_ARRAY_ORDER, TREATING_NULL_AS_ABSENT, IGNORING_EXTRA_FIELDS)
            .isEqualTo(originalPodJson);
    }

    private String loadJson(String path) {
        try (InputStream resourceAsStream = getClass().getResourceAsStream(path)) {
            final Scanner scanner = new Scanner(resourceAsStream).useDelimiter("\\A");
            return scanner.hasNext() ? scanner.next() : "";
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}
