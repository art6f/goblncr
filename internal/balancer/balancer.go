package balancer

import (
	"context"
	"log/slog"

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
		pods:   make(PodsMap, 0),
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

	WatchPods(balancer)
}

func (balancer *Balancer) GetActivePodsIp() []string {
	var podList []string
	for _, podData := range balancer.pods {
		if podData.Ready {
			podList = append(podList, podData.Ip.String())
		}
	}
	return podList
}
