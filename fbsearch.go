package instabot

import (
	"context"
	"fmt"

	"github.com/winterssy/sreq"
)

const (
	apiGetSuggestedSearches = apiV1 + "fbsearch/suggested_searches/?type=%s"
	apiGetTopSearches       = apiV1 + "fbsearch/topsearch_flat/"
	apiGetRecentSearches    = apiV1 + "fbsearch/recent_searches/"
	apiSearchPlaces         = apiV1 + "fbsearch/places/"
)

func (bot *Bot) GetSuggestedSearches(ctx context.Context, _type string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetSuggestedSearches, _type),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetTopSearches(ctx context.Context, query string, count int) (sreq.H, error) {
	params := sreq.Params{
		"query":           query,
		"count":           count,
		"search_surface":  "top_search_page",
		"timezone_offset": bot.timeOffset,
		"context":         "blended",
	}
	req, _ := sreq.NewRequest(sreq.MethodGet, apiGetTopSearches,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetRecentSearches(ctx context.Context) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, apiGetRecentSearches,
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) SearchPlaces(ctx context.Context, query string, count int) (sreq.H, error) {
	params := sreq.Params{
		"query":           query,
		"count":           count,
		"search_surface":  "places_search_page",
		"timezone_offset": bot.timeOffset,
	}
	req, _ := sreq.NewRequest(sreq.MethodGet, apiSearchPlaces,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}
