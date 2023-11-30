package plugin

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/gorilla/mux"
)

func (ds *Datasource) resourceRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/namespaces", ds.listNamespaces).Methods(http.MethodGet)
	router.HandleFunc("/api/metricNames", ds.listMetrics).Methods(http.MethodGet)
	router.HandleFunc("/api/dimensions", ds.listDimensions).Methods(http.MethodGet)
	return router
}

type Namespace struct {
	Name string `json:"name"`
}

type MetricName struct {
	Name string `json:"name"`
}

type Dimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (ds *Datasource) listNamespaces(w http.ResponseWriter, r *http.Request) {
	metrics, err := ds.cloudwatch.ListMetrics(r.Context(), &cloudwatch.ListMetricsInput{})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m := make(map[string]struct{}, len(metrics.Metrics))
	namespaces := make([]Namespace, 0, len(metrics.Metrics))
	for _, metric := range metrics.Metrics {
		if _, has := m[*metric.Namespace]; !has {
			m[*metric.Namespace] = struct{}{}
			namespaces = append(namespaces, Namespace{Name: *metric.Namespace})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(namespaces)
}

func (ds *Datasource) listMetrics(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	metrics, err := ds.cloudwatch.ListMetrics(r.Context(), &cloudwatch.ListMetricsInput{
		Namespace: aws.String(namespace),
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m := make(map[string]struct{}, len(metrics.Metrics))
	metricNames := make([]MetricName, 0, len(metrics.Metrics))
	for _, metric := range metrics.Metrics {
		if _, has := m[*metric.MetricName]; !has {
			m[*metric.MetricName] = struct{}{}
			metricNames = append(metricNames, MetricName{Name: *metric.MetricName})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(metricNames)
}

func (ds *Datasource) listDimensions(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	metricName := r.URL.Query().Get("metricName")
	metrics, err := ds.cloudwatch.ListMetrics(r.Context(), &cloudwatch.ListMetricsInput{
		Namespace:  aws.String(namespace),
		MetricName: aws.String(metricName),
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(metrics.Metrics) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}
	dimensions := make([]Dimension, 0, len(metrics.Metrics[0].Dimensions))
	for _, dimension := range metrics.Metrics[0].Dimensions {
		dimensions = append(dimensions, Dimension{
			Name:  *dimension.Name,
			Value: *dimension.Value,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dimensions)
}
