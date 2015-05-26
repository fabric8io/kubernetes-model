package io.fabric8.kubernetes.internal;

import io.fabric8.kubernetes.api.model.HasMetadata;

import java.util.Collection;
import java.util.SortedSet;
import java.util.TreeSet;

public class HasMetadatSet extends TreeSet<HasMetadata> {

    public HasMetadatSet() {
        super(new HasMetadataComparator());
    }

    public HasMetadatSet(Collection<? extends HasMetadata> c) {
        super(new HasMetadataComparator());
        addAll(c);
    }

    public HasMetadatSet(SortedSet<HasMetadata> s) {
        super(new HasMetadataComparator());
        addAll(s);
    }

}
