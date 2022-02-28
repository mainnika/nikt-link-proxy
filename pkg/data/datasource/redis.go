package datasource

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/go-redis/redis/v8"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/utils"
)

const (
	redisKeyLastID      = "lastID"
	redisKeyURLMetadata = "%s:metadata"
	redisKeyURLReversed = "%s:reversed"
	redisKeyURLFull     = "%s:full"
)

var _ DataSource = (*RedisSource)(nil)

// RedisSource uses redis database to handle link data
type RedisSource struct {
	redis.UniversalClient

	p big.Int
	q big.Int
}

// RedisOpt is a functor to initialize redis source values
type RedisOpt func(*RedisSource) *RedisSource

// WithRedisPQ initializes P and Q for data id generation
func WithRedisPQ(p string, q uint64) RedisOpt {
	return func(s *RedisSource) *RedisSource {
		s.p.SetString(p, 10)
		s.q.SetUint64(q)
		return s
	}
}

// NewRedisSource creates a new redis source
func NewRedisSource(uc redis.UniversalClient, opts ...RedisOpt) (source *RedisSource) {

	source = &RedisSource{UniversalClient: uc}
	for _, f := range opts {
		source = f(source)
	}

	return source
}

// Sync does nothing in redis
func (r *RedisSource) Sync(ctx context.Context) (err error) {
	return
}

// CreateShortID returns the next unique shortID
func (r *RedisSource) CreateShortID(ctx context.Context) (shortID string, err error) {

	intCmd := r.Incr(ctx, r.getKeyLastID())

	lastID, err := intCmd.Uint64()
	if err != nil {
		return
	}

	bigLastID := &big.Int{}
	bigLastID.SetUint64(lastID)

	bigShortID := &big.Int{}
	bigShortID.Mul(bigLastID, &r.p)
	bigShortID.Mod(bigShortID, &r.q)

	var shortIDbytes [shortIDbytesLength]byte
	bigShortID.FillBytes(shortIDbytes[:])

	shortIDencoded := makeEncoded(shortIDbytes)
	shortIDtrimmed := makeTrimmed(shortIDencoded)

	shortID = string(shortIDtrimmed)

	return
}

// InsertURL saves a pair short→full
func (r *RedisSource) InsertURL(ctx context.Context, shortID, fullURL string, metadata ...MetadataOpts) (err error) {

	metadataValues := url.Values{}
	for _, m := range metadata {
		metadataValues = m(metadataValues)
	}

	metadataEncoded := metadataValues.Encode()
	hashed, err := utils.HashStrings(fullURL, metadataEncoded)
	if err != nil {
		return
	}

	p := r.Pipeline()
	_ = p.SetNX(ctx, r.getKeyURLFull(shortID), fullURL, 0)
	_ = p.SetNX(ctx, r.getKeyURLMetadata(shortID), metadataEncoded, 0)
	_ = p.Set(ctx, r.getKeyURLReversed(hashed[:]), shortID, 0)

	_, err = p.Exec(ctx)

	return
}

// GetFull loads a saved full URL by a short ID
func (r *RedisSource) GetFull(ctx context.Context, shortID string) (fullURL string, err error) {

	stringCmd := r.Get(ctx, r.getKeyURLFull(shortID))
	fullURL, err = stringCmd.Result()

	return
}

func (r *RedisSource) getKeyLastID() string {
	return redisKeyLastID
}
func (r *RedisSource) getKeyURLMetadata(id string) string {
	return fmt.Sprintf(redisKeyURLMetadata, id)
}
func (r *RedisSource) getKeyURLReversed(hash []byte) string {
	return fmt.Sprintf(redisKeyURLReversed, hash)
}
func (r *RedisSource) getKeyURLFull(id string) string {
	return fmt.Sprintf(redisKeyURLFull, id)
}
