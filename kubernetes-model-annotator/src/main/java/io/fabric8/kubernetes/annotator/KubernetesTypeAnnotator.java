/**
 * Copyright (C) 2011 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package io.fabric8.kubernetes.annotator;

import com.fasterxml.jackson.databind.JsonDeserializer;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;
import com.sun.codemodel.JAnnotationUse;
import com.sun.codemodel.JClassAlreadyExistsException;
import com.sun.codemodel.JCodeModel;
import com.sun.codemodel.JDefinedClass;
import com.sun.codemodel.JEnumConstant;
import com.sun.codemodel.JFieldVar;
import com.sun.codemodel.JMethod;
import io.sundr.builder.annotations.Buildable;
import io.sundr.builder.annotations.Inline;
import lombok.EqualsAndHashCode;
import lombok.ToString;
import org.jsonschema2pojo.Jackson2Annotator;

import java.util.Map;

public class KubernetesTypeAnnotator extends Jackson2Annotator {

    private final String nameIsDNS952LabelPattern = "[a-z]([-a-z0-9]*[a-z0-9])?";
    private final int nameIsDNS952LabelLength = 24;

    private final String nameIsDNS1123LabelPattern = "[a-z0-9]([-a-z0-9]*[a-z0-9])?";
    private final int nameIsDNS1123LabelLength = 63;

    private final String nameIsDNS1123SubdomainPattern = nameIsDNS1123LabelPattern + "(\\." + nameIsDNS1123LabelPattern + ")*";
    private final int nameIsDNS1123SubdomainLength = 253;

    @Override
    public void propertyOrder(JDefinedClass clazz, JsonNode propertiesNode) {
        //We just want to make sure we avoid infinite loops
        clazz.annotate(JsonDeserialize.class)
                .param("using", JsonDeserializer.None.class);
        clazz.annotate(ToString.class);
        clazz.annotate(EqualsAndHashCode.class);
        try {
            clazz.annotate(Buildable.class)
                    .param("editableEnabled", false)
                    .param("validationEnabled", true)
                    .param("generateBuilderPackage", true)
                    .param("builderPackage", "io.fabric8.kubernetes.api.builder")
                    .annotationParam("inline", Inline.class)
                    .param("type", new JCodeModel()._class("io.fabric8.kubernetes.api.model.Doneable"))
                    .param("prefix", "Doneable")
                    .param("value", "done");

            annotateMetatadataValidator(clazz);
        } catch (JClassAlreadyExistsException e) {
            e.printStackTrace();
        }
    }

    private int getObjectNameMaxLength(JDefinedClass clazz) {
        String kind = clazz.name();
        if (kind.equals("Service")) {
            return nameIsDNS952LabelLength;
        }
        if (kind.equals("Namespace") || kind.equals("Project")) {
            return nameIsDNS1123LabelLength;
        }
        return nameIsDNS1123SubdomainLength;
    }

    private String getObjectNamePattern(JDefinedClass clazz) {
        String kind = clazz.name();
        if (kind.equals("Service")) {
            return nameIsDNS952LabelPattern;
        }
        if (kind.equals("Namespace") || kind.equals("Project")) {
            return nameIsDNS1123LabelPattern;
        }
        return nameIsDNS1123SubdomainPattern;
    }

    private void annotateMetatadataValidator(JDefinedClass clazz) {
        if (clazz.name().equals("PodTemplateSpec")) {
            return;
        }

        for (Map.Entry<String, JFieldVar> f : clazz.fields().entrySet()) {
            if (f.getKey().equals("metadata") && f.getValue().type().name().equals("ObjectMeta")) {
                try {
                    JAnnotationUse annotation = f.getValue().annotate(new JCodeModel()._class("io.fabric8.kubernetes.api.model.validators.CheckObjectMeta"));

                    if (isMinimal(clazz)) {
                        annotation.param("minimal", true);
                    } else {
                        annotation
                            .param("regexp", "^" + getObjectNamePattern(clazz) + "$")
                            .param("max", getObjectNameMaxLength(clazz));
                    }
                } catch (JClassAlreadyExistsException e) {
                    e.printStackTrace();
                }
                return;
            }
        }
    }

    private boolean isMinimal(JDefinedClass clazz) {
        String kind = clazz.name();
        return kind.equals("Group") || kind.equals("User") || kind.equals("Identity") || kind.equals("UserIdentityMapping")
            || kind.equals("ClusterNetwork") || kind.equals("HostSubnet") || kind.equals("NetNamespace")
            || kind.equals("Image") || kind.equals("ImageStream") || kind.equals("ImageStreamMapping") || kind.equals("ImageStreamTag") || kind.equals("ImageStreamImport")
            || kind.equals("Policy") || kind.equals("PolicyBinding") || kind.equals("Role") || kind.equals("RoleBinding")
            || kind.equals("OAuthAccessToken") || kind.equals("OAuthAuthorizeToken") || kind.equals("OAuthClientAuthorization");
    }

    @Override
    public void propertyInclusion(JDefinedClass clazz, JsonNode schema) {

    }

    @Override
    public void propertyField(JFieldVar field, JDefinedClass clazz, String propertyName, JsonNode propertyNode) {

    }

    @Override
    public void propertyGetter(JMethod getter, String propertyName) {

    }

    @Override
    public void propertySetter(JMethod setter, String propertyName) {

    }

    @Override
    public void anyGetter(JMethod getter) {

    }

    @Override
    public void anySetter(JMethod setter) {

    }

    @Override
    public void enumCreatorMethod(JMethod creatorMethod) {

    }

    @Override
    public void enumValueMethod(JMethod valueMethod) {

    }

    @Override
    public void enumConstant(JEnumConstant constant, String value) {

    }

    @Override
    public boolean isAdditionalPropertiesSupported() {
        return true;
    }

    @Override
    public void additionalPropertiesField(JFieldVar field, JDefinedClass clazz, String propertyName) {

    }
}
