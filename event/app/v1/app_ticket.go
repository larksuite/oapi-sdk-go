package v1

import (
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/store"
	"github.com/larksuite/oapi-sdk-go/event"
	"github.com/larksuite/oapi-sdk-go/event/core/model"
	"math"
)

type AppTicketEventData struct {
	*model.BaseEventData
	AppTicket string `json:"app_ticket"`
}

type AppTicketEvent struct {
	*model.BaseEvent
	Event *AppTicketEventData `json:"event"`
}

type AppTicketHandler struct {
	event *AppTicketEvent
}

func (h *AppTicketHandler) GetEvent() interface{} {
	h.event = &AppTicketEvent{}
	return h.event
}

func (h *AppTicketHandler) Handle(ctx *core.Context, event interface{}) error {
	appTicketEvent := event.(*AppTicketEvent)
	conf := config.ByCtx(ctx)
	return conf.GetStore().Put(ctx, store.AppTicketKey(appTicketEvent.Event.AppID), appTicketEvent.Event.AppTicket, math.MaxInt32)
}

func SetAppTicketHandler(conf *config.Config) {
	if conf.GetAppSettings().AppType == constants.AppTypeInternal {
		return
	}
	event.SetTypeHandler(conf, "app_ticket", &AppTicketHandler{})
}
