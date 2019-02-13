// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	time "time"

	gardenv1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	versioned "github.com/gardener/gardener/pkg/client/garden/clientset/versioned"
	internalinterfaces "github.com/gardener/gardener/pkg/client/garden/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/gardener/gardener/pkg/client/garden/listers/garden/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CloudProfileInformer provides access to a shared informer and lister for
// CloudProfiles.
type CloudProfileInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.CloudProfileLister
}

type cloudProfileInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewCloudProfileInformer constructs a new informer for CloudProfile type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCloudProfileInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCloudProfileInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredCloudProfileInformer constructs a new informer for CloudProfile type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCloudProfileInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GardenV1beta1().CloudProfiles().List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GardenV1beta1().CloudProfiles().Watch(options)
			},
		},
		&gardenv1beta1.CloudProfile{},
		resyncPeriod,
		indexers,
	)
}

func (f *cloudProfileInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCloudProfileInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *cloudProfileInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&gardenv1beta1.CloudProfile{}, f.defaultInformer)
}

func (f *cloudProfileInformer) Lister() v1beta1.CloudProfileLister {
	return v1beta1.NewCloudProfileLister(f.Informer().GetIndexer())
}
