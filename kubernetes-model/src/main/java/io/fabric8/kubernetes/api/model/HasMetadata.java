package io.fabric8.kubernetes.api.model;

public interface HasMetadata extends HasKind {

  ObjectMeta getMetadata();
  void setMetadata(ObjectMeta metadata);
}
