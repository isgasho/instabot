package instabot

import (
	"context"
	"fmt"

	"github.com/winterssy/sreq"
)

const (
	apiSearchTags     = apiV1 + "tags/search/"
	apiGetTagInfo     = apiV1 + "tags/%s/info/"
	apiGetTagSections = apiV1 + "tags/%s/sections/"
	apiGetTagStory    = apiV1 + "tags/%s/story/"
	apiFollowTag      = apiV1 + "tags/follow/%s/"
	apiUnfollowTag    = apiV1 + "tags/unfollow/%s/"
)

func (bot *Bot) SearchTags(ctx context.Context, query string, count int) (sreq.H, error) {
	params := sreq.Params{
		"query":           query,
		"count":           count,
		"search_surface":  "hashtag_search_page",
		"timezone_offset": bot.timeOffset,
	}
	req, _ := sreq.NewRequest(sreq.MethodGet, apiSearchTags,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetTagInfo(ctx context.Context, name string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetTagInfo, name),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

// supported_tabs: ["top","recent"]
func (bot *Bot) GetTagSections(ctx context.Context, name string, tab string, page int) (sreq.H, error) {
	form := sreq.Form{
		"tab":                tab,
		"page":               page,
		"_csrftoken":         bot.GetCSRFToken(),
		"_uuid":              bot.GetUUID(),
		"rank_token":         GenerateUUID(),
		"include_persistent": false,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(apiGetTagSections, name),
		sreq.WithForm(form),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetTagStory(ctx context.Context, name string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetTagStory, name),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) FollowTag(ctx context.Context, name string, undo bool) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.GetUUID(),
	}

	url := apiFollowTag
	if undo {
		url = apiUnfollowTag
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(url, name),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}
