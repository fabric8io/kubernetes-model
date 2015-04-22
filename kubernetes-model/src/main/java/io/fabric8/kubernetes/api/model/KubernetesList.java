
package io.fabric8.kubernetes.api.model;

import com.fasterxml.jackson.annotation.JsonAnyGetter;
import com.fasterxml.jackson.annotation.JsonAnySetter;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;
import io.fabric8.openshift.api.model.BuildConfig;
import io.fabric8.openshift.api.model.DeploymentConfig;
import io.fabric8.openshift.api.model.ImageRepository;
import io.fabric8.openshift.api.model.Route;
import org.apache.commons.lang.builder.EqualsBuilder;
import org.apache.commons.lang.builder.HashCodeBuilder;
import org.apache.commons.lang.builder.ToStringBuilder;

import javax.annotation.Generated;
import javax.validation.Valid;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;


/**
 * 
 * 
 */
@JsonInclude(JsonInclude.Include.NON_NULL)
@Generated("org.jsonschema2pojo")
@JsonPropertyOrder({
    "annotations",
    "apiVersion",
    "creationTimestamp",
    "id",
    "items",
    "kind",
    "namespace",
    "resourceVersion",
    "selfLink",
    "uid"
})
public class KubernetesList {

    /**
     * map of string keys and values that can be used by external tooling to store and retrieve arbitrary metadata about the object
     * 
     */
    @JsonProperty("annotations")
    @Valid
    private Map<String, String> annotations;
    /**
     * version of the schema the object should have
     * 
     */
    @JsonProperty("apiVersion")
    private java.lang.String apiVersion = "v1beta2";
    /**
     * RFC 3339 date and time at which the object was created; recorded by the system; null for lists
     *
     */
    @JsonProperty("creationTimestamp")
    private String creationTimestamp;
    /**
     * name of the object; must be a DNS_SUBDOMAIN and unique among all objects of the same kind within the same namespace; used in resource URLs
     *
     */
    @JsonProperty("id")
    private String id;
    
    /**
     * list of objects
     *
     */
    @JsonProperty("items")
    @Valid
    private List<Object> items = new ArrayList<Object>();

    /**
     * list of services
     * Note: This is not to be used. Added for influencing the generation of fluent nested builders
     */
    @JsonIgnore
    private final List<Service> services = Collections.emptyList();

    /**
     * list of replication controllers
     * Note: This is not to be used. Added for influencing the generation of fluent nested builders
     */
    @JsonIgnore
    private final List<ReplicationController> replicationControllers = Collections.emptyList();

    /**
     * list of pods. 
     * Note: This is not to be used. Added for influencing the generation of fluent nested builders 
     *
     */
    @JsonIgnore
    private final List<Pod> pods = Collections.emptyList();
    
    /**
     * list of build configs
     * Note: This is not to be used. Added for influencing the generation of fluent nested builders
     *
     */
    @JsonIgnore
    private final List<BuildConfig> buildConfigs = Collections.emptyList();

    /**
     * list of deployment configs
     * Note: This is not to be used. Added for influencing the generation of fluent nested builders
     */
    @JsonIgnore
    private final List<DeploymentConfig> deploymentConfigs = Collections.emptyList();

    /**
     * list of image repositories
     * Note: This is not to be used. Added for influencing the generation of fluent nested builders
     */
    @JsonIgnore
    private final List<ImageRepository> imageRepositories = Collections.emptyList();

    /**
     * list of routes
     * Note: This is not to be used. Added for influencing the generation of fluent nested builders.
     */
    @JsonIgnore
    private final List<Route> routes = Collections.emptyList();
    
    
    /**
     *
     *
     */
    @JsonProperty("kind")
    private String kind = "List";
    /**
     * namespace to which the object belongs; must be a DNS_SUBDOMAIN; 'default' by default
     *
     */
    @JsonProperty("namespace")
    private String namespace;
    /**
     * string that identifies the internal version of this object that can be used by clients to determine when objects have changed; value must be treated as opaque by clients and passed unmodified back to the server
     *
     */
    @JsonProperty("resourceVersion")
    private Integer resourceVersion;
    /**
     * URL for the object
     *
     */
    @JsonProperty("selfLink")
    private String selfLink;
    /**
     * UUID assigned by the system upon creation
     *
     */
    @JsonProperty("uid")
    private String uid;
    @JsonIgnore
    private Map<String, Object> additionalProperties = new HashMap<String, Object>();

    /**
     * No args constructor for use in serialization
     *
     */
    public KubernetesList() {
    }

    /**
     *
     * @param uid
     * @param id
     * @param apiVersion
     * @param items
     * @param resourceVersion
     * @param selfLink
     * @param creationTimestamp
     * @param annotations
     * @param kind
     * @param namespace
     */
    public KubernetesList(Map<String, String> annotations, String apiVersion, String creationTimestamp, String id,
                          List<Service> services,
                          List<ReplicationController> replicationControllers,
                          List<Pod> pods,
                          List<BuildConfig> buildConfigs,
                          List<DeploymentConfig> deploymentConfigs,
                          List<ImageRepository> imageRepositories,
                          List<Route> routes,
                          List<Object> items,
                          String kind, String namespace, Integer resourceVersion, String selfLink, String uid) {
        this.annotations = annotations;
        this.apiVersion = apiVersion;
        this.creationTimestamp = creationTimestamp;
        this.id = id;
        this.kind = kind;
        this.namespace = namespace;
        this.resourceVersion = resourceVersion;
        this.selfLink = selfLink;
        this.uid = uid;
        Set<Object> allItems = new LinkedHashSet<>();
        allItems.addAll(items != null ? items : Collections.emptyList());
        allItems.addAll(services != null ? services : Collections.<Service>emptyList());
        allItems.addAll(replicationControllers != null ? replicationControllers : Collections.<ReplicationController>emptyList());
        allItems.addAll(pods != null ? pods : Collections.<Pod>emptyList());
        allItems.addAll(buildConfigs != null ? buildConfigs : Collections.<BuildConfig>emptyList());
        allItems.addAll(deploymentConfigs != null ? deploymentConfigs : Collections.<DeploymentConfig>emptyList());
        allItems.addAll(imageRepositories != null ? imageRepositories : Collections.<ImageRepository>emptyList());
        allItems.addAll(routes != null ? routes : Collections.<Route>emptyList());
        this.items = new ArrayList<>(allItems);
    }

    /**
     * map of string keys and values that can be used by external tooling to store and retrieve arbitrary metadata about the object
     *
     * @return
     *     The annotations
     */
    @JsonProperty("annotations")
    public Map<String, String> getAnnotations() {
        return annotations;
    }

    /**
     * map of string keys and values that can be used by external tooling to store and retrieve arbitrary metadata about the object
     *
     * @param annotations
     *     The annotations
     */
    @JsonProperty("annotations")
    public void setAnnotations(Map<String, String> annotations) {
        this.annotations = annotations;
    }

    /**
     * version of the schema the object should have
     *
     * @return
     *     The apiVersion
     */
    @JsonProperty("apiVersion")
    public String getApiVersion() {
        return apiVersion;
    }

    /**
     * version of the schema the object should have
     *
     * @param apiVersion
     *     The apiVersion
     */
    @JsonProperty("apiVersion")
    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    /**
     * RFC 3339 date and time at which the object was created; recorded by the system; null for lists
     *
     * @return
     *     The creationTimestamp
     */
    @JsonProperty("creationTimestamp")
    public String getCreationTimestamp() {
        return creationTimestamp;
    }

    /**
     * RFC 3339 date and time at which the object was created; recorded by the system; null for lists
     *
     * @param creationTimestamp
     *     The creationTimestamp
     */
    @JsonProperty("creationTimestamp")
    public void setCreationTimestamp(String creationTimestamp) {
        this.creationTimestamp = creationTimestamp;
    }

    /**
     * name of the object; must be a DNS_SUBDOMAIN and unique among all objects of the same kind within the same namespace; used in resource URLs
     *
     * @return
     *     The id
     */
    @JsonProperty("id")
    public String getId() {
        return id;
    }

    /**
     * name of the object; must be a DNS_SUBDOMAIN and unique among all objects of the same kind within the same namespace; used in resource URLs
     *
     * @param id
     *     The id
     */
    @JsonProperty("id")
    public void setId(String id) {
        this.id = id;
    }

    /**
     * list of objects
     *
     * @return
     *     The items
     */
    @JsonProperty("items")
    public List<Object> getItems() {
        return items;
    }

    /**
     * list of objects
     *
     * @param items
     *     The items
     */
    @JsonProperty("items")
    public void setItems(List<Object> items) {
        this.items = items;
    }

    @JsonIgnore
    public List<Service> getServices() {
        return services;
    }
    
    @JsonIgnore
    public List<ReplicationController> getReplicationControllers() {
        return replicationControllers;
    }

    @JsonIgnore
    public List<Pod> getPods() {
        return pods;
    }
    
    @JsonIgnore
    public List<BuildConfig> getBuildConfigs() {
        return buildConfigs;
    }
    
    @JsonIgnore
    public List<DeploymentConfig> getDeploymentConfigs() {
        return deploymentConfigs;
    }
    
    @JsonIgnore
    public List<ImageRepository> getImageRepositories() {
        return imageRepositories;
    }
    
    @JsonIgnore
    public List<Route> getRoutes() {
        return routes;
    }

    /**
     *
     *
     * @return
     *     The kind
     */
    @JsonProperty("kind")
    public String getKind() {
        return kind;
    }

    /**
     *
     *
     * @param kind
     *     The kind
     */
    @JsonProperty("kind")
    public void setKind(String kind) {
        this.kind = kind;
    }

    /**
     * namespace to which the object belongs; must be a DNS_SUBDOMAIN; 'default' by default
     *
     * @return
     *     The namespace
     */
    @JsonProperty("namespace")
    public String getNamespace() {
        return namespace;
    }

    /**
     * namespace to which the object belongs; must be a DNS_SUBDOMAIN; 'default' by default
     *
     * @param namespace
     *     The namespace
     */
    @JsonProperty("namespace")
    public void setNamespace(String namespace) {
        this.namespace = namespace;
    }

    /**
     * string that identifies the internal version of this object that can be used by clients to determine when objects have changed; value must be treated as opaque by clients and passed unmodified back to the server
     *
     * @return
     *     The resourceVersion
     */
    @JsonProperty("resourceVersion")
    public Integer getResourceVersion() {
        return resourceVersion;
    }

    /**
     * string that identifies the internal version of this object that can be used by clients to determine when objects have changed; value must be treated as opaque by clients and passed unmodified back to the server
     *
     * @param resourceVersion
     *     The resourceVersion
     */
    @JsonProperty("resourceVersion")
    public void setResourceVersion(Integer resourceVersion) {
        this.resourceVersion = resourceVersion;
    }

    /**
     * URL for the object
     *
     * @return
     *     The selfLink
     */
    @JsonProperty("selfLink")
    public String getSelfLink() {
        return selfLink;
    }

    /**
     * URL for the object
     *
     * @param selfLink
     *     The selfLink
     */
    @JsonProperty("selfLink")
    public void setSelfLink(String selfLink) {
        this.selfLink = selfLink;
    }

    /**
     * UUID assigned by the system upon creation
     *
     * @return
     *     The uid
     */
    @JsonProperty("uid")
    public String getUid() {
        return uid;
    }

    /**
     * UUID assigned by the system upon creation
     *
     * @param uid
     *     The uid
     */
    @JsonProperty("uid")
    public void setUid(String uid) {
        this.uid = uid;
    }

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    @JsonAnyGetter
    public Map<String, Object> getAdditionalProperties() {
        return this.additionalProperties;
    }

    @JsonAnySetter
    public void setAdditionalProperty(String name, Object value) {
        this.additionalProperties.put(name, value);
    }

    @Override
    public int hashCode() {
        return new HashCodeBuilder().append(annotations).append(apiVersion).append(creationTimestamp).append(id).append(items).append(kind).append(namespace).append(resourceVersion).append(selfLink).append(uid).append(additionalProperties).toHashCode();
    }

    @Override
    public boolean equals(Object other) {
        if (other == this) {
            return true;
        }
        if ((other instanceof KubernetesList) == false) {
            return false;
        }
        KubernetesList rhs = ((KubernetesList) other);
        return new EqualsBuilder().append(annotations, rhs.annotations).append(apiVersion, rhs.apiVersion).append(creationTimestamp, rhs.creationTimestamp).append(id, rhs.id).append(items, rhs.items).append(kind, rhs.kind).append(namespace, rhs.namespace).append(resourceVersion, rhs.resourceVersion).append(selfLink, rhs.selfLink).append(uid, rhs.uid).append(additionalProperties, rhs.additionalProperties).isEquals();
    }

}
