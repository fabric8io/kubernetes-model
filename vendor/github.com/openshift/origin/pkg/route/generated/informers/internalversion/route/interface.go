// This file was automatically generated by informer-gen

package route

import (
	internalinterfaces "github.com/openshift/origin/pkg/route/generated/informers/internalversion/internalinterfaces"
	internalversion "github.com/openshift/origin/pkg/route/generated/informers/internalversion/route/internalversion"
)

// Interface provides access to each of this group's versions.
type Interface interface {
	// InternalVersion provides access to shared informers for resources in InternalVersion.
	InternalVersion() internalversion.Interface
}

type group struct {
	internalinterfaces.SharedInformerFactory
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory) Interface {
	return &group{f}
}

// InternalVersion returns a new internalversion.Interface.
func (g *group) InternalVersion() internalversion.Interface {
	return internalversion.New(g.SharedInformerFactory)
}
