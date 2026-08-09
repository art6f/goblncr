// Package balancer implements the load balancer.
package balancer

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/art6f/goblncr/config"
	"github.com/art6f/goblncr/internal/balancer/strategies"
	"github.com/art6f/goblncr/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Balancer struct {
	Client   *kubernetes.Clientset
	Config   *config.AppConfig
	pods     k8s.PodsMap
	strategy strategies.Strategy
}

func NewBalancer(config *config.AppConfig, strategy strategies.Strategy) *Balancer {
	return &Balancer{
		Client:   k8s.NewClient(),
		Config:   config,
		pods:     make(k8s.PodsMap, 0),
		strategy: strategy,
	}
}

func (balancer *Balancer) Run() {
	slog.Info("Starting Balancer...")
	slog.Info("Querying pods", "namespace", balancer.Config.Target.Namespace, "label", balancer.Config.Target.Selector)

	pods, err := balancer.Client.CoreV1().Pods(balancer.Config.Target.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: balancer.Config.Target.Selector,
	})

	if err != nil {
		slog.Error("Errror retrieving pods", "errro", err.Error())
		panic(err.Error())
	}

	if len(pods.Items) == 0 {
		slog.Warn("No pods found! Balancer will watch the pods and idle...")
	}

	balancer.strategy.SetPods(&balancer.pods)

	k8s.WatchPods(balancer.Config.Target.Namespace, balancer.Config.Target.Selector, balancer.Client, &balancer.pods)
}

func (balancer *Balancer) SelectServer(_ *http.Request) (string, error) {
	activePods := balancer.pods.GetActivePodsIp()

	if len(activePods) == 0 {
		return "", errors.New("NO PODS AVAILABLE")
	}

	return balancer.strategy.Select(activePods)
}

func (balancer *Balancer) GetPods() *k8s.PodsMap {
	return &balancer.pods
}
