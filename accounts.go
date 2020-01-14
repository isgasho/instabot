package instabot

import (
	"context"

	"github.com/winterssy/sreq"
)

const (
	apiLogin               = apiV1 + "accounts/login/"
	apiTwoFactorLogin      = apiV1 + "accounts/two_factor_login/"
	apiLogout              = apiV1 + "accounts/logout/"
	apiSetAccountPrivate   = apiV1 + "accounts/set_private/"
	apiSetAccountPublic    = apiV1 + "accounts/set_public/"
	apiSetAccountBiography = apiV1 + "accounts/set_biography/"
	apiSetAccountGender    = apiV1 + "accounts/set_gender/"
)

func (bot *Bot) Login(ctx context.Context, interactive bool) (sreq.H, error) {
	form := sreq.Form{
		"username":            bot.Username,
		"password":            bot.Password,
		"_csrftoken":          bot.GetCSRFToken(),
		"guid":                bot.uuid,
		"phone_id":            bot.phoneId,
		"device_id":           bot.deviceId,
		"login_attempt_count": "0",
	}

	req, _ := sreq.NewRequest(sreq.MethodPost, apiLogin,
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	data, err := bot.sendRequest(req)
	if err != nil && data != nil {
		if data.GetBool("two_factor_required") {
			if interactive {
				twoFactorId := data.GetH("two_factor_info").GetString("two_factor_identifier")
				twoFactorCode := readUserInput("Enter 2FA verification code: ")
				return bot.TwoFactorLogin(ctx, twoFactorId, twoFactorCode)
			}
			return data, ErrTwoFactorRequired
		}
		if data.GetString("error_type") == "checkpoint_challenge_required" {
			return data, ErrChallengeRequired
		}
	}

	return data, err
}

func (bot *Bot) TwoFactorLogin(ctx context.Context, twoFactorId string, twoFactorCode string) (sreq.H, error) {
	form := sreq.Form{
		"two_factor_identifier": twoFactorId,
		"verification_code":     twoFactorCode,
		"username":              bot.Username,
		"password":              bot.Password,
		"device_id":             bot.deviceId,
		"ig_sig_key_version":    igSigKeyVersion,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, apiTwoFactorLogin,
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) Logout(ctx context.Context) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"guid":       bot.uuid,
		"_uuid":      bot.uuid,
		"phone_id":   bot.phoneId,
		"device_id":  bot.deviceId,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, apiLogout,
		sreq.WithForm(form),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) SetAccountPrivate(ctx context.Context, undo bool) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.uuid,
		"_uuid":      bot.uuid,
	}

	url := apiSetAccountPrivate
	if undo {
		url = apiSetAccountPublic
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, url,
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) SetAccountBiography(ctx context.Context, biography string) (sreq.H, error) {
	form := sreq.Form{
		"raw_text":   biography,
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.uuid,
		"_uuid":      bot.uuid,
		"device_id":  bot.deviceId,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, apiSetAccountBiography,
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) SetAccountGender(ctx context.Context, gender string) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.uuid,
		"_uuid":      bot.uuid,
		"device_id":  bot.deviceId,
	}

	switch gender {
	case "1", "2", "3":
		form.Set("gender", gender)
	default:
		form.Set("gender", "4")
		form.Set("custom_gender", gender)
	}

	req, _ := sreq.NewRequest(sreq.MethodPost, apiSetAccountGender,
		sreq.WithForm(form),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}
