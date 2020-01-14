package instabot

import (
	"context"
	"fmt"

	"github.com/winterssy/sreq"
)

const (
	apiGetUserFriendship = apiV1 + "friendships/show/%s/"
	apiGetUserFollowings = apiV1 + "friendships/%s/following/?max_id=%s"
	apiGetUserFollowers  = apiV1 + "friendships/%s/followers/?max_id=%s"
	apiFollowUser        = apiV1 + "friendships/create/%s/"
	apiUnfollowUser      = apiV1 + "friendships/destroy/%s/"
	apiBlockUser         = apiV1 + "friendships/block/%s/"
	apiUnblockUser       = apiV1 + "friendships/unblock/%s/"
)

func (bot *Bot) GetUserFriendship(ctx context.Context, userId string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetUserFriendship, userId),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetUserFollowings(ctx context.Context, userId string, maxId string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetUserFollowings, userId, maxId),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetUserFollowers(ctx context.Context, userId string, maxId string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetUserFollowers, userId, maxId),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetSelfFollowings(ctx context.Context, maxId string) (sreq.H, error) {
	return bot.GetUserFollowings(ctx, bot.GetUserId(), maxId)
}

func (bot *Bot) GetSelfFollowers(ctx context.Context, maxId string) (sreq.H, error) {
	return bot.GetUserFollowers(ctx, bot.GetUserId(), maxId)
}

func (bot *Bot) FollowUser(ctx context.Context, userId string, undo bool) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"user_id":    userId,
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}

	url := apiFollowUser
	if undo {
		url = apiUnfollowUser
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(url, userId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) BlockUser(ctx context.Context, userId string, undo bool) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"user_id":    userId,
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}

	url := apiBlockUser
	if undo {
		url = apiUnblockUser
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(url, userId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}
