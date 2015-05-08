package io.fabric8.kubernetes.api.model;

import com.fasterxml.jackson.core.JsonProcessingException;
import org.hamcrest.CoreMatchers;
import org.junit.Test;

import static org.junit.Assert.*;

public class KubernetesListTest {

    @Test
    public void testDefaultValues() throws JsonProcessingException {
        Service service = new ServiceBuilder()
                .withNewMetadata()
                    .withName("test-service")
                .endMetadata()
                .build();
        assertNotNull(service.getApiVersion());
        assertEquals(service.getKind(), "Service");
        
        ReplicationController replicationController = new ReplicationControllerBuilder()
                .withNewMetadata()
                .withName("test-controller")
                .endMetadata()
                .build();
        assertNotNull(replicationController.getApiVersion());
        assertEquals(replicationController.getKind(), "ReplicationController");
        
        KubernetesList kubernetesList = new KubernetesListBuilder()
                .addNewService()
                .withNewMetadata()
                    .withName("test-service")
                .endMetadata()
                .and()
                .addNewReplicationController()
                .withNewMetadata()
                    .withName("test-controller")
                .endMetadata()
                .and()
                .build();
        
        assertNotNull(kubernetesList.getApiVersion());
        assertEquals(kubernetesList.getKind(), "List");
        assertThat(kubernetesList.getItems(), CoreMatchers.hasItem(service));
        assertThat(kubernetesList.getItems(), CoreMatchers.hasItem(replicationController));
    }

    @Test
    public void testVisitor() throws JsonProcessingException {
        KubernetesList list = new KubernetesListBuilder()
                .addNewPod()
                    .withNewSpec()
                        .addNewContainer()
                            .withName("my-container")
                            .withImage("my/image")
                        .endContainer()
                    .endSpec()
                .endPod()
                .build();

        list = new KubernetesListBuilder(list).accept(new io.fabric8.common.Visitor() {
            public void visit(Object item) {
                if (item instanceof io.fabric8.kubernetes.api.model.PodSpecBuilder) {
                    ((io.fabric8.kubernetes.api.model.PodSpecBuilder)item).addNewContainer()
                            .withName("other-container")
                            .withImage("other/image")
                            .endContainer();
                }
            }
        }).build();
    }
}