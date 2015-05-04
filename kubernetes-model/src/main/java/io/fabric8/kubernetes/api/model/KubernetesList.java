
package io.fabric8.kubernetes.api.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonPropertyOrder;
import io.fabric8.openshift.api.model.BuildConfig;
import io.fabric8.openshift.api.model.DeploymentConfig;
import io.fabric8.openshift.api.model.ImageStream;
import io.fabric8.openshift.api.model.Route;
import io.fabric8.openshift.api.model.template.Template;

import javax.annotation.Generated;
import java.util.*;


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
    "deletionTimestamp",
    "generateName",
    "id",
    "items",
    "kind",
    "namespace",
    "resourceVersion",
    "selfLink",
    "uid"
})
public class KubernetesList extends BaseKubernetesList {

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
    private final List<ImageStream> imageStreams = Collections.emptyList();

    /**
     * list of routes
     * Note: This is not to be used. Added for influencing the generation of fluent nested builders.
     */
    @JsonIgnore
    private final List<Route> routes = Collections.emptyList();

    /**
     * list of routes
     * Note: This is not to be used. Added for influencing the generation of fluent nested builders.
     */
    @JsonIgnore
    private final List<Template> templates = Collections.emptyList();
    
    /**
     * No args constructor for use in serialization
     *
     */
    public KubernetesList() {
        super();
    }

    public KubernetesList(KubernetesList.ApiVersion apiVersion,
                          List<Object> items,
                          String kind,
                          String resourceVersion,
                          String selfLink,
                          List<Service> services,
                          List<ReplicationController> replicationControllers,
                          List<Pod> pods,
                          List<BuildConfig> buildConfigs,
                          List<DeploymentConfig> deploymentConfigs,
                          List<ImageStream> imageStreams,
                          List<Route> routes,
                          List<Template> templates) {
        super(apiVersion, items, kind, resourceVersion, selfLink);
        Set<Object> allItems = new LinkedHashSet<>();
        allItems.addAll(items != null ? items : Collections.emptyList());
        allItems.addAll(services != null ? services : Collections.<Service>emptyList());
        allItems.addAll(replicationControllers != null ? replicationControllers : Collections.<ReplicationController>emptyList());
        allItems.addAll(pods != null ? pods : Collections.<Pod>emptyList());
        allItems.addAll(buildConfigs != null ? buildConfigs : Collections.<BuildConfig>emptyList());
        allItems.addAll(deploymentConfigs != null ? deploymentConfigs : Collections.<DeploymentConfig>emptyList());
        allItems.addAll(imageStreams != null ? imageStreams : Collections.<ImageStream>emptyList());
        allItems.addAll(routes != null ? routes : Collections.<Route>emptyList());
        allItems.addAll(templates != null ? templates : Collections.<Template>emptyList());
        setItems(new ArrayList<>(allItems));
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
    public List<ImageStream> getImageStreams() {
        return imageStreams;
    }
    
    @JsonIgnore
    public List<Route> getRoutes() {
        return routes;
    }

    @JsonIgnore
    public List<Template> getTemplates() {
        return templates;
    }

}
