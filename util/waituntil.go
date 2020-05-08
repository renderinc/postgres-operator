package util

/*
 Copyright 2017-2018 Crunchy Data Solutions, Inc.
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

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	twatch "k8s.io/client-go/tools/watch"
)

// WaitUntilPod ...
//podPhase is v1.PodRunning
//timeout := time.Minute
func WaitUntilPod(clientset *kubernetes.Clientset, lo meta_v1.ListOptions, podPhase v1.PodPhase, timeout time.Duration, namespace string) error {
	var err error

	pods := clientset.CoreV1().Pods(namespace)
	lw := cache.ListWatch{
		ListFunc: func(_ meta_v1.ListOptions) (runtime.Object, error) {
			return pods.List(lo)
		},
		WatchFunc: func(_ meta_v1.ListOptions) (watch.Interface, error) {
			return pods.Watch(lo)
		},
	}

	conditions := []twatch.ConditionFunc{
		func(event watch.Event) (bool, error) {
			log.Debug("watch Modified called")
			gotpod2 := event.Object.(*v1.Pod)
			log.Debug("pod2 phase=" + gotpod2.Status.Phase)
			if gotpod2.Status.Phase == podPhase {
				return true, nil
			}
			return false, nil
		},
	}

	log.Debug("before watch.Until")

	var lastEvent *watch.Event
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	lastEvent, err = twatch.ListWatchUntil(ctx, &lw, conditions...)
	if err != nil {
		log.Errorf("timeout waiting for %v error=%s", podPhase, err.Error())
		return err
	}
	if lastEvent == nil {
		log.Error("expected event")
		return err
	}
	log.Debug("after watch.Until")
	return err

}

// WaitUntilPodIsDeleted timeout := time.Minute

// WaitUntilDeploymentIsDeleted timeout := time.Minute
