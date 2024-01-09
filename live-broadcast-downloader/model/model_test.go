package model

import "testing"

func TestTask_M3U8Url(t *testing.T) {
	task := &Task{
		Prefix:  "https://example.com/video_stream/ebimayo/1234",
		SaveTo:  "/Users/Shared/Videos",
		PageUrl: "https://example.com/video/ebimayo/v/1234",
		Spec: &M3U8Spec{
			Filename: "index_1080p.m3u8",
			KeyName:  "crypt.key",
			RawQuery: "__token=qJjrt6Ky81yLuPiQxOloemZdm0k5Es+0eC02sf59Pq1r4wj0Ys1AQ723EwHSgOu7JGv/09S0JfmwCQINybR9PraX4pi6IZmcojr/YQYEiDw=gO8tI6SS0UwGoNGN",
		},
	}
	t.Log(task.M3U8Url())
}
