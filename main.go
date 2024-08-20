package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	gatewayapiGatewayclassInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gatewayapi_gatewayclass_info",
			Help: "Information about a GatewayClass",
		},
		[]string{"name"},
	)
)

func init() {
	prometheus.MustRegister(gatewayapiGatewayclassInfo)
}

func main() {
	// Get optional environment variables
	apiserver := os.Getenv("APISERVER")
	kubeconfig := os.Getenv("KUBECONFIG")

	// Set up Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
		if err != nil {
			fmt.Printf("Error creating Kubernetes client. Is APISERVER and/or KUBECONFIG set?\n(%v)\n", err)
			os.Exit(1)
		}
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes dynamic client\n(%v)\n", err)
		os.Exit(1)
	}

	// Set up Prometheus endpoint
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		collectMetrics(dynClient)
		promhttp.Handler().ServeHTTP(w, r)
	})
	server := &http.Server{Addr: ":8080"}

	// Handle shutdown gracefully
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Shutting down server...")
		server.Close()
	}()

	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func collectMetrics(dynClient *dynamic.DynamicClient) {
	// Clear previous metric values
	gatewayapiGatewayclassInfo.Reset()

	// Define the GroupVersionResource for the GatewayClasses CRD
	gvr := schema.GroupVersionResource{
		Group:    "gateway.networking.k8s.io",
		Version:  "v1",
		Resource: "gatewayclasses",
	}

	// List GatewayClasses
	gatewayClasses, err := dynClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error fetching GatewayClasses: %v\n", err)
		return
	}

	for _, gatewayClass := range gatewayClasses.Items {
		gatewayapiGatewayclassInfo.WithLabelValues(gatewayClass.GetName()).Set(1)
	}
}
