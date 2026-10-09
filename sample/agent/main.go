// Agent v1 scenarios. See README.md for credentials and run commands.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	core "github.com/larksuite/oapi-sdk-go/v3/core"
)

type example struct {
	client *lark.Client
	token  string
}

func required(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" || strings.ContainsAny(value, "\r\n\x00") {
		return "", fmt.Errorf("set a valid %s environment variable", name)
	}
	return value, nil
}

func newExample() (*example, error) {
	appID, err := required("APP_ID")
	if err != nil {
		return nil, err
	}
	secret, err := required("APP_SECRET")
	if err != nil {
		return nil, err
	}
	token, err := required("USER_ACCESS_TOKEN")
	if err != nil {
		return nil, err
	}
	base := os.Getenv("OPEN_BASE_URL")
	if base == "" {
		base = "https://open.feishu.cn"
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return nil, errors.New(
			"OPEN_BASE_URL must be HTTPS (HTTP allowed only for localhost testing)",
		)
	}
	localHost := parsed.Hostname() == "127.0.0.1" || parsed.Hostname() == "localhost"
	localHTTP := parsed.Scheme == "http" && localHost
	invalidEndpoint := parsed.Hostname() == "" || parsed.User != nil ||
		parsed.RawQuery != "" || parsed.Fragment != ""
	if invalidEndpoint || (parsed.Scheme != "https" && !localHTTP) {
		return nil, errors.New(
			"OPEN_BASE_URL must be HTTPS (HTTP allowed only for localhost testing)",
		)
	}
	client := lark.NewClient(
		appID,
		secret,
		lark.WithOpenBaseUrl(base),
		lark.WithReqTimeout(10*time.Second),
		lark.WithLogLevel(core.LogLevelError),
	)
	return &example{client: client, token: token}, nil
}

func (e *example) run(ctx context.Context, mode string) error {
	switch mode {
	case "chat":
		return e.chat(ctx)
	case "upload":
		return e.upload(ctx)
	case "resume", "close":
		sid, err := required("AGENT_SESSION_ID")
		if err != nil {
			return err
		}
		queryID, err := required("AGENT_QUERY_MESSAGE_ID")
		if err != nil {
			return err
		}
		cursor := ""
		if mode == "resume" {
			cursor, err = required("AGENT_LAST_EVENT_ID")
			if err != nil {
				return err
			}
		}
		return e.read(ctx, sid, queryID, cursor, mode == "close")
	case "interrupt":
		return e.interrupt(ctx)
	default:
		return errors.New("usage: go run ./agent chat|upload|resume|close|interrupt")
	}
}

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		log.Print("usage: go run ./agent chat|upload|resume|close|interrupt")
		os.Exit(1)
	}
	example, err := newExample()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	ctx, cancel := context.WithTimeout(signalCtx, 180*time.Second)
	err = example.run(ctx, os.Args[1])
	cancel()
	stop()
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Print(err)
		os.Exit(1)
	}
}
