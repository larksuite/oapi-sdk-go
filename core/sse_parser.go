package larkcore

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

// SSEEvent is a protocol event. Data is never decoded as business JSON.
type SSEEvent struct {
	Event string
	Data  string
	ID    *string
}

type sseParser struct {
	onControl func(string, *int64)
	reader    *bufio.Reader
	limit     int
	size      int
	first     bool
	skipLF    bool
	data      []string
	event     string
	id        *string
	cursor    string
	retry     *int64
}

func newSSEParser(reader io.Reader, limit int, cursor string) *sseParser {
	return &sseParser{reader: bufio.NewReaderSize(reader, 4096), limit: limit, first: true, cursor: cursor}
}

func (p *sseParser) readLine() (string, error) {
	var line []byte
	for {
		b, err := p.reader.ReadByte()
		if err != nil {
			return "", err
		}
		if p.skipLF {
			p.skipLF = false
			if b == '\n' {
				continue
			}
		}
		p.size++
		if p.size > p.limit {
			return "", errors.New("SSE event exceeds maxEventBytes")
		}
		if b == '\r' || b == '\n' {
			p.skipLF = b == '\r'
			value := strings.ToValidUTF8(string(line), "\uFFFD")
			if p.first {
				p.first = false
				value = strings.TrimPrefix(value, "\uFEFF")
			}
			return value, nil
		}
		line = append(line, b)
	}
}

func (p *sseParser) consume(line string) *SSEEvent {
	if line == "" {
		return p.finish()
	}
	if strings.HasPrefix(line, ":") {
		return nil
	}
	parts := strings.SplitN(line, ":", 2)
	value := ""
	if len(parts) == 2 {
		value = strings.TrimPrefix(parts[1], " ")
	}
	switch parts[0] {
	case "data":
		p.data = append(p.data, value)
	case "event":
		p.event = value
	case "id":
		if !strings.ContainsRune(value, 0) {
			p.id = &value
		}
	case "retry":
		if value != "" && strings.Trim(value, "0123456789") == "" {
			if n, err := strconv.ParseInt(value, 10, 64); err == nil {
				p.retry = &n
			}
		}
	}
	if p.onControl != nil {
		p.onControl(p.cursor, p.retry)
	}
	return nil
}

func (p *sseParser) finish() *SSEEvent {
	if p.id != nil {
		p.cursor = *p.id
	}
	if p.onControl != nil {
		p.onControl(p.cursor, p.retry)
	}
	var result *SSEEvent
	if len(p.data) > 0 {
		name := p.event
		if name == "" {
			name = "message"
		}
		result = &SSEEvent{Event: name, Data: strings.Join(p.data, "\n"), ID: p.id}
	}
	p.data = nil
	p.event = ""
	p.id = nil
	p.size = 0
	return result
}

func (p *sseParser) next() (*SSEEvent, error) {
	for {
		line, err := p.readLine()
		if err != nil {
			return nil, err
		}
		if event := p.consume(line); event != nil {
			return event, nil
		}
	}
}
