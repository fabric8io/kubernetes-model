package io.fabric8.kubernetes.internal;

import io.fabric8.kubernetes.api.model.HasMetadata;

import java.util.Comparator;

public class HasMetadataComparator implements Comparator<HasMetadata> {

    private enum KindOrder {
        Secret,
        ServiceAccount,
        OAuthClient,
        Service;
    }

    private Integer getKindValue(String kind) {
        try {
            KindOrder kindOrder = KindOrder.valueOf(kind);
            switch (kindOrder) {
                case Secret:
                    return 0;
                case ServiceAccount:
                    return 1;
                case OAuthClient:
                    return 2;
                case Service:
                    return 3;
                default:
                    return 100;
            }
        } catch (IllegalArgumentException e) {
            return 100;
        }
    }

    @Override
    public int compare(HasMetadata a, HasMetadata b) {
        if (a == null || b == null) {
            throw new NullPointerException("Cannot compare null HasMetadata objects");
        }
        if (a == b || a.equals(b)) {
            return 0;
        }

        String classNameA = a.getClass().getSimpleName();
        String classNameB = b.getClass().getSimpleName();

        int kindOrderCompare = getKindValue(classNameA).compareTo(getKindValue(classNameB));
        if (kindOrderCompare != 0) {
            return kindOrderCompare;
        }

        int classCompare = classNameA.compareTo(classNameB);
        if (classCompare != 0) {
            return classCompare;
        }
        return a.getMetadata().getName().compareTo(b.getMetadata().getName());
    }
}
