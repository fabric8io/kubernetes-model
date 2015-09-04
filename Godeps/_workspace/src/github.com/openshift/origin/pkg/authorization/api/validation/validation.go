package validation

import (
	"fmt"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/validation"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/util/fielderrors"

	oapi "github.com/openshift/origin/pkg/api"
	authorizationapi "github.com/openshift/origin/pkg/authorization/api"
	uservalidation "github.com/openshift/origin/pkg/user/api/validation"
)

func ValidateSubjectAccessReview(review *authorizationapi.SubjectAccessReview) fielderrors.ValidationErrorList {
	allErrs := fielderrors.ValidationErrorList{}

	if len(review.Action.Verb) == 0 {
		allErrs = append(allErrs, fielderrors.NewFieldRequired("verb"))
	}
	if len(review.Action.Resource) == 0 {
		allErrs = append(allErrs, fielderrors.NewFieldRequired("resource"))
	}

	return allErrs
}

func ValidateResourceAccessReview(review *authorizationapi.ResourceAccessReview) fielderrors.ValidationErrorList {
	allErrs := fielderrors.ValidationErrorList{}

	if len(review.Action.Verb) == 0 {
		allErrs = append(allErrs, fielderrors.NewFieldRequired("verb"))
	}
	if len(review.Action.Resource) == 0 {
		allErrs = append(allErrs, fielderrors.NewFieldRequired("resource"))
	}

	return allErrs
}

func ValidateLocalSubjectAccessReview(review *authorizationapi.LocalSubjectAccessReview) fielderrors.ValidationErrorList {
	allErrs := fielderrors.ValidationErrorList{}

	if len(review.Action.Verb) == 0 {
		allErrs = append(allErrs, fielderrors.NewFieldRequired("verb"))
	}
	if len(review.Action.Resource) == 0 {
		allErrs = append(allErrs, fielderrors.NewFieldRequired("resource"))
	}

	return allErrs
}

func ValidateLocalResourceAccessReview(review *authorizationapi.LocalResourceAccessReview) fielderrors.ValidationErrorList {
	allErrs := fielderrors.ValidationErrorList{}

	if len(review.Action.Verb) == 0 {
		allErrs = append(allErrs, fielderrors.NewFieldRequired("verb"))
	}
	if len(review.Action.Resource) == 0 {
		allErrs = append(allErrs, fielderrors.NewFieldRequired("resource"))
	}

	return allErrs
}

func ValidatePolicyName(name string, prefix bool) (bool, string) {
	if ok, reason := oapi.MinimalNameRequirements(name, prefix); !ok {
		return ok, reason
	}

	if name != authorizationapi.PolicyName {
		return false, "name must be " + authorizationapi.PolicyName
	}

	return true, ""
}

func ValidateLocalPolicy(policy *authorizationapi.Policy) fielderrors.ValidationErrorList {
	return ValidatePolicy(policy, true)
}

func ValidateLocalPolicyUpdate(policy *authorizationapi.Policy, oldPolicy *authorizationapi.Policy) fielderrors.ValidationErrorList {
	return ValidatePolicyUpdate(policy, oldPolicy, true)
}

func ValidateClusterPolicy(policy *authorizationapi.ClusterPolicy) fielderrors.ValidationErrorList {
	return ValidatePolicy(authorizationapi.ToPolicy(policy), false)
}

func ValidateClusterPolicyUpdate(policy *authorizationapi.ClusterPolicy, oldPolicy *authorizationapi.ClusterPolicy) fielderrors.ValidationErrorList {
	return ValidatePolicyUpdate(authorizationapi.ToPolicy(policy), authorizationapi.ToPolicy(oldPolicy), false)
}

func ValidatePolicy(policy *authorizationapi.Policy, isNamespaced bool) fielderrors.ValidationErrorList {
	allErrs := fielderrors.ValidationErrorList{}
	allErrs = append(allErrs, validation.ValidateObjectMeta(&policy.ObjectMeta, isNamespaced, ValidatePolicyName).Prefix("metadata")...)

	for roleKey, role := range policy.Roles {
		if role == nil {
			allErrs = append(allErrs, fielderrors.NewFieldRequired("roles."+roleKey))
		}

		if roleKey != role.Name {
			allErrs = append(allErrs, fielderrors.NewFieldInvalid("roles."+roleKey+".metadata.name", role.Name, "must be "+roleKey))
		}

		allErrs = append(allErrs, ValidateRole(role, isNamespaced).Prefix("roles."+roleKey)...)
	}

	return allErrs
}

func ValidatePolicyUpdate(policy *authorizationapi.Policy, oldPolicy *authorizationapi.Policy, isNamespaced bool) fielderrors.ValidationErrorList {
	allErrs := ValidatePolicy(policy, isNamespaced)
	allErrs = append(allErrs, validation.ValidateObjectMetaUpdate(&policy.ObjectMeta, &oldPolicy.ObjectMeta).Prefix("metadata")...)

	return allErrs
}

func PolicyBindingNameValidator(policyRefNamespace string) validation.ValidateNameFunc {
	return func(name string, prefix bool) (bool, string) {
		if ok, reason := oapi.MinimalNameRequirements(name, prefix); !ok {
			return ok, reason
		}

		if name != authorizationapi.GetPolicyBindingName(policyRefNamespace) {
			return false, "name must be " + authorizationapi.GetPolicyBindingName(policyRefNamespace)
		}

		return true, ""
	}
}

func ValidateLocalPolicyBinding(policy *authorizationapi.PolicyBinding) fielderrors.ValidationErrorList {
	return ValidatePolicyBinding(policy, true)
}

func ValidateLocalPolicyBindingUpdate(policy *authorizationapi.PolicyBinding, oldPolicyBinding *authorizationapi.PolicyBinding) fielderrors.ValidationErrorList {
	return ValidatePolicyBindingUpdate(policy, oldPolicyBinding, true)
}

func ValidateClusterPolicyBinding(policy *authorizationapi.ClusterPolicyBinding) fielderrors.ValidationErrorList {
	return ValidatePolicyBinding(authorizationapi.ToPolicyBinding(policy), false)
}

func ValidateClusterPolicyBindingUpdate(policy *authorizationapi.ClusterPolicyBinding, oldPolicyBinding *authorizationapi.ClusterPolicyBinding) fielderrors.ValidationErrorList {
	return ValidatePolicyBindingUpdate(authorizationapi.ToPolicyBinding(policy), authorizationapi.ToPolicyBinding(oldPolicyBinding), false)
}

func ValidatePolicyBinding(policyBinding *authorizationapi.PolicyBinding, isNamespaced bool) fielderrors.ValidationErrorList {
	allErrs := fielderrors.ValidationErrorList{}
	allErrs = append(allErrs, validation.ValidateObjectMeta(&policyBinding.ObjectMeta, isNamespaced, PolicyBindingNameValidator(policyBinding.PolicyRef.Namespace)).Prefix("metadata")...)

	if !isNamespaced {
		if len(policyBinding.PolicyRef.Namespace) > 0 {
			allErrs = append(allErrs, fielderrors.NewFieldInvalid("policyRef.namespace", policyBinding.PolicyRef.Namespace, "may not reference another namespace"))
		}
	}

	for roleBindingKey, roleBinding := range policyBinding.RoleBindings {
		if roleBinding == nil {
			allErrs = append(allErrs, fielderrors.NewFieldRequired("roleBindings."+roleBindingKey))
		}

		if roleBinding.RoleRef.Namespace != policyBinding.PolicyRef.Namespace {
			allErrs = append(allErrs, fielderrors.NewFieldInvalid("roleBindings."+roleBindingKey+".roleRef.namespace", policyBinding.PolicyRef.Namespace, "must be "+policyBinding.PolicyRef.Namespace))
		}

		if roleBindingKey != roleBinding.Name {
			allErrs = append(allErrs, fielderrors.NewFieldInvalid("roleBindings."+roleBindingKey+".metadata.name", roleBinding.Name, "must be "+roleBindingKey))
		}

		allErrs = append(allErrs, ValidateRoleBinding(roleBinding, isNamespaced).Prefix("roleBindings."+roleBindingKey)...)
	}

	return allErrs
}

func ValidatePolicyBindingUpdate(policyBinding *authorizationapi.PolicyBinding, oldPolicyBinding *authorizationapi.PolicyBinding, isNamespaced bool) fielderrors.ValidationErrorList {
	allErrs := ValidatePolicyBinding(policyBinding, isNamespaced)
	allErrs = append(allErrs, validation.ValidateObjectMetaUpdate(&policyBinding.ObjectMeta, &oldPolicyBinding.ObjectMeta).Prefix("metadata")...)

	if oldPolicyBinding.PolicyRef.Namespace != policyBinding.PolicyRef.Namespace {
		allErrs = append(allErrs, fielderrors.NewFieldInvalid("policyRef.namespace", policyBinding.PolicyRef.Namespace, "cannot change policyRef"))
	}

	return allErrs
}

func ValidateLocalRole(policy *authorizationapi.Role) fielderrors.ValidationErrorList {
	return ValidateRole(policy, true)
}

func ValidateLocalRoleUpdate(policy *authorizationapi.Role, oldRole *authorizationapi.Role) fielderrors.ValidationErrorList {
	return ValidateRoleUpdate(policy, oldRole, true)
}

func ValidateClusterRole(policy *authorizationapi.ClusterRole) fielderrors.ValidationErrorList {
	return ValidateRole(authorizationapi.ToRole(policy), false)
}

func ValidateClusterRoleUpdate(policy *authorizationapi.ClusterRole, oldRole *authorizationapi.ClusterRole) fielderrors.ValidationErrorList {
	return ValidateRoleUpdate(authorizationapi.ToRole(policy), authorizationapi.ToRole(oldRole), false)
}

func ValidateRole(role *authorizationapi.Role, isNamespaced bool) fielderrors.ValidationErrorList {
	allErrs := fielderrors.ValidationErrorList{}
	allErrs = append(allErrs, validation.ValidateObjectMeta(&role.ObjectMeta, isNamespaced, oapi.MinimalNameRequirements).Prefix("metadata")...)

	return allErrs
}

func ValidateRoleUpdate(role *authorizationapi.Role, oldRole *authorizationapi.Role, isNamespaced bool) fielderrors.ValidationErrorList {
	allErrs := ValidateRole(role, isNamespaced)
	allErrs = append(allErrs, validation.ValidateObjectMetaUpdate(&role.ObjectMeta, &oldRole.ObjectMeta).Prefix("metadata")...)

	return allErrs
}

func ValidateLocalRoleBinding(policy *authorizationapi.RoleBinding) fielderrors.ValidationErrorList {
	return ValidateRoleBinding(policy, true)
}

func ValidateLocalRoleBindingUpdate(policy *authorizationapi.RoleBinding, oldRoleBinding *authorizationapi.RoleBinding) fielderrors.ValidationErrorList {
	return ValidateRoleBindingUpdate(policy, oldRoleBinding, true)
}

func ValidateClusterRoleBinding(policy *authorizationapi.ClusterRoleBinding) fielderrors.ValidationErrorList {
	return ValidateRoleBinding(authorizationapi.ToRoleBinding(policy), false)
}

func ValidateClusterRoleBindingUpdate(policy *authorizationapi.ClusterRoleBinding, oldRoleBinding *authorizationapi.ClusterRoleBinding) fielderrors.ValidationErrorList {
	return ValidateRoleBindingUpdate(authorizationapi.ToRoleBinding(policy), authorizationapi.ToRoleBinding(oldRoleBinding), false)
}

func ValidateRoleBinding(roleBinding *authorizationapi.RoleBinding, isNamespaced bool) fielderrors.ValidationErrorList {
	allErrs := fielderrors.ValidationErrorList{}
	allErrs = append(allErrs, validation.ValidateObjectMeta(&roleBinding.ObjectMeta, isNamespaced, oapi.MinimalNameRequirements).Prefix("metadata")...)

	// roleRef namespace is empty when referring to global policy.
	if (len(roleBinding.RoleRef.Namespace) > 0) && !util.IsDNS1123Subdomain(roleBinding.RoleRef.Namespace) {
		allErrs = append(allErrs, fielderrors.NewFieldInvalid("roleRef.namespace", roleBinding.RoleRef.Namespace, "roleRef.namespace must be a valid subdomain"))
	}

	if len(roleBinding.RoleRef.Name) == 0 {
		allErrs = append(allErrs, fielderrors.NewFieldRequired("roleRef.name"))
	} else {
		if valid, err := oapi.MinimalNameRequirements(roleBinding.RoleRef.Name, false); !valid {
			allErrs = append(allErrs, fielderrors.NewFieldInvalid("roleRef.name", roleBinding.RoleRef.Name, err))
		}
	}

	for i, subject := range roleBinding.Subjects {
		allErrs = append(allErrs, ValidateRoleBindingSubject(subject, isNamespaced).Prefix(fmt.Sprintf("subjects[%d]", i))...)
	}

	return allErrs
}

func ValidateRoleBindingSubject(subject kapi.ObjectReference, isNamespaced bool) fielderrors.ValidationErrorList {
	allErrs := fielderrors.ValidationErrorList{}

	if len(subject.Name) == 0 {
		allErrs = append(allErrs, fielderrors.NewFieldRequired("name"))
	}
	if len(subject.UID) != 0 {
		allErrs = append(allErrs, fielderrors.NewFieldForbidden("uid", subject.UID))
	}
	if len(subject.APIVersion) != 0 {
		allErrs = append(allErrs, fielderrors.NewFieldForbidden("apiVersion", subject.APIVersion))
	}
	if len(subject.ResourceVersion) != 0 {
		allErrs = append(allErrs, fielderrors.NewFieldForbidden("resourceVersion", subject.ResourceVersion))
	}
	if len(subject.FieldPath) != 0 {
		allErrs = append(allErrs, fielderrors.NewFieldForbidden("fieldPath", subject.FieldPath))
	}

	switch subject.Kind {
	case authorizationapi.ServiceAccountKind:
		if valid, reason := validation.ValidateServiceAccountName(subject.Name, false); len(subject.Name) > 0 && !valid {
			allErrs = append(allErrs, fielderrors.NewFieldInvalid("name", subject.Name, reason))
		}
		if !isNamespaced && len(subject.Namespace) == 0 {
			allErrs = append(allErrs, fielderrors.NewFieldRequired("namespace"))
		}

	case authorizationapi.UserKind:
		if valid, reason := uservalidation.ValidateUserName(subject.Name, false); len(subject.Name) > 0 && !valid {
			allErrs = append(allErrs, fielderrors.NewFieldInvalid("name", subject.Name, reason))
		}

	case authorizationapi.GroupKind:
		if valid, reason := uservalidation.ValidateGroupName(subject.Name, false); len(subject.Name) > 0 && !valid {
			allErrs = append(allErrs, fielderrors.NewFieldInvalid("name", subject.Name, reason))
		}

	case authorizationapi.SystemUserKind:
		isValidSAName, _ := validation.ValidateServiceAccountName(subject.Name, false)
		isValidUserName, _ := uservalidation.ValidateUserName(subject.Name, false)
		if isValidSAName || isValidUserName {
			allErrs = append(allErrs, fielderrors.NewFieldInvalid("name", subject.Name, "conforms to User.name or ServiceAccount.name restrictions"))
		}

	case authorizationapi.SystemGroupKind:
		if valid, _ := uservalidation.ValidateGroupName(subject.Name, false); len(subject.Name) > 0 && valid {
			allErrs = append(allErrs, fielderrors.NewFieldInvalid("name", subject.Name, "conforms to Group.name restrictions"))
		}

	default:
		allErrs = append(allErrs, fielderrors.NewFieldValueNotSupported("kind", subject.Kind, []string{authorizationapi.ServiceAccountKind, authorizationapi.UserKind, authorizationapi.GroupKind, authorizationapi.SystemGroupKind, authorizationapi.SystemUserKind}))
	}

	return allErrs
}

func ValidateRoleBindingUpdate(roleBinding *authorizationapi.RoleBinding, oldRoleBinding *authorizationapi.RoleBinding, isNamespaced bool) fielderrors.ValidationErrorList {
	allErrs := ValidateRoleBinding(roleBinding, isNamespaced)
	allErrs = append(allErrs, validation.ValidateObjectMetaUpdate(&roleBinding.ObjectMeta, &oldRoleBinding.ObjectMeta).Prefix("metadata")...)

	if oldRoleBinding.RoleRef != roleBinding.RoleRef {
		allErrs = append(allErrs, fielderrors.NewFieldInvalid("roleRef", roleBinding.RoleRef, "cannot change roleRef"))
	}

	return allErrs
}
