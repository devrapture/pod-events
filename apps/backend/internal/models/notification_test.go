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
			destination: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
			wantErr:     false,
		},
		{
			name:        "slack invalid host",
			channelType: ChannelTypeSlack,
			destination: "https://example.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
			wantErr:     true,
		},
		{
			name:        "slack non-https",
			channelType: ChannelTypeSlack,
			destination: "http://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
			wantErr:     true,
		},
		{
			name:        "slack empty path",
			channelType: ChannelTypeSlack,
			destination: "https://hooks.slack.com",
			wantErr:     true,
		},
		{
			name:        "slack services alone",
			channelType: ChannelTypeSlack,
			destination: "https://hooks.slack.com/services",
			wantErr:     true,
		},
		{
			name:        "slack partial path",
			channelType: ChannelTypeSlack,
			destination: "https://hooks.slack.com/services/T00000000/B00000000",
			wantErr:     true,
		},
		{
			name:        "slack wrong prefix",
			channelType: ChannelTypeSlack,
			destination: "https://hooks.slack.com/api/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
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
			name:        "discord missing token",
			channelType: ChannelTypeDiscord,
			destination: "https://discord.com/api/webhooks/123456789012345678",
			wantErr:     true,
		},
		{
			name:        "discord unrelated path",
			channelType: ChannelTypeDiscord,
			destination: "https://discord.com/foo",
			wantErr:     true,
		},
		{
			name:        "discord wrong prefix",
			channelType: ChannelTypeDiscord,
			destination: "https://discord.com/api/other/123456789012345678/abcdefghijklmnop",
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
