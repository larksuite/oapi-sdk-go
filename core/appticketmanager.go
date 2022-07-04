package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var appTicketManager AppTicketManager = AppTicketManager{cache: cache}

func GetAppTicketManager() *AppTicketManager {
	return &appTicketManager
}

type AppTicketManager struct {
	cache Cache
}

func (m *AppTicketManager) Get(ctx context.Context, config *Config) (string, error) {
	ticket, err := m.cache.Get(ctx, appTicketKey(config.AppId))
	if err != nil {
		return "", err
	}
	if ticket == "" {
		applyAppTicket(ctx, config)
	}
	return ticket, nil
}

func (m *AppTicketManager) Set(ctx context.Context, appId, value string, ttl time.Duration) error {
	return m.cache.Set(ctx, appTicketKey(appId), value, ttl)
}

func appTicketKey(appID string) string {
	return fmt.Sprintf("%s-%s", appTicketKeyPrefix, appID)
}

func applyAppTicket(ctx context.Context, config *Config) {
	rawResp, err := SendRequest(ctx, config, http.MethodPost, applyAppTicketPath, []AccessTokenType{accessTokenTypeNone}, &applyAppTicketReq{
		AppID:     config.AppId,
		AppSecret: config.AppSecret,
	})
	if err != nil {
		config.Logger.Error(ctx, fmt.Sprintf("apply app_ticket, error: %v", err))
		return
	}
	if !strings.Contains(rawResp.Header.Get(contentTypeHeader), contentTypeJson) {
		config.Logger.Error(ctx, fmt.Sprintf("apply app_ticket, response content-type not json, response: %v", rawResp))
		return
	}
	codeError := &CodeError{}
	err = json.Unmarshal(rawResp.RawBody, codeError)
	if err != nil {
		config.Logger.Error(ctx, fmt.Sprintf("apply app_ticket, json unmarshal error: %v", err))
		return
	}
	if codeError.Code != 0 {
		config.Logger.Error(ctx, fmt.Sprintf("apply app_ticket, response error: %+v", codeError))
		return
	}
}
