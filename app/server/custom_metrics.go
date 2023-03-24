package server

import "encoding/json"

type CustomMetric struct {
	Name          string  `json:"name"`
	Value         float64 `json:"value"`
	Unit          string  `json:"unit"`
	AppGUID       string  `json:"app_guid"`
	InstanceIndex uint32  `json:"instance_index"`
}

type MetricsConsumer struct {
	InstanceIndex uint32          `json:"instance_index"`
	CustomMetrics []*CustomMetric `json:"metrics"`
}

func createCustomMetricsPayload(appId string, metricsValue float64) []byte {
	customMetrics := []*CustomMetric{
		{
			Name:          "tooManyRequestCustomMetrics",
			Value:         metricsValue,
			Unit:          "test-unit",
			AppGUID:       appId,
			InstanceIndex: 0,
		},
	}
	metricsValueBodyJson, _ := json.Marshal(MetricsConsumer{
		InstanceIndex: 0,
		CustomMetrics: customMetrics,
	})
	return metricsValueBodyJson
}
