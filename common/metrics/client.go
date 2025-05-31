package metrics

type Client interface {
	CountCall(name string, status string)
	RecordTime(name string, value float64)
}
