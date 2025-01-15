package config

import (
	"bytes"
	"os"
	"testing"
)

func TestLoggerFlags(t *testing.T) {
	logger := &Logger{}
	flags := logger.Flags()

	if len(flags) != 3 {
		t.Errorf("expected 3 flags, got %d", len(flags))
	}

	if flags[0].Names()[0] != "log-level" {
		t.Errorf("expected 'log-level', got '%s'", flags[0].Names()[0])
	}
	if flags[1].Names()[0] != "log-format" {
		t.Errorf("expected 'log-format', got '%s'", flags[1].Names()[0])
	}
	if flags[2].Names()[0] != "log-output" {
		t.Errorf("expected 'log-output', got '%s'", flags[2].Names()[0])
	}
}

func TestLoggerNew(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test.log")
	if err != nil {
		t.Fatalf("expected no error creating temp file, got %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tests := []struct {
		level  string
		format string
		output string
		valid  bool
	}{
		{"debug", "json", "stdout", true},
		{"info", "text", "stderr", true},
		{"warn", "json", tmpfile.Name(), true},
		{"error", "text", "stdout", true},
		{"invalid", "json", "stdout", false},
		{"info", "invalid", "stdout", false},
	}

	for _, tt := range tests {
		t.Run(tt.level+"_"+tt.format+"_"+tt.output, func(t *testing.T) {
			logger := &Logger{
				level:  tt.level,
				format: tt.format,
				output: tt.output,
			}

			log, closer, err := logger.New()
			if closer != nil {
				defer closer()
			}

			if tt.valid {
				if log == nil {
					t.Error("expected logger to be not nil")
				}
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if log != nil {
					t.Error("expected logger to be nil")
				}
				if err == nil {
					t.Error("expected an error")
				}
			}
		})
	}
}

func TestLoggerOutputToFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test.log")
	if err != nil {
		t.Fatalf("expected no error creating temp file, got %v", err)
	}
	defer os.Remove(tmpfile.Name())

	logger := &Logger{
		level:  "info",
		format: "text",
		output: tmpfile.Name(),
	}

	log, closer, err := logger.New()
	defer closer()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if log == nil {
		t.Fatal("expected logger to be not nil")
	}

	log.Info("test message")

	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("expected no error reading file, got %v", err)
	}
	if !bytes.Contains(data, []byte("test message")) {
		t.Errorf("expected file to contain 'test message'")
	}
}
