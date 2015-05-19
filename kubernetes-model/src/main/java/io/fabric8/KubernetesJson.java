package io.fabric8;

import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import io.fabric8.kubernetes.api.model.KubernetesList;
import io.fabric8.openshift.api.model.template.Template;

@JsonTypeInfo(
        use = JsonTypeInfo.Id.NAME,
        include = JsonTypeInfo.As.PROPERTY,
        property = "kind")
@JsonSubTypes({
        @JsonSubTypes.Type(value = KubernetesList.class, name = "List"),
        @JsonSubTypes.Type(value = Template.class, name = "Template")
})
public interface KubernetesJson {

}
