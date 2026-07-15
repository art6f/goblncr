package balancer

import (
	"context"
	"fmt"
	"log/slog"
	"net/netip"

	"github.com/art6f/goblncr/config"
	"github.com/art6f/goblncr/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Balancer struct {
	Client *kubernetes.Clientset
	Config *config.AppConfig
	pods   PodsMap
}

func NewBalancer(config *config.AppConfig) *Balancer {
	return &Balancer{
		Client: k8s.NewClient(),
		Config: config,
		pods:   make(map[string]netip.Addr, 0),
	}
}

func (balancer *Balancer) Run() {
	slog.Info("Starting Balancer...")
	slog.Info("Querying pods", "namespace", balancer.Config.Namespace, "label", balancer.Config.PodsSelector)

	// balancer.config.Namespace
	pods, err := balancer.Client.CoreV1().Pods(balancer.Config.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: balancer.Config.PodsSelector,
	})

	if err != nil {
		slog.Error("Errror retrieving pods", "errro", err.Error())
		panic(err.Error())
	}

	if len(pods.Items) == 0 {
		slog.Error("No pods found")
	} else {
		slog.Info(fmt.Sprintf("Pods found: %v", len(pods.Items)))
		podsInfo := ""
		
		for _, pod := range pods.Items {
			podsInfo += fmt.Sprintf("\n\t%s (%s)", pod.Status.PodIP, pod.Name)
		
			if ip, err := netip.ParseAddr(pod.Status.PodIP); err == nil {
				balancer.pods[pod.Status.PodIP] = ip
			}
		}
		slog.Info("Pods IPs to serve: " + podsInfo)
	}
	
	InformerWatchPods(balancer)
}