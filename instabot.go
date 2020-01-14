package instabot

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/winterssy/sreq"
)

const (
	apiHost          = "https://i.instagram.com/"
	apiV1            = apiHost + "api/v1/"
	defaultUserAgent = "Instagram 123.0.0.21.114 Android (28/9; 480dpi; 1080x2232; meizu; 16s; 16s; qcom; zh_CN; 188791674)"
	igCapabilities   = "3brTvwM="
	igAppId          = "567067343352427"
	bloksVersionId   = "7ab39aa203b17c94cc6787d6cd9052d221683361875eee1e1bfe30b8e9debd74"
)

var (
	deviceSettings = map[string]interface{}{
		"manufacturer":    "meizu",
		"model":           "16s",
		"android_version": 28,
		"android_release": "9",
	}

	defaultHeaders = sreq.Headers{
		"X-IG-Connection-Type": "WIFI",
		"X-IG-Capabilities":    igCapabilities,
		"X-IG-App-ID":          igAppId,
		"User-Agent":           defaultUserAgent,
		"Accept-Language":      "zh-CN, en-US",
		"X-FB-HTTP-Engine":     "Liger",
		"Connection":           "keep-alive",
	}
)

var (
	ErrTwoFactorRequired = errors.New("two-factor authentication required")
	ErrChallengeRequired = errors.New("challenge required")
	ErrFeedbackRequired  = errors.New("feedback required")
)

type (
	Bot struct {
		Client   *sreq.Client `json:"-"`
		Username string       `json:"username"`
		Password string       `json:"password"`

		phoneId    string
		uuid       string
		igDeviceId string
		deviceId   string
		timeOffset int
	}
)

func New(username string, password string) *Bot {
	client := sreq.New()
	client.OnBeforeRequest(sreq.SetDefaultHeaders(defaultHeaders))
	client.Head(apiHost)
	return &Bot{
		Client:     client,
		Username:   username,
		Password:   password,
		phoneId:    GenerateUUID(),
		uuid:       GenerateUUID(),
		igDeviceId: GenerateUUID(),
		deviceId:   generateDeviceId(generateMD5Hash(username + password)),
		timeOffset: timeOffset(),
	}
}

func (bot *Bot) GetCSRFToken() string {
	cookie, err := bot.Client.FilterCookie(apiV1, "csrftoken")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func (bot *Bot) GetUserId() string {
	cookie, err := bot.Client.FilterCookie(apiV1, "ds_user_id")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func (bot *Bot) GetUUID() string {
	return bot.uuid
}

func (bot *Bot) GetDeviceId() string {
	return bot.deviceId
}

func (bot *Bot) SendRequest(req *sreq.Request) *sreq.Response {
	ts := float64(time.Now().UnixNano()) / 1e9
	req.SetHeaders(sreq.Headers{
		"X-IG-App-Locale":             "zh-CN",
		"X-IG-Device-Locale":          "zh-CN",
		"X-Pigeon-Session-Id":         GenerateUUID(),
		"X-Pigeon-Rawclienttime":      fmt.Sprintf("%.3f", ts),
		"X-IG-Connection-Speed":       "-1kbps",
		"X-IG-Bandwidth-Speed-KBPS":   "-1.000",
		"X-IG-Bandwidth-TotalBytes-B": "0",
		"X-IG-Bandwidth-TotalTime-MS": "0",
		"X-Bloks-Version-Id":          bloksVersionId,
		"X-Bloks-Is-Layout-RTL":       "false",
		"X-Bloks-Enable-RenderCore":   "false",
		"X-IG-Device-ID":              bot.igDeviceId,
		"X-IG-Android-ID":             bot.deviceId,
	})
	return bot.Client.Do(req)
}

func (bot *Bot) sendRequest(req *sreq.Request) (sreq.H, error) {
	data, err := bot.SendRequest(req).H()
	if err != nil {
		return nil, err
	}

	if data.GetStringDefault("status", "fail") != "ok" {
		message := data.GetStringDefault("message", "unknown exception")
		if strings.Contains(message, "feedback_required") {
			return data, ErrFeedbackRequired
		}
		return data, errors.New(message)
	}

	return data, nil
}
