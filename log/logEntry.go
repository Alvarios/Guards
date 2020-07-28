package log

import (
	"time"
)

type LogEntry struct {
	ReceivedTime      time.Time
	RequestMethod     string
	RequestURL        string
	RequestHeaderSize int64
	RequestBodySize   int64
	UserAgent         string
	Referer           string
	Proto             string

	RemoteIP string
	ServerIP string

	Status             int
	ResponseHeaderSize int64
	ResponseBodySize   int64
	Latency            time.Duration
}
