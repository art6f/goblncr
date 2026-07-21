package balancer

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func WatchPods(balancer *Balancer) {
	factory := informers.NewSharedInformerFactoryWithOptions(
		balancer.Client,
		5*time.Second,
		informers.WithNamespace(balancer.Config.Namespace),
		informers.WithTweakListOptions(func(lo *metav1.ListOptions) {
			lo.LabelSelector = balancer.Config.PodsSelector
		}),
	)
	podInformer := factory.Core().V1().Pods().Informer()

	_, err := podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				log.Printf("Pod ADDED: %s", key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod, ok := oldObj.(*corev1.Pod)
			if !ok {
				return
			}

			newPod, ok := newObj.(*corev1.Pod)
			if !ok {
				return
			}

			if oldPod.Name != newPod.Name {
				slog.Info(fmt.Sprintf("Pod name has changed '%s' -> '%s'", oldPod.Name, newPod.Name))
			}

			if newPod.Status.PodIP != oldPod.Status.PodIP {

			}

			fmt.Printf("Pod change\nFrom\t%v: %v\nFrom\t%v: %v\n", oldPod.Name, oldPod.Status.PodIP, newPod.Name, newPod.Status.PodIP)

			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				log.Printf("Pod UPDATED: %s", key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				log.Printf("Pod DELETED: %s", key)
			}
		},
	})
	if err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		return
	}

	// Graceful shutdown requires a two-channel pattern.
	//
	// The first channel, `sigCh`, is used by the `signal` package to send us
	// OS signals (e.g., Ctrl+C). This channel must be of type `chan os.Signal`.
	//
	// The second channel, `stopCh`, is used to tell the informer factory to
	// stop. The informer factory's `Start` method expects a channel of type
	// `<-chan struct{}`. It will stop when this channel is closed.
	//
	// The goroutine below is the "translator" that connects these two channels.
	// It waits for a signal on `sigCh`, and when it receives one, it closes
	// `stopCh`, which in turn tells the informer factory to shut down.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	stopCh := make(chan struct{})
	go func() {
		<-sigCh
		close(stopCh)
	}()

	// Start the informer.
	factory.Start(stopCh)

	// Wait for the initial cache sync.
	if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced) {
		log.Println("Timed out waiting for caches to sync")
		return
	}

	log.Println("Informer has synced. Watching for Pod events...")

	// Wait for the stop signal.
	<-stopCh
	log.Println("Shutting down...")
}
