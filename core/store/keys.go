package store

import "fmt"

const (
	appTicketKeyPrefix         = "app_ticket"
	appAccessTokenKeyPrefix    = "app_access_token"
	tenantAccessTokenKeyPrefix = "tenant_access_token"
)

func AppTicketKey(appID string) string {
	return fmt.Sprintf("%s-%s", appTicketKeyPrefix, appID)
}

func AppAccessTokenKey(appID string) string {
	return fmt.Sprintf("%s-%s", appAccessTokenKeyPrefix, appID)
}

func TenantAccessTokenKey(appID, tenantKey string) string {
	return fmt.Sprintf("%s-%s-%s", tenantAccessTokenKeyPrefix, appID, tenantKey)
}
