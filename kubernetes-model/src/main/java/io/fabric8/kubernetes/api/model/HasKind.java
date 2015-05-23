package io.fabric8.kubernetes.api.model;

import com.fasterxml.jackson.annotation.JsonSubTypes;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import io.fabric8.openshift.api.model.BuildConfig;
import io.fabric8.openshift.api.model.DeploymentConfig;
import io.fabric8.openshift.api.model.ImageStream;
import io.fabric8.openshift.api.model.OAuthAccessToken;
import io.fabric8.openshift.api.model.OAuthClient;
import io.fabric8.openshift.api.model.OAuthClientAuthorization;
import io.fabric8.openshift.api.model.Route;
import io.fabric8.openshift.api.model.template.Template;

@JsonTypeInfo(
        use = JsonTypeInfo.Id.NAME,
        include = JsonTypeInfo.As.PROPERTY,
        property = "kind")
@JsonSubTypes({
        @JsonSubTypes.Type(value = KubernetesList.class, name = "List"),
        @JsonSubTypes.Type(value = Service.class, name = "Service"),
        @JsonSubTypes.Type(value = Pod.class, name = "Pod"),
        @JsonSubTypes.Type(value = ReplicationController.class, name = "ReplicationController"),
        @JsonSubTypes.Type(value = BuildConfig.class, name = "BuildConfig"),
        @JsonSubTypes.Type(value = DeploymentConfig.class, name = "DeploymentConfig"),
        @JsonSubTypes.Type(value = ImageStream.class, name = "ImageStream"),
        @JsonSubTypes.Type(value = Route.class, name = "Route"),
        @JsonSubTypes.Type(value = Template.class, name = "Template"),
        @JsonSubTypes.Type(value = OAuthClient.class, name = "OAuthClient"),
        @JsonSubTypes.Type(value = OAuthClientAuthorization.class, name = "OAuthClientAuthorization"),
        @JsonSubTypes.Type(value = OAuthAccessToken.class, name = "OAuthAccessToken"),
        @JsonSubTypes.Type(value = Namespace.class, name = "Namespace"),
        @JsonSubTypes.Type(value = Secret.class, name = "Secrets")
})
public interface HasKind {

}
