package io.fabric8.kubernetes.api.model;

import io.sundr.builder.annotations.ExternalBuildables;

@ExternalBuildables(builderPackage = "io.fabric8.common" ,value = {
        "io.fabric8.docker.client.dockerclient.Config",
        "io.fabric8.docker.client.dockerclient.Image"
        
})
public class DockerConfig {
}
