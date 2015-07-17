package io.fabric8.kubernetes.api.model;

public interface Doneable<T> {

    T done();
}
