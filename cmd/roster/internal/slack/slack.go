package slack

import (
	"log/slog"

	goslack "github.com/slack-go/slack"
)

type SlackClient struct {
        SlackClient *goslack.Client
        ChannelID string
}

var Client *SlackClient

func Init(token string, channel_id string) error {
	Client = &SlackClient {
		SlackClient: goslack.New(token),
		ChannelID: channel_id,
	}

	return nil
}

func (s *SlackClient) SendMessage(message string) error {
        ChannelID, timestamp, err := s.SlackClient.PostMessage(s.ChannelID, goslack.MsgOptionText(message, false))

        if err != nil {
                return err
        }

        slog.Info("sendSlackMessage", "channel", ChannelID, "message", message, "timestamp",timestamp)
        return nil
}
