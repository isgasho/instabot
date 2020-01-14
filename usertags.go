package instabot

import (
	"context"
	"fmt"

	"github.com/winterssy/sreq"
)

const (
	apiGetUserTags = apiV1 + "usertags/%s/feed/"
)

func (bot *Bot) GetUserTags(ctx context.Context, userId string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetUserTags, userId),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetSelfTags(ctx context.Context) (sreq.H, error) {
	return bot.GetUserTags(ctx, bot.GetUserId())
}
