package io.fabric8.openshift.api.model;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;
import io.sundr.builder.annotations.Buildable;

import javax.annotation.Generated;

@JsonInclude(JsonInclude.Include.NON_NULL)
@Generated("org.jsonschema2pojo")
@JsonPropertyOrder({
        "created",
        "dockerImageReference",
        "image",
})
public class TagEvent {

    @JsonProperty("created")
    private String created;
    
    @JsonProperty("dockerImageReference")
    private java.lang.String dockerImageReference;

    @JsonProperty("image")
    private Image image;

    public TagEvent() {
    }
    
    @Buildable
    public TagEvent(String created, String dockerImageReference, Image image) {
        this.created = created;
        this.dockerImageReference = dockerImageReference;
        this.image = image;
    }

    @JsonProperty("created")
    public String getCreated() {
        return created;
    }

    @JsonProperty("created")
    public void setCreated(String created) {
        this.created = created;
    }

    @JsonProperty("dockerImageReference")
    public String getDockerImageReference() {
        return dockerImageReference;
    }

    @JsonProperty("dockerImageReference")
    public void setDockerImageReference(String dockerImageReference) {
        this.dockerImageReference = dockerImageReference;
    }

    @JsonProperty("image")
    public Image getImage() {
        return image;
    }

    @JsonProperty("image")
    public void setImage(Image image) {
        this.image = image;
    }
}
