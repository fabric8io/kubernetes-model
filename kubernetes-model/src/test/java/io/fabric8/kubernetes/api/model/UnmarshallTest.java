package io.fabric8.kubernetes.api.model;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.Test;

import static org.junit.Assert.assertEquals;

public class UnmarshallTest {

    @Test
    public void testUnmarshallInt64ToLong() throws Exception {
        ObjectMapper mapper = new ObjectMapper(); // can reuse, share globally
        ReplicationController rc = mapper.readValue(getClass().getResourceAsStream("/meteor-controller.json"), ReplicationController.class);
        assertEquals(rc.getDesiredState().getPodTemplate().getDesiredState().getManifest().getContainers().get(0).getMemory(), 500000000l);
    }
}