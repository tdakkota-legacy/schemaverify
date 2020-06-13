package parseutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pascalToSnake(t *testing.T) {
	type pair struct {
		Pascal string
		Snake  string
	}

	testTable := []pair{
		{
			Pascal: "VideoVideoFull",
			Snake:  "video_video_full",
		},
		{
			Pascal: "WallWallpostToID",
			Snake:  "wall_wallpost_to_id",
		},
		{
			Pascal: "IDUser",
			Snake:  "id_user",
		},
	}

	for _, pair := range testTable {
		assert.Equal(t, pair.Snake, PascalToSnake(pair.Pascal))
	}
}
