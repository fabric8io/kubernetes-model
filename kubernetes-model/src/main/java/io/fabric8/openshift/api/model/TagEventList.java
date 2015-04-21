package io.fabric8.openshift.api.model;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;
import io.sundr.builder.annotations.Buildable;

import javax.annotation.Generated;
import javax.validation.Valid;
import java.util.ArrayList;
import java.util.List;

/**
 *
 *
 */
@JsonInclude(JsonInclude.Include.NON_NULL)
@Generated("org.jsonschema2pojo")
@JsonPropertyOrder({
        "items"
})
public class TagEventList {

    /**
     * list of objects
     *
     */
    @JsonProperty("items")
    @Valid
    private List<TagEvent> items = new ArrayList<TagEvent>();

    public TagEventList(){
        
    }

    @Buildable
    public TagEventList(List<TagEvent> items) {
        this.items = items;
    }

    @JsonProperty("items")
    public List<TagEvent> getItems() {
        return items;
    }

    @JsonProperty("items")
    public void setItems(List<TagEvent> items) {
        this.items = items;
    }
}
