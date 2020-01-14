package instabot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/winterssy/sreq"
)

const (
	apiUploadPhoto           = "https://i.instagram.com/rupload_igphoto/%s"
	apiConfigureMedia        = apiV1 + "media/configure/"
	apiConfigureSidecarMedia = apiV1 + "media/configure_sidecar/"
	apiLikeMedia             = apiV1 + "media/%s/like/"
	apiUnlikeMedia           = apiV1 + "media/%s/unlike/"
	apiSaveMedia             = apiV1 + "media/%s/save/"
	apiUnsaveMedia           = apiV1 + "media/%s/unsave/"
	apiSendMediaComment      = apiV1 + "media/%s/comment/"
	apiLikeMediaComment      = apiV1 + "media/%s/comment_like/"
	apiUnlikeMediaComment    = apiV1 + "media/%s/comment_unlike/"
	apiGetMediaComments      = apiV1 + "media/%s/comments/"
	apiDeleteMediaComments   = apiV1 + "media/%s/comment/bulk_delete/"
	apiDisableMediaComments  = apiV1 + "media/%s/disable_comments/"
	apiEnableMediaComments   = apiV1 + "media/%s/enable_comments/"
	apiEditMedia             = apiV1 + "media/%s/edit_media/"
	apiArchiveMedia          = apiV1 + "media/%s/only_me/"
	apiUndoArchiveMedia      = apiV1 + "media/%s/undo_only_me/"
	apiDeleteMedia           = apiV1 + "media/%s/delete/"

	imageHorizontalLimit = 2.0
	imageVerticalLimit   = 0.8
	imageMaxWidth        = 1080
	imageMaxHeight       = imageMaxWidth
)

func adjustImage(img image.Image) (image.Image, error) {
	rectangle := img.Bounds()
	width, height := rectangle.Dx(), rectangle.Dy()
	ratio := float64(width) / float64(height)

	if width > height {
		if ratio > imageHorizontalLimit {
			width = width - int(math.Ceil(float64(width)-float64(height)*2))
			img = imaging.CropCenter(img, width, height)
		}
		if width > 1080 {
			height = int(math.Ceil(imageMaxWidth * float64(height) / float64(width)))
			img = imaging.Resize(img, imageMaxWidth, height, imaging.Lanczos)
		}
	} else if width < height {
		if ratio < imageVerticalLimit {
			height = height - int(math.Ceil(float64(height)-float64(width)*1.25))
			img = imaging.CropCenter(img, width, height)
		}
		if height > imageMaxHeight {
			width = int(math.Ceil(imageMaxHeight * float64(width) / float64(height)))
			img = imaging.Resize(img, width, imageMaxHeight, imaging.Lanczos)
		}
	} else {
		if width > imageMaxWidth {
			img = imaging.Resize(img, imageMaxWidth, imageMaxHeight, imaging.Lanczos)
		}
	}

	return img, nil
}

func (bot *Bot) UploadPhoto(ctx context.Context, r io.Reader, isSidecar bool) (sreq.H, error) {
	img, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	img, err = adjustImage(img)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = imaging.Encode(&buf, img, imaging.JPEG)
	if err != nil {
		return nil, err
	}

	ts := time.Now().UnixNano() / 1e6
	rUploadParams := sreq.Values{
		"upload_id":         ts,
		"media_type":        "1",
		"retry_context":     "{\"num_reupload\":0,\"num_step_auto_retry\":0,\"num_step_manual_retry\":0}",
		"image_compression": "{\"lib_name\":\"moz\",\"lib_version\":\"3.1.m\",\"quality\":\"95\"}",
	}
	if isSidecar {
		rUploadParams.Set("is_sidecar", "1")
	}

	uploadName := fmt.Sprintf("%d_0_%d", ts, 1e9+rand.Intn(1e10-1e9))
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(apiUploadPhoto, uploadName),
		sreq.WithBody(&buf),
		sreq.WithHeaders(sreq.Headers{
			"X_FB_PHOTO_WATERFALL_ID":    GenerateUUID(),
			"X-Entity-Length":            buf.Len(),
			"X-Entity-Name":              uploadName,
			"X-Instagram-Rupload-Params": rUploadParams.Marshal(),
			"X-Entity-Type":              "image/jpeg",
			"Offset":                     0,
			"Content-Type":               "application/octet-stream",
		}),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) uploadPhoto(ctx context.Context, imageFile string, isSidecar bool) (sreq.H, error) {
	file, err := os.Open(imageFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return bot.UploadPhoto(ctx, file, isSidecar)
}

func (bot *Bot) ConfigureMedia(ctx context.Context, uploadId string,
	caption string, disableComments bool) (sreq.H, error) {
	form := sreq.Form{
		"timezone_offset": bot.timeOffset,
		"_csrftoken":      bot.GetCSRFToken(),
		"source_type":     "4",
		"_uid":            bot.GetUserId(),
		"_uuid":           bot.uuid,
		"device_id":       bot.deviceId,
		"caption":         caption,
		"upload_id":       uploadId,
		"device":          deviceSettings,
	}
	if disableComments {
		form.Set("disable_comments", "1")
	}

	req, _ := sreq.NewRequest(sreq.MethodPost, apiConfigureMedia,
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) PostPhoto(ctx context.Context, imageFile string, caption string,
	disableComments bool) (sreq.H, error) {
	data, err := bot.uploadPhoto(ctx, imageFile, false)
	if err != nil {
		return data, err
	}

	return bot.ConfigureMedia(ctx, data.GetString("upload_id"), caption, disableComments)
}

func (bot *Bot) ConfigureSidecarMedia(ctx context.Context, uploadIds []string, caption string,
	disableComments bool) (sreq.H, error) {
	ts := time.Now().UnixNano() / 1e6
	form := sreq.Form{
		"timezone_offset":   bot.timeOffset,
		"_csrftoken":        bot.GetCSRFToken(),
		"source_type":       "4",
		"_uid":              bot.GetUserId(),
		"device_id":         bot.deviceId,
		"_uuid":             bot.uuid,
		"caption":           caption,
		"upload_id":         ts,
		"client_sidecar_id": ts,
		"device":            deviceSettings,
	}

	childrenMetadata := make([]map[string]interface{}, len(uploadIds))
	for i, uploadId := range uploadIds {
		childrenMetadata[i] = map[string]interface{}{
			"upload_id":       uploadId,
			"timezone_offset": bot.timeOffset,
			"source_type":     "4",
			"device":          deviceSettings,
		}
	}
	form.Set("children_metadata", childrenMetadata)
	if disableComments {
		form.Set("disable_comments", "1")
	}

	req, _ := sreq.NewRequest(sreq.MethodPost, apiConfigureSidecarMedia,
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) PostAlbum(ctx context.Context, imageFiles []string, caption string,
	disableComments bool) (sreq.H, error) {
	n := len(imageFiles)
	if n < 2 || n > 10 {
		return nil, errors.New("album requires 2 to 10 photos")
	}

	uploadIds := make([]string, 0, n)
	for _, imageFile := range imageFiles {
		data, err := bot.uploadPhoto(ctx, imageFile, true)
		if err != nil {
			return data, err
		}
		uploadIds = append(uploadIds, data.GetString("upload_id"))
	}
	return bot.ConfigureSidecarMedia(ctx, uploadIds, caption, disableComments)
}

func (bot *Bot) SaveMediaToCollections(ctx context.Context, mediaId string, undo bool,
	collectionIds ...string) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}

	v := "[" + strings.Join(collectionIds, ",") + "]"
	if undo {
		form.Set("removed_collection_ids", v)
	} else {
		form.Set("added_collection_ids", v)
	}

	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(apiSaveMedia, mediaId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) GetMediaComments(ctx context.Context, mediaId string) (sreq.H, error) {
	req, _ := sreq.NewRequest(sreq.MethodGet, fmt.Sprintf(apiGetMediaComments, mediaId),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) SendMediaComment(ctx context.Context, mediaId string, commentText string) (sreq.H, error) {
	form := sreq.Form{
		"comment_text": commentText,
		"_csrftoken":   bot.GetCSRFToken(),
		"_uid":         bot.GetUserId(),
		"_uuid":        bot.uuid,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(apiSendMediaComment, mediaId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) ReplyMediaComment(ctx context.Context, mediaId string, commentId string,
	commentText string) (sreq.H, error) {
	form := sreq.Form{
		"replied_to_comment_id": commentId,
		"comment_text":          commentText,
		"_csrftoken":            bot.GetCSRFToken(),
		"_uid":                  bot.GetUserId(),
		"_uuid":                 bot.uuid,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(apiSendMediaComment, mediaId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) LikeMediaComment(ctx context.Context, commentId string, undo bool) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}

	url := apiLikeMediaComment
	if undo {
		url = apiUnlikeMediaComment
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(url, commentId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) DeleteMediaComments(ctx context.Context, mediaId string, commentIds ...string) (sreq.H, error) {
	form := sreq.Form{
		"comment_ids_to_delete": strings.Join(commentIds, ","),
		"_csrftoken":            bot.GetCSRFToken(),
		"_uid":                  bot.GetUserId(),
		"_uuid":                 bot.uuid,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(apiDeleteMediaComments, mediaId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) DisableMediaComments(ctx context.Context, mediaId string, undo bool) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}

	url := apiDisableMediaComments
	if undo {
		url = apiEnableMediaComments
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(url, mediaId),
		sreq.WithForm(form),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) LikeMedia(ctx context.Context, mediaId string, undo bool) (sreq.H, error) {
	form := sreq.Form{
		"media_id":   mediaId,
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}
	signedForm := GenerateSignedForm(form)
	doubleTap := rand.Intn(1)
	signedForm.Set("d", doubleTap)

	url := apiLikeMedia
	if undo {
		url = apiUnlikeMedia
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(url, mediaId),
		sreq.WithForm(signedForm),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) SaveMedia(ctx context.Context, mediaId string, undo bool) (sreq.H, error) {
	form := sreq.Form{
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}

	url := apiSaveMedia
	if undo {
		url = apiUnsaveMedia
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(url, mediaId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) EditMedia(ctx context.Context, mediaId string, captionText string) (sreq.H, error) {
	form := sreq.Form{
		"caption_text": captionText,
		"_csrftoken":   bot.GetCSRFToken(),
		"_uid":         bot.GetUserId(),
		"_uuid":        bot.uuid,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(apiEditMedia, mediaId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) ArchiveMedia(ctx context.Context, mediaId string, undo bool) (sreq.H, error) {
	form := sreq.Form{
		"media_id":   mediaId,
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}

	url := apiArchiveMedia
	if undo {
		url = apiUndoArchiveMedia
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(url, mediaId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}

func (bot *Bot) DeleteMedia(ctx context.Context, mediaId string) (sreq.H, error) {
	form := sreq.Form{
		"media_id":   mediaId,
		"_csrftoken": bot.GetCSRFToken(),
		"_uid":       bot.GetUserId(),
		"_uuid":      bot.uuid,
	}
	req, _ := sreq.NewRequest(sreq.MethodPost, fmt.Sprintf(apiDeleteMedia, mediaId),
		sreq.WithForm(GenerateSignedForm(form)),
		sreq.WithContext(ctx),
	)
	return bot.sendRequest(req)
}
