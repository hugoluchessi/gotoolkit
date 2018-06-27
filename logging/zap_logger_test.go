package logging

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestNewZapLogger(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	cfg := LoggerConfig{w}
	cfgs := []LoggerConfig{cfg}

	l := NewZapLogger(cfgs)

	if l == nil {
		t.Error("[l] must not be nil.")
	}
}

func createZapLoggingContext() (*bytes.Buffer, *ZapLogger) {
	var b bytes.Buffer
	cfg := LoggerConfig{&b}
	cfgs := []LoggerConfig{cfg}

	l := NewZapLogger(cfgs)

	return &b, l
}

func TestZapLog(t *testing.T) {
	b, l := createZapLoggingContext()

	l.Log(
		"Message test",
		map[string]interface{}{
			"test1": "value1",
			"test2": "value2",
		},
	)

	content := b.String()

	ell := "\"level\":\"info\""
	etest1 := "\"test1\":\"value1\""
	etest2 := "\"test2\":\"value2\""

	if !strings.Contains(content, ell) {
		t.Error("Wrong log level.")
	}

	if !strings.Contains(content, etest1) {
		t.Error("Wrong test map params 1.")
	}

	if !strings.Contains(content, etest2) {
		t.Error("Wrong test map params 2.")
	}
}

func TestZapError(t *testing.T) {
	b, l := createZapLoggingContext()

	l.Error(
		"Message test",
		map[string]interface{}{
			"test1": "value1",
			"test2": "value2",
		},
	)

	content := b.String()

	ell := "\"level\":\"error\""
	etest1 := "\"test1\":\"value1\""
	etest2 := "\"test2\":\"value2\""

	if !strings.Contains(content, ell) {
		t.Error("Wrong log level.")
	}

	if !strings.Contains(content, etest1) {
		t.Error("Wrong test map params 1.")
	}

	if !strings.Contains(content, etest2) {
		t.Error("Wrong test map params 2.")
	}
}
