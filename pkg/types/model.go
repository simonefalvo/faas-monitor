package types

type Function struct {
	Name           string           `json:"name"`
	Replicas       int              `json:"replicas"`
	ResponseTime   float64          `json:"response_time"`
	ProcessingTime float64          `json:"processing_time"`
	Throughput     float64          `json:"throughput"`
	ColdStart      float64          `json:"cold_start"`
	Cpu            map[string]int64 `json:"cpu"`
	Mem            map[string]int64 `json:"mem"`
}

type Node struct {
	Name string  `json:"name"`
	Cpu  float64 `json:"cpu"`
	Mem  float64 `json:"mem"`
}

type Message struct {
	Functions []Function `json:"functions"`
	Nodes     []Node     `json:"nodes"`
	Timestamp int64      `json:"timestamp"`
}