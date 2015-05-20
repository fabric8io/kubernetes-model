package io.fabric8.kubernetes.api.model;

public interface HasMetadata {
  ObjectMeta getMetadata();
  void setMetadata(ObjectMeta metadata);
}
