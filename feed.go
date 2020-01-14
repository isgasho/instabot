package instabot

import (
	"context"
	"fmt"

	"github.com/winterssy/sreq"
)

const (
	apiGetUserFeed = apiV1 + "feed/user/%s/?max_id=%s"
	apiGetTimeline = apiV1 + "feed/timeline/"
)

func (bot *Bot) GetUserFeed(ctx context.Context, userId string, maxId string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetUserFeed, userId, maxId),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetSelfFeed(ctx context.Context, maxId string) (sreq.H, error) {
	return bot.GetUserFeed(ctx, bot.GetUserId(), maxId)
}

func (bot *Bot) GetTimeline(ctx context.Context, maxId string) (sreq.H, error) {
	form := sreq.Form{
		"phone_id":            bot.phoneId,
		"max_id":              maxId,
		"timezone_offset":     bot.timeOffset,
		"_csrftoken":          bot.GetCSRFToken(),
		"device_id":           bot.igDeviceId,
		"request_id":          GenerateUUID(),
		"_uuid":               bot.uuid,
		"session_id":          GenerateUUID(),
		"bloks_versioning_id": bloksVersionId,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, apiGetTimeline,
		sreq.WithForm(form),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}
