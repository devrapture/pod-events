package models

import "testing"

func TestValidateDestination(t *testing.T) {
	tests := []struct {
		name        string
		channelType ChannelType
		destination string
		wantErr     bool
	}{
		{
			name:        "slack valid",
			channelType: ChannelTypeSlack,
			destination: "https://hooks.slack.com/services/test/webhook",
			wantErr:     false,
		},
		{
			name:        "slack invalid host",
			channelType: ChannelTypeSlack,
			destination: "https://example.com/services/test/webhook",
			wantErr:     true,
		},
		{
			name:        "slack non-https",
			channelType: ChannelTypeSlack,
			destination: "http://hooks.slack.com/services/test/webhook",
			wantErr:     true,
		},
		{
			name:        "slack empty path",
			channelType: ChannelTypeSlack,
			destination: "https://hooks.slack.com",
			wantErr:     true,
		},
		{
			name:        "discord valid",
			channelType: ChannelTypeDiscord,
			destination: "https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnop",
			wantErr:     false,
		},
		{
			name:        "discord invalid host",
			channelType: ChannelTypeDiscord,
			destination: "https://example.com/api/webhooks/123456789012345678/abcdefghijklmnop",
			wantErr:     true,
		},
		{
			name:        "telegram not validated",
			channelType: ChannelTypeTelegram,
			destination: "123456789",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.channelType.ValidateDestination(tt.destination)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateDestination() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
