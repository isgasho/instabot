package instabot

import (
	"context"

	"github.com/winterssy/sreq"
)

const (
	apiExplore = apiV1 + "discover/topical_explore/"
)

func (bot *Bot) Explore(ctx context.Context) (sreq.H, error) {
	params := sreq.Params{
		"is_prefetch":                true,
		"omit_cover_media":           true,
		"use_sectional_payload":      true,
		"timezone_offset":            bot.timeOffset,
		"include_fixed_destinations": true,
		"session_id":                 GenerateUUID(),
	}
	req, _ := sreq.NewRequest(sreq.MethodGet, apiExplore,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) ExploreCluster(ctx context.Context, clusterId string, maxId string) (sreq.H, error) {
	params := sreq.Params{
		"cluster_id":                 clusterId,
		"max_id":                     maxId,
		"module":                     "explore_popular",
		"is_prefetch":                false,
		"omit_cover_media":           true,
		"use_sectional_payload":      true,
		"timezone_offset":            bot.timeOffset,
		"include_fixed_destinations": true,
		"session_id":                 GenerateUUID(),
	}
	req, _ := sreq.NewRequest(sreq.MethodGet, apiExplore,
		sreq.WithQuery(params),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}
