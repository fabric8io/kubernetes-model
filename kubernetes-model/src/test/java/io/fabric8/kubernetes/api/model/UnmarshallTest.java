package io.fabric8.kubernetes.api.model;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.fabric8.common.Visitor;
import io.fabric8.kubernetes.api.model.resource.Quantity;
import org.junit.Assert;
import org.junit.Test;

import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.Assert.assertEquals;

public class UnmarshallTest {

    @Test
    public void testUnmarshallInt64ToLong() throws Exception {
        ObjectMapper mapper = new ObjectMapper(); // can reuse, share globally
        Pod pod = mapper.readValue(getClass().getResourceAsStream("/valid-pod.json"), Pod.class);
        assertEquals(pod.getSpec().getContainers().get(0).getResources().getLimits().get("memory"), new Quantity("5Mi"));
        assertEquals(pod.getSpec().getContainers().get(0).getResources().getLimits().get("cpu"), new Quantity("1"));
    }

    @Test
    public void testUnmarshallWithVisitors() throws Exception {
        ObjectMapper mapper = new ObjectMapper(); // can reuse, share globally
        KubernetesList list = mapper.readValue(getClass().getResourceAsStream("/simple-list.json"), KubernetesList.class);
        final AtomicInteger integer = new AtomicInteger();
        new KubernetesListBuilder(list).accept(new Visitor() {
            public void visit(Object o) {
                integer.incrementAndGet();
            }
        });

        //We just want to make sure that it visits nested objects when deserialization from json is used.
        // The exact number is volatile so we just care about the minimum number of objects (list, pod and service).
        Assert.assertTrue(integer.intValue() >= 3);

    }
}