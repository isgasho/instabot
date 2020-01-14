package instabot

import (
	"context"
	"fmt"

	"github.com/winterssy/sreq"
)

const (
	apiGetUserByName = apiV1 + "users/%s/usernameinfo/"
	apiGetUserById   = apiV1 + "users/%s/info/"
	apiSearchUser    = apiV1 + "users/search/"
)

func (bot *Bot) GetUserByName(ctx context.Context, username string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetUserByName, username),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetUserById(ctx context.Context, userId string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetUserById, userId),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetMe(ctx context.Context) (sreq.H, error) {
	return bot.GetUserById(ctx, bot.GetUserId())
}

func (bot *Bot) SearchUser(ctx context.Context, query string, count int) (sreq.H, error) {
	params := sreq.Params{
		"query":           query,
		"count":           count,
		"search_surface":  "user_search_page",
		"timezone_offset": bot.timeOffset,
	}
	req, _ := sreq.NewRequest(sreq.MethodGet, apiSearchUser,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}
