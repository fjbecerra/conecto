package _http

import (
	"regexp"
	"github.com/tidwall/gjson"
)


type CursorExtractor interface {
    Extract(resp *HttpResponse) (*PageCursor, bool)
}

type JSONCursorExtractor struct {
   Path string
}

func (e *JSONCursorExtractor) Extract(resp *HttpResponse,) (*PageCursor, bool) {

	next := gjson.GetBytes(resp.Body,e.Path,
	).String()

	if next == "" {
		return nil, false
	}

	return &PageCursor{
		Value: next,
	}, true

}

type LinkHeaderExtractor struct{}

func (e *LinkHeaderExtractor) Extract(resp *HttpResponse,) (*PageCursor, bool) {

	link := resp.Headers.Get("Link")

	if link == "" {
		return nil, false
	}

	re := regexp.MustCompile(`page_info=([^&>]+)`)

	match := re.FindStringSubmatch(link)

	if len(match) < 2 {
		return nil, false
	}

	return &PageCursor{
		Value: match[1],
	}, true

}