package spotify

import "testing"

func TestSpotifyEpisodeParsedReleaseDate(t *testing.T) {
	tests := []struct {
		releaseDate string
		want        string
	}{
		{releaseDate: "2024", want: "2024-01-01"},
		{releaseDate: "2024-06", want: "2024-06-01"},
		{releaseDate: "2024-06-15", want: "2024-06-15"},
	}

	for _, tt := range tests {
		t.Run(tt.releaseDate, func(t *testing.T) {
			episode := SpotifyEpisode{ReleaseDate: tt.releaseDate}
			got, err := episode.ParsedReleaseDate()
			if err != nil {
				t.Fatalf("ParsedReleaseDate() error = %v", err)
			}
			if got.Format("2006-01-02") != tt.want {
				t.Errorf("ParsedReleaseDate() = %s, want %s", got.Format("2006-01-02"), tt.want)
			}
		})
	}
}

func TestSpotifyEpisodeParsedReleaseDateRejectsInvalidDate(t *testing.T) {
	episode := SpotifyEpisode{ReleaseDate: "not-a-date"}
	if _, err := episode.ParsedReleaseDate(); err == nil {
		t.Fatal("ParsedReleaseDate() error = nil, want an error")
	}
}
