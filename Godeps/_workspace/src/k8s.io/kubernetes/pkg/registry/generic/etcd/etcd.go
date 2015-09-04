/*
Copyright 2014 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package etcd

import (
	"fmt"
	"reflect"
	"time"

	"k8s.io/kubernetes/pkg/api"
	kubeerr "k8s.io/kubernetes/pkg/api/errors"
	etcderr "k8s.io/kubernetes/pkg/api/errors/etcd"
	"k8s.io/kubernetes/pkg/api/rest"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/generic"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/storage"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/golang/glog"
)

// Etcd implements generic.Registry, backing it with etcd storage.
// It's intended to be embeddable, so that you can implement any
// non-generic functions if needed.
// You must supply a value for every field below before use; these are
// left public as it's meant to be overridable if need be.
// This object is intended to be copyable so that it can be used in
// different ways but share the same underlying behavior.
//
// The intended use of this type is embedding within a Kind specific
// RESTStorage implementation. This type provides CRUD semantics on
// a Kubelike resource, handling details like conflict detection with
// ResourceVersion and semantics. The RESTCreateStrategy and
// RESTUpdateStrategy are generic across all backends, and encapsulate
// logic specific to the API.
//
// TODO: make the default exposed methods exactly match a generic RESTStorage
type Etcd struct {
	// Called to make a new object, should return e.g., &api.Pod{}
	NewFunc func() runtime.Object

	// Called to make a new listing object, should return e.g., &api.PodList{}
	NewListFunc func() runtime.Object

	// Used for error reporting
	EndpointName string

	// Used for listing/watching; should not include trailing "/"
	KeyRootFunc func(ctx api.Context) string

	// Called for Create/Update/Get/Delete. Note that 'namespace' can be
	// gotten from ctx.
	KeyFunc func(ctx api.Context, name string) (string, error)

	// Called to get the name of an object
	ObjectNameFunc func(obj runtime.Object) (string, error)

	// Return the TTL objects should be persisted with. Update is true if this
	// is an operation against an existing object. Existing is the current TTL
	// or the default for this operation.
	TTLFunc func(obj runtime.Object, existing uint64, update bool) (uint64, error)

	// Returns a matcher corresponding to the provided labels and fields.
	PredicateFunc func(label labels.Selector, field fields.Selector) generic.Matcher

	// Called on all objects returned from the underlying store, after
	// the exit hooks are invoked. Decorators are intended for integrations
	// that are above etcd and should only be used for specific cases where
	// storage of the value in etcd is not appropriate, since they cannot
	// be watched.
	Decorator rest.ObjectFunc
	// Allows extended behavior during creation, required
	CreateStrategy rest.RESTCreateStrategy
	// On create of an object, attempt to run a further operation.
	AfterCreate rest.ObjectFunc
	// Allows extended behavior during updates, required
	UpdateStrategy rest.RESTUpdateStrategy
	// On update of an object, attempt to run a further operation.
	AfterUpdate rest.ObjectFunc
	// Allows extended behavior during updates, optional
	DeleteStrategy rest.RESTDeleteStrategy
	// On deletion of an object, attempt to run a further operation.
	AfterDelete rest.ObjectFunc
	// If true, return the object that was deleted. Otherwise, return a generic
	// success status response.
	ReturnDeletedObject bool

	// Used for all etcd access functions
	Storage storage.Interface
}

// NamespaceKeyRootFunc is the default function for constructing etcd paths to resource directories enforcing namespace rules.
func NamespaceKeyRootFunc(ctx api.Context, prefix string) string {
	key := prefix
	ns, ok := api.NamespaceFrom(ctx)
	if ok && len(ns) > 0 {
		key = key + "/" + ns
	}
	return key
}

// NamespaceKeyFunc is the default function for constructing etcd paths to a resource relative to prefix enforcing namespace rules.
// If no namespace is on context, it errors.
func NamespaceKeyFunc(ctx api.Context, prefix string, name string) (string, error) {
	key := NamespaceKeyRootFunc(ctx, prefix)
	ns, ok := api.NamespaceFrom(ctx)
	if !ok || len(ns) == 0 {
		return "", kubeerr.NewBadRequest("Namespace parameter required.")
	}
	if len(name) == 0 {
		return "", kubeerr.NewBadRequest("Name parameter required.")
	}
	key = key + "/" + name
	return key, nil
}

// New implements RESTStorage
func (e *Etcd) New() runtime.Object {
	return e.NewFunc()
}

// NewList implements RESTLister
func (e *Etcd) NewList() runtime.Object {
	return e.NewListFunc()
}

// List returns a list of items matching labels and field
func (e *Etcd) List(ctx api.Context, label labels.Selector, field fields.Selector) (runtime.Object, error) {
	return e.ListPredicate(ctx, e.PredicateFunc(label, field))
}

// ListPredicate returns a list of all the items matching m.
func (e *Etcd) ListPredicate(ctx api.Context, m generic.Matcher) (runtime.Object, error) {
	list := e.NewListFunc()
	trace := util.NewTrace("List " + reflect.TypeOf(list).String())
	defer trace.LogIfLong(600 * time.Millisecond)
	if name, ok := m.MatchesSingle(); ok {
		trace.Step("About to read single object")
		key, err := e.KeyFunc(ctx, name)
		if err != nil {
			return nil, err
		}
		err = e.Storage.GetToList(key, list)
		trace.Step("Object extracted")
		if err != nil {
			return nil, err
		}
	} else {
		trace.Step("About to list directory")
		err := e.Storage.List(e.KeyRootFunc(ctx), list)
		trace.Step("List extracted")
		if err != nil {
			return nil, err
		}
	}
	defer trace.Step("List filtered")
	return generic.FilterList(list, m, generic.DecoratorFunc(e.Decorator))
}

// CreateWithName inserts a new item with the provided name
// DEPRECATED: use Create instead
func (e *Etcd) CreateWithName(ctx api.Context, name string, obj runtime.Object) error {
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return err
	}
	if e.CreateStrategy != nil {
		if err := rest.BeforeCreate(e.CreateStrategy, ctx, obj); err != nil {
			return err
		}
	}
	ttl, err := e.calculateTTL(obj, 0, false)
	if err != nil {
		return err
	}
	err = e.Storage.Create(key, obj, nil, ttl)
	err = etcderr.InterpretCreateError(err, e.EndpointName, name)
	if err == nil && e.Decorator != nil {
		err = e.Decorator(obj)
	}
	return err
}

// Create inserts a new item according to the unique key from the object.
func (e *Etcd) Create(ctx api.Context, obj runtime.Object) (runtime.Object, error) {
	trace := util.NewTrace("Create " + reflect.TypeOf(obj).String())
	defer trace.LogIfLong(time.Second)
	if err := rest.BeforeCreate(e.CreateStrategy, ctx, obj); err != nil {
		return nil, err
	}
	name, err := e.ObjectNameFunc(obj)
	if err != nil {
		return nil, err
	}
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}
	ttl, err := e.calculateTTL(obj, 0, false)
	if err != nil {
		return nil, err
	}
	trace.Step("About to create object")
	out := e.NewFunc()
	if err := e.Storage.Create(key, obj, out, ttl); err != nil {
		err = etcderr.InterpretCreateError(err, e.EndpointName, name)
		err = rest.CheckGeneratedNameError(e.CreateStrategy, err, obj)
		return nil, err
	}
	trace.Step("Object created")
	if e.AfterCreate != nil {
		if err := e.AfterCreate(out); err != nil {
			return nil, err
		}
	}
	if e.Decorator != nil {
		if err := e.Decorator(obj); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// UpdateWithName updates the item with the provided name
// DEPRECATED: use Update instead
func (e *Etcd) UpdateWithName(ctx api.Context, name string, obj runtime.Object) error {
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return err
	}
	ttl, err := e.calculateTTL(obj, 0, true)
	if err != nil {
		return err
	}
	err = e.Storage.Set(key, obj, nil, ttl)
	err = etcderr.InterpretUpdateError(err, e.EndpointName, name)
	if err == nil && e.Decorator != nil {
		err = e.Decorator(obj)
	}
	return err
}

// Update performs an atomic update and set of the object. Returns the result of the update
// or an error. If the registry allows create-on-update, the create flow will be executed.
// A bool is returned along with the object and any errors, to indicate object creation.
func (e *Etcd) Update(ctx api.Context, obj runtime.Object) (runtime.Object, bool, error) {
	trace := util.NewTrace("Update " + reflect.TypeOf(obj).String())
	defer trace.LogIfLong(time.Second)
	name, err := e.ObjectNameFunc(obj)
	if err != nil {
		return nil, false, err
	}
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, false, err
	}
	// If AllowUnconditionalUpdate() is true and the object specified by the user does not have a resource version,
	// then we populate it with the latest version.
	// Else, we check that the version specified by the user matches the version of latest etcd object.
	resourceVersion, err := e.Storage.Versioner().ObjectResourceVersion(obj)
	if err != nil {
		return nil, false, err
	}
	doUnconditionalUpdate := resourceVersion == 0 && e.UpdateStrategy.AllowUnconditionalUpdate()
	// TODO: expose TTL
	creating := false
	out := e.NewFunc()
	err = e.Storage.GuaranteedUpdate(key, out, true, func(existing runtime.Object, res storage.ResponseMeta) (runtime.Object, *uint64, error) {
		version, err := e.Storage.Versioner().ObjectResourceVersion(existing)
		if err != nil {
			return nil, nil, err
		}
		if version == 0 {
			if !e.UpdateStrategy.AllowCreateOnUpdate() {
				return nil, nil, kubeerr.NewNotFound(e.EndpointName, name)
			}
			creating = true
			if err := rest.BeforeCreate(e.CreateStrategy, ctx, obj); err != nil {
				return nil, nil, err
			}
			ttl, err := e.calculateTTL(obj, 0, false)
			if err != nil {
				return nil, nil, err
			}
			return obj, &ttl, nil
		}

		creating = false
		if doUnconditionalUpdate {
			// Update the object's resource version to match the latest etcd object's resource version.
			err = e.Storage.Versioner().UpdateObject(obj, res.Expiration, res.ResourceVersion)
			if err != nil {
				return nil, nil, err
			}
		} else {
			// Check if the object's resource version matches the latest resource version.
			newVersion, err := e.Storage.Versioner().ObjectResourceVersion(obj)
			if err != nil {
				return nil, nil, err
			}
			if newVersion != version {
				return nil, nil, kubeerr.NewConflict(e.EndpointName, name, fmt.Errorf("the object has been modified; please apply your changes to the latest version and try again"))
			}
		}
		if err := rest.BeforeUpdate(e.UpdateStrategy, ctx, obj, existing); err != nil {
			return nil, nil, err
		}
		ttl, err := e.calculateTTL(obj, res.TTL, true)
		if err != nil {
			return nil, nil, err
		}
		if int64(ttl) != res.TTL {
			return obj, &ttl, nil
		}
		return obj, nil, nil
	})

	if err != nil {
		if creating {
			err = etcderr.InterpretCreateError(err, e.EndpointName, name)
			err = rest.CheckGeneratedNameError(e.CreateStrategy, err, obj)
		} else {
			err = etcderr.InterpretUpdateError(err, e.EndpointName, name)
		}
		return nil, false, err
	}
	if creating {
		if e.AfterCreate != nil {
			if err := e.AfterCreate(out); err != nil {
				return nil, false, err
			}
		}
	} else {
		if e.AfterUpdate != nil {
			if err := e.AfterUpdate(out); err != nil {
				return nil, false, err
			}
		}
	}
	if e.Decorator != nil {
		if err := e.Decorator(obj); err != nil {
			return nil, false, err
		}
	}
	return out, creating, nil
}

// Get retrieves the item from etcd.
func (e *Etcd) Get(ctx api.Context, name string) (runtime.Object, error) {
	obj := e.NewFunc()
	trace := util.NewTrace("Get " + reflect.TypeOf(obj).String())
	defer trace.LogIfLong(time.Second)
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}
	trace.Step("About to read object")
	if err := e.Storage.Get(key, obj, false); err != nil {
		return nil, etcderr.InterpretGetError(err, e.EndpointName, name)
	}
	trace.Step("Object read")
	if e.Decorator != nil {
		if err := e.Decorator(obj); err != nil {
			return nil, err
		}
	}
	return obj, nil
}

// Delete removes the item from etcd.
func (e *Etcd) Delete(ctx api.Context, name string, options *api.DeleteOptions) (runtime.Object, error) {
	key, err := e.KeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}

	obj := e.NewFunc()
	trace := util.NewTrace("Delete " + reflect.TypeOf(obj).String())
	defer trace.LogIfLong(time.Second)
	trace.Step("About to read object")
	if err := e.Storage.Get(key, obj, false); err != nil {
		return nil, etcderr.InterpretDeleteError(err, e.EndpointName, name)
	}

	// support older consumers of delete by treating "nil" as delete immediately
	if options == nil {
		options = api.NewDeleteOptions(0)
	}
	graceful, pendingGraceful, err := rest.BeforeDelete(e.DeleteStrategy, ctx, obj, options)
	if err != nil {
		return nil, err
	}
	if pendingGraceful {
		return e.finalizeDelete(obj, false)
	}
	if graceful && *options.GracePeriodSeconds != 0 {
		trace.Step("Graceful deletion")
		out := e.NewFunc()
		if err := e.Storage.Set(key, obj, out, uint64(*options.GracePeriodSeconds)); err != nil {
			return nil, etcderr.InterpretUpdateError(err, e.EndpointName, name)
		}
		return e.finalizeDelete(out, true)
	}

	// delete immediately, or no graceful deletion supported
	out := e.NewFunc()
	trace.Step("About to delete object")
	if err := e.Storage.Delete(key, out); err != nil {
		return nil, etcderr.InterpretDeleteError(err, e.EndpointName, name)
	}
	return e.finalizeDelete(out, true)
}

func (e *Etcd) finalizeDelete(obj runtime.Object, runHooks bool) (runtime.Object, error) {
	if runHooks && e.AfterDelete != nil {
		if err := e.AfterDelete(obj); err != nil {
			return nil, err
		}
	}
	if e.ReturnDeletedObject {
		if e.Decorator != nil {
			if err := e.Decorator(obj); err != nil {
				return nil, err
			}
		}
		return obj, nil
	}
	return &api.Status{Status: api.StatusSuccess}, nil
}

// Watch makes a matcher for the given label and field, and calls
// WatchPredicate. If possible, you should customize PredicateFunc to produre a
// matcher that matches by key. generic.SelectionPredicate does this for you
// automatically.
func (e *Etcd) Watch(ctx api.Context, label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	return e.WatchPredicate(ctx, e.PredicateFunc(label, field), resourceVersion)
}

// WatchPredicate starts a watch for the items that m matches.
func (e *Etcd) WatchPredicate(ctx api.Context, m generic.Matcher, resourceVersion string) (watch.Interface, error) {
	version, err := storage.ParseWatchResourceVersion(resourceVersion, e.EndpointName)
	if err != nil {
		return nil, err
	}

	filterFunc := func(obj runtime.Object) bool {
		matches, err := m.Matches(obj)
		if err != nil {
			glog.Errorf("unable to match watch: %v", err)
			return false
		}
		if matches && e.Decorator != nil {
			if err := e.Decorator(obj); err != nil {
				glog.Errorf("unable to decorate watch: %v", err)
				return false
			}
		}
		return matches
	}

	if name, ok := m.MatchesSingle(); ok {
		key, err := e.KeyFunc(ctx, name)
		if err != nil {
			return nil, err
		}
		return e.Storage.Watch(key, version, filterFunc)
	}

	return e.Storage.WatchList(e.KeyRootFunc(ctx), version, filterFunc)
}

// calculateTTL is a helper for retrieving the updated TTL for an object or returning an error
// if the TTL cannot be calculated. The defaultTTL is changed to 1 if less than zero. Zero means
// no TTL, not expire immediately.
func (e *Etcd) calculateTTL(obj runtime.Object, defaultTTL int64, update bool) (ttl uint64, err error) {
	// etcd may return a negative TTL for a node if the expiration has not occurred due
	// to server lag - we will ensure that the value is at least set.
	if defaultTTL < 0 {
		defaultTTL = 1
	}
	ttl = uint64(defaultTTL)
	if e.TTLFunc != nil {
		ttl, err = e.TTLFunc(obj, ttl, update)
	}
	return ttl, err
}
