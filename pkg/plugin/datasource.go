package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ backend.CallResourceHandler   = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

func NewDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	cfg, err := GetSettings(settings)
	if err != nil {
		return nil, err
	}
	rt, err := newSigningRoundTripper(ctx, cfg)
	if err != nil {
		return nil, err
	}
	c, err := newConfig(ctx,
		config.WithRegion(cfg.MonitoringRegion),
		config.WithHTTPClient(&http.Client{Transport: rt}),
	)
	if err != nil {
		return nil, err
	}
	cwClient := newCloudWatch(c)
	if err != nil {
		return nil, err
	}

	return &Datasource{
		cloudwatch: cwClient,
	}, nil
}

func (ds *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return ds.resourceHandler().CallResource(ctx, req, sender)
}

func (ds *Datasource) resourceHandler() backend.CallResourceHandler {
	ds.rcHandlerOnce.Do(func() {
		ds.rcHandler = httpadapter.New(ds.resourceRouter())
	})
	return ds.rcHandler
}

type TokenIssuer interface {
	IssueAccessToken(ctx context.Context) (string, error)
}

type Datasource struct {
	cloudwatch    *cloudwatch.Client
	rcHandlerOnce sync.Once
	rcHandler     backend.CallResourceHandler
}

func (ds *Datasource) Dispose() {}

func (ds *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		res := ds.query(ctx, req.PluginContext, q)
		response.Responses[q.RefID] = res
	}
	return response, nil
}

type queryModel struct {
	Namespace  string           `json:"namespace"`
	MetricName string           `json:"metricName"`
	Statistic  string           `json:"statistic"`
	Dimensions []queryDimension `json:"dimensions"`
}

type queryDimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (q *queryModel) isDefaultQuery() bool {
	return q.Namespace == "" || q.MetricName == "" || q.Statistic == ""
}

func (ds *Datasource) query(ctx context.Context, _ backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	if qm.isDefaultQuery() {
		return response
	}

	dims := make([]types.Dimension, 0, len(qm.Dimensions))
	for _, dim := range qm.Dimensions {
		dims = append(dims, types.Dimension{
			Name:  aws.String(dim.Name),
			Value: aws.String(dim.Value),
		})
	}
	metrics, err := ds.cloudwatch.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
		StartTime: aws.Time(query.TimeRange.From),
		EndTime:   aws.Time(query.TimeRange.To),
		MetricDataQueries: []types.MetricDataQuery{
			{
				Id: aws.String("metric1"),
				MetricStat: &types.MetricStat{
					Period: aws.Int32(60),
					Stat:   aws.String(qm.Statistic),
					Metric: &types.Metric{
						Dimensions: dims,
						MetricName: aws.String(qm.MetricName),
						Namespace:  aws.String(qm.Namespace),
					},
				},
			},
		},
	})
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("cloudwatch get metric data: %v", err.Error()))
	}

	frame := data.NewFrame("response")
	times := make([]time.Time, 0, len(metrics.MetricDataResults))
	values := make([]float64, 0, len(metrics.MetricDataResults))
	for _, metric := range metrics.MetricDataResults {
		for _, timestamp := range metric.Timestamps {
			times = append(times, timestamp)
		}
		for _, value := range metric.Values {
			values = append(values, value)
		}
	}
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, times),
		data.NewField("values", nil, values),
	)
	response.Frames = append(response.Frames, frame)
	return response
}

func (ds *Datasource) CheckHealth(_ context.Context, _ *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	const status = backend.HealthStatusOk
	const message = "Success"

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
