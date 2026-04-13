package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Logger menulis log ke file (plain text terformat) dan stdout secara bersamaan.
type Logger struct {
	appName string
	version string
	env     string
	rotator *dailyRotator
}

// dailyRotator mengelola file log harian.
type dailyRotator struct {
	mu      sync.Mutex
	logDir  string
	appName string
	date    string
	file    *os.File
}

func (r *dailyRotator) writeString(s string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != r.date {
		if err := r.rotate(today); err != nil {
			return err
		}
	}
	_, err := r.file.WriteString(s)
	return err
}

func (r *dailyRotator) rotate(today string) error {
	if r.file != nil {
		_ = r.file.Close()
	}
	filename := filepath.Join(r.logDir, fmt.Sprintf("%s-%s.log", r.appName, today))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", filename, err)
	}
	r.file = file
	r.date = today
	return nil
}

// separator line
const line = "-------------------------------------------------------------------------------"

// New membuat Logger baru dengan daily rotation.
func New(logDir, appName, version, env string) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	rotator := &dailyRotator{logDir: logDir, appName: appName}
	today := time.Now().Format("2006-01-02")
	if err := rotator.rotate(today); err != nil {
		return nil, err
	}

	return &Logger{
		appName: appName,
		version: version,
		env:     env,
		rotator: rotator,
	}, nil
}

// PrintBanner mencetak info startup ke stdout dan file log.
func (l *Logger) PrintBanner(dbStatus, port string) {
	now := time.Now().Format("2006-01-02 15:04:05 MST")

	block := fmt.Sprintf(
		"\n%s\nService Name : %s\nVersion      : %s\nDate         : %s\nDatabase     : %s\nPort         : %s\nEnvironment  : %s\n%s\n",
		line,
		l.appName,
		l.version,
		now,
		dbStatus,
		port,
		l.env,
		line,
	)

	fmt.Print(block)
	_ = l.rotator.writeString(block)
}

// ── Level helpers ──────────────────────────────────────────────────────────────

func (l *Logger) Info(message string, fields ...Field) {
	l.write("INFO", message, fields...)
}

func (l *Logger) Error(message string, fields ...Field) {
	l.write("ERROR", message, fields...)
}

func (l *Logger) Fatal(message string, fields ...Field) {
	l.write("FATAL", message, fields...)
	os.Exit(1)
}

func (l *Logger) Warn(message string, fields ...Field) {
	l.write("WARN", message, fields...)
}

func (l *Logger) Debug(message string, fields ...Field) {
	if l.env == "production" {
		return
	}
	l.write("DEBUG", message, fields...)
}

// ── Field builder ──────────────────────────────────────────────────────────────

type Field struct {
	Key   string
	Value interface{}
}

func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// ── HTTP request/response log ──────────────────────────────────────────────────

type HTTPLogEntry struct {
	// Request
	RequestID string
	Method    string
	Path      string
	UserAgent string
	ReqBody   string

	// Response
	RespBody   string
	StatusCode int
	LatencyMS  int64
	Latency    string
	Message    string
	Errors     string
}

// WriteHTTPLog menulis log request+response dalam format yang rapi dan mudah dibaca.
func (l *Logger) WriteHTTPLog(e HTTPLogEntry) {
	now := time.Now().Format("2006-01-02 15:04:05 MST")

	level := "INFO"
	switch {
	case e.StatusCode >= 500:
		level = "ERROR"
	case e.StatusCode >= 400:
		level = "WARN"
	}

	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%-14s %s\n", "Time", ":  "+now))
	sb.WriteString(fmt.Sprintf("%-14s %s\n", "Level", ":  "+level))
	sb.WriteString("\n")

	// ── Request ───────────────────────────────────────────────────────────────
	sb.WriteString("Request\n")
	sb.WriteString(fmt.Sprintf("  %-14s %s\n", "ID", ":  "+e.RequestID))
	sb.WriteString(fmt.Sprintf("  %-14s %s\n", "Method", ":  "+e.Method))
	sb.WriteString(fmt.Sprintf("  %-14s %s\n", "Path", ":  "+e.Path))
	sb.WriteString(fmt.Sprintf("  %-14s %s\n", "User Agent", ":  "+e.UserAgent))

	if e.ReqBody != "" {
		sb.WriteString(fmt.Sprintf("  %-14s\n", "Body"))
		sb.WriteString(indentBlock(e.ReqBody, "    "))
	} else {
		sb.WriteString(fmt.Sprintf("  %-14s %s\n", "Body", ":  -"))
	}

	sb.WriteString("\n")

	// ── Response ──────────────────────────────────────────────────────────────
	sb.WriteString("Response\n")
	sb.WriteString(fmt.Sprintf("  %-14s %s\n", "ID", ":  "+e.RequestID))
	sb.WriteString(fmt.Sprintf("  %-14s %d\n", "Status", e.StatusCode))

	if e.RespBody != "" {
		sb.WriteString(fmt.Sprintf("  %-14s\n", "Body"))
		sb.WriteString(indentBlock(e.RespBody, "    "))
	} else {
		sb.WriteString(fmt.Sprintf("  %-14s %s\n", "Body", ":  -"))
	}

	sb.WriteString(fmt.Sprintf("  %-14s %s\n", "Time", ":  "+now))
	sb.WriteString(fmt.Sprintf("  %-14s %s\n", "Message", ":  "+e.Message))
	sb.WriteString(fmt.Sprintf("  %-14s %d ms\n", "Latency MS", e.LatencyMS))
	sb.WriteString(fmt.Sprintf("  %-14s %s\n", "Latency Human", ":  "+e.Latency))

	if e.Errors != "" {
		sb.WriteString(fmt.Sprintf("  %-14s %s\n", "Errors", ":  "+e.Errors))
	}

	sb.WriteString(line + "\n")

	output := sb.String()
	fmt.Print(output)
	_ = l.rotator.writeString(output)
}

// ── Internal ───────────────────────────────────────────────────────────────────

func (l *Logger) write(level, message string, fields ...Field) {
	now := time.Now().Format("2006-01-02 15:04:05 MST")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] [%s] %s", now, level, message))

	for _, f := range fields {
		sb.WriteString(fmt.Sprintf(" | %s=%v", f.Key, f.Value))
	}
	sb.WriteString("\n")

	output := sb.String()
	fmt.Print(output)
	_ = l.rotator.writeString(output)
}

// indentBlock mem-pretty-print JSON atau menambahkan indentasi pada plain text.
func indentBlock(body, prefix string) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(body), prefix, "  "); err == nil {
		return prefix + buf.String() + "\n"
	}
	// Bukan JSON — tampilkan apa adanya
	return prefix + body + "\n"
}
