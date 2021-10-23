package platform

import (
	"github.com/jdkato/pb/internal/config"
	msdk "github.com/medium/medium-sdk-go"
)

func (c Converter) toMedium() error {
	post, err := toMediumMarkdown(c.body)
	if err != nil {
		return err
	}
	m := msdk.NewClientWithAccessToken(config.Auth.Medium)

	u, err := m.GetUser("")
	if err != nil {
		return err
	}

	_, err = m.CreatePost(msdk.CreatePostOptions{
		UserID:        u.ID,
		Title:         post.meta.Title,
		Content:       post.body,
		Tags:          post.meta.Tags,
		ContentFormat: msdk.ContentFormatHTML,
		PublishStatus: msdk.PublishStatusDraft,
	})

	return err
}
