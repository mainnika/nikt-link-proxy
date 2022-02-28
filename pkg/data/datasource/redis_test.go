package datasource_test

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data/datasource"
)

func TestRedisSource_CreateShortID(t *testing.T) {
	t.Run("must generate uniq", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		uc := NewMockUniversalClient(ctrl)
		rs := datasource.NewRedisSource(uc, datasource.WithRedisPQ("17248111889498283943", 1024))
		uniqtable := map[string]bool{}

		t.Run("first Q must be uniq", func(t *testing.T) {
			for i := 0; i < 1024; i++ {
				incrResult := &redis.IntCmd{}
				incrResult.SetVal(int64(i))
				uc.EXPECT().Incr(gomock.Any(), gomock.Any()).Return(incrResult)

				sid, err := rs.CreateShortID(context.Background())
				hasValue, _ := uniqtable[sid]
				uniqtable[sid] = true

				assert.False(t, hasValue)
				assert.NoError(t, err)
			}
		})
		t.Run("Q+1 must collide", func(t *testing.T) {
			incrResult := &redis.IntCmd{}
			incrResult.SetVal(int64(1024))
			uc.EXPECT().Incr(gomock.Any(), gomock.Any()).Return(incrResult)

			sid, err := rs.CreateShortID(context.Background())
			hasValue, _ := uniqtable[sid]

			assert.True(t, hasValue)
			assert.NoError(t, err)
		})
	})
}
