package balancer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/art6f/goblncr/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Balancer struct {
	client *kubernetes.Clientset
}

func NewBalancer() *Balancer {
	return &Balancer{
		client: k8s.NewClient(),
	}
}

func (balancer *Balancer) Run() {
	slog.Info("Starting Balancer")

	pods, err := balancer.client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		slog.Error("Errror retrieving pods", "errro", err.Error())
		panic(err.Error())
	}

	slog.Info(fmt.Sprintf("Pods found: %v", len(pods.Items)))
}
