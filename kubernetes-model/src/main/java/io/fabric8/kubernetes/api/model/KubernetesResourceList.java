package io.fabric8.kubernetes.api.model;

import java.util.List;

import com.fasterxml.jackson.databind.annotation.JsonDeserialize;
import io.fabric8.kubernetes.internal.KubernetesDeserializer;

@JsonDeserialize(using = KubernetesDeserializer.class)
public interface KubernetesResourceList<E extends io.fabric8.kubernetes.api.model.HasMetadata> {

  List<E> getItems();

}
