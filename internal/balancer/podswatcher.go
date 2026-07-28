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
		informers.WithNamespace(balancer.Config.Target.Namespace),
		informers.WithTweakListOptions(func(lo *metav1.ListOptions) {
			lo.LabelSelector = balancer.Config.Target.Selector
		}),
	)
	podInformer := factory.Core().V1().Pods().Informer()

	_, err := podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				log.Printf("Pod error resolving add cache: %s", err.Error())
				return
			}

			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}

			slog.Info(fmt.Sprintf("[WATCHER] New pod added: %s - %s @ %s", key, pod.Name, pod.Status.PodIP))

			balancer.pods.Add(pod)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err != nil {
				log.Printf("Pod error resolving change cache: %s", err.Error())
				return
			}

			oldPod, ok := oldObj.(*corev1.Pod)
			if !ok {
				return
			}

			newPod, ok := newObj.(*corev1.Pod)
			if !ok {
				return
			}

			if oldPod.Status.Phase != newPod.Status.Phase {
				slog.Info(fmt.Sprintf("[WATCHER] Pod '%s' phase has changed '%s' -> '%s'", key, oldPod.Status.Phase, newPod.Status.Phase))
				balancer.pods.Update(newPod)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err != nil {
				log.Printf("Pod error resolving delete cache: %s", err.Error())
				return
			}

			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}

			slog.Info(fmt.Sprintf("[WATCHER] Pod was deleted: %s - %s @ %s", key, pod.Name, pod.Status.PodIP))

			balancer.pods.Delete(pod)

			if len(balancer.pods) == 0 {
				slog.Warn("No more pods left!")
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
