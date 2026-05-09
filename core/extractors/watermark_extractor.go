package extractors

import (
	"conecto/core"
	"time"
	"github.com/tidwall/gjson"
)

type WatermarkExtractor interface {
	Extract(ev core.Event) time.Time
}

type JsonWatermarkExtractor struct {
	Path string
}

func (e *JsonWatermarkExtractor) Extract(ev core.Event) time.Time {
	return gjson.GetBytes(ev.Payload, e.Path).Time()
}