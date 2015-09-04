package api

import kapi "k8s.io/kubernetes/pkg/api"

// Auth system gets identity name and provider
// POST to UserIdentityMapping, get back error or a filled out UserIdentityMapping object

type User struct {
	kapi.TypeMeta
	kapi.ObjectMeta

	FullName string

	Identities []string

	Groups []string
}

type UserList struct {
	kapi.TypeMeta
	kapi.ListMeta
	Items []User
}

type Identity struct {
	kapi.TypeMeta
	kapi.ObjectMeta

	// ProviderName is the source of identity information
	ProviderName string

	// ProviderUserName uniquely represents this identity in the scope of the provider
	ProviderUserName string

	// User is a reference to the user this identity is associated with
	// Both Name and UID must be set
	User kapi.ObjectReference

	Extra map[string]string
}

type IdentityList struct {
	kapi.TypeMeta
	kapi.ListMeta
	Items []Identity
}

type UserIdentityMapping struct {
	kapi.TypeMeta
	kapi.ObjectMeta

	Identity kapi.ObjectReference
	User     kapi.ObjectReference
}

// Group represents a referenceable set of Users
type Group struct {
	kapi.TypeMeta
	kapi.ObjectMeta

	Users []string
}

type GroupList struct {
	kapi.TypeMeta
	kapi.ListMeta
	Items []Group
}

func (*GroupList) IsAnAPIObject()           {}
func (*Group) IsAnAPIObject()               {}
func (*User) IsAnAPIObject()                {}
func (*UserList) IsAnAPIObject()            {}
func (*Identity) IsAnAPIObject()            {}
func (*IdentityList) IsAnAPIObject()        {}
func (*UserIdentityMapping) IsAnAPIObject() {}
