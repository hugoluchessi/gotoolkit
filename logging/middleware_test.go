package logging

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hugoluchessi/gotoolkit/tctx"
)

func TestNewContextLoggerMiddleware(t *testing.T) {
	l := NewMockLogger()
	ctxl := NewContextLogger(l)

	mw := NewContextLoggerMiddleware(ctxl)

	if mw == nil {
		t.Error("NewContextLoggerMiddleware cannot return nil.")
	}
}

func TestContextLoggerHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	h := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("some", "test")
	})

	l := NewMockLogger()
	ctxl := NewContextLogger(l)

	tid := uuid.New()
	tms := time.Now()
	ctx := req.Context()
	ctx = tctx.Create(ctx, tid, tms)

	tctx.AddRequestHeaders(ctx, req)

	mw := NewContextLoggerMiddleware(ctxl)
	mwh := mw.Handler(h)

	mwh.ServeHTTP(res, req)

	loggedContent := l.String()

	loggedContents := strings.Split(loggedContent, "\n")

	if !strings.Contains(loggedContents[0], requestStartedMsg) {
		t.Error("Missing request log part 'requestStartedMsg'.")
	}

	if !strings.Contains(loggedContents[0], requestHostLogEntry) {
		t.Error("Missing request log part 'requestHostLogEntry'.")
	}

	if !strings.Contains(loggedContents[0], requestRemoteAddrLogEntry) {
		t.Error("Missing request log part 'requestRemoteAddrLogEntry'.")
	}

	if !strings.Contains(loggedContents[0], requestMethodLogEntry) {
		t.Error("Missing request log part 'requestMethodLogEntry'.")
	}

	if !strings.Contains(loggedContents[0], requestURILogEntry) {
		t.Error("Missing request log part 'requestURILogEntry'.")
	}

	if !strings.Contains(loggedContents[0], requestProtoLogEntry) {
		t.Error("Missing request log part 'requestProtoLogEntry'.")
	}

	if !strings.Contains(loggedContents[0], requestUserAgentLogEntry) {
		t.Error("Missing request log part 'requestUserAgentLogEntry'.")
	}

	if !strings.Contains(loggedContents[1], requestEndedMsg) {
		t.Error("Missing request log part 'requestEndedMsg'.")
	}

	if !strings.Contains(loggedContents[1], requestDurationMilliseconds) {
		t.Error("Missing request log part 'requestDurationMilliseconds'.")
	}
}

func TestContextLoggerHandlerWithoutContext(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	h := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("some", "test")
	})

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	cfg := LoggerConfig{w}
	cfgs := []LoggerConfig{cfg}
	l := NewZapLogger(cfgs)
	ctxl := NewContextLogger(l)

	mw := NewContextLoggerMiddleware(ctxl)
	mwh := mw.Handler(h)

	mwh.ServeHTTP(res, req)

	response := res.Result()

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Request without transaction headers must fail")
	}
}