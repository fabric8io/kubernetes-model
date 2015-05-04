package io.fabric8.kubernetes.api.model;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.fabric8.kubernetes.api.model.resource.Quantity;
import org.junit.Test;

import static org.junit.Assert.assertEquals;

public class UnmarshallTest {

    @Test
    public void testUnmarshallInt64ToLong() throws Exception {
        ObjectMapper mapper = new ObjectMapper(); // can reuse, share globally
        Pod pod = mapper.readValue(getClass().getResourceAsStream("/valid-pod.json"), Pod.class);
        assertEquals(pod.getSpec().getContainers().get(0).getResources().getLimits().get("memory"), new Quantity("5Mi"));
        assertEquals(pod.getSpec().getContainers().get(0).getResources().getLimits().get("cpu"), new Quantity("1"));
    }
}