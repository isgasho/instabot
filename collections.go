package instabot

import (
	"context"
	"fmt"

	"github.com/winterssy/sreq"
)

const (
	apiGetCollections   = apiV1 + "collections/list/"
	apiCreateCollection = apiV1 + "collections/create/"
	apiEditCollection   = apiV1 + "collections/%s/edit/"
	apiDeleteCollection = apiV1 + "collections/%s/delete/"
)

func (bot *Bot) GetCollections(ctx context.Context) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, apiGetCollections,
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) CreateCollection(ctx context.Context, name string) (sreq.H, error) {
	form := sreq.Form{
		"name":       name,
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, apiCreateCollection,
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) EditCollection(ctx context.Context, collectionId string, name string) (sreq.H, error) {
	form := sreq.Form{
		"name":       name,
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(apiEditCollection, collectionId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) DeleteCollection(ctx context.Context, collectionId string) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(apiDeleteCollection, collectionId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}
