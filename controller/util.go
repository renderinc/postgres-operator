package controller

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func WrapListWatchWithErrorHandler(listWatch *cache.ListWatch, handleError func(error)) *cache.ListWatch {
	return &cache.ListWatch{
		ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
			obj, err := listWatch.ListFunc(options)
			if err != nil {
				handleError(err)
			}
			return obj, err
		},
		WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
			i, err := listWatch.WatchFunc(options)
			if err != nil {
				handleError(err)
			}
			return i, err
		},
		DisableChunking: listWatch.DisableChunking,
	}
}
