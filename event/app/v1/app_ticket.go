package v1

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/store"
	"github.com/larksuite/oapi-sdk-go/event/core/handlers"
	"github.com/larksuite/oapi-sdk-go/event/core/model"

	"time"
)

type AppTicketEventData struct {
	*model.BaseEventData
	AppTicket string `json:"app_ticket"`
}

type AppTicketEvent struct {
	*model.BaseEvent
	Event *AppTicketEventData `json:"event"`
}

type AppTicketEventHandler struct {
	event *AppTicketEvent
}

func (h *AppTicketEventHandler) GetEvent() interface{} {
	h.event = &AppTicketEvent{}
	return h.event
}

func (h *AppTicketEventHandler) Handle(ctx *core.Context, event interface{}) error {
	appTicketEvent := event.(*AppTicketEvent)
	conf := config.ByCtx(ctx)
	return conf.GetStore().Put(ctx, store.AppTicketKey(appTicketEvent.Event.AppID), appTicketEvent.Event.AppTicket, time.Hour*12)
}

func SetAppTicketEventHandler(conf *config.Config) {
	if conf.GetAppSettings().AppType == constants.AppTypeInternal {
		return
	}
	handlers.SetTypeHandler(conf, "app_ticket", &AppTicketEventHandler{})
}
