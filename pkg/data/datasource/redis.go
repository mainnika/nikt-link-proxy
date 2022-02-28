package datasource

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/url"

	"github.com/go-redis/redis/v8"
)

const (
	redisKeyLastID      = "lastID"
	redisKeyURLMetadata = "%s:metadata"
	redisKeyURLFull     = "%s:full"
	redisKeyURLMetric   = "%s:metrica:%d"
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

// Sync initiates redis source initial state
func (r *RedisSource) Sync(ctx context.Context) (err error) {
	boolCmd := r.SetNX(ctx, redisKeyLastID, startID, 0)
	return boolCmd.Err()
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

	var shortIDbytes [8]byte
	bigShortID.FillBytes(shortIDbytes[:])

	shortID = base64.RawURLEncoding.EncodeToString(shortIDbytes[:])

	return
}

// InsertURL saves a pair short→full
func (r *RedisSource) InsertURL(ctx context.Context, shortID, fullURL string, metadata ...MetadataOpts) (err error) {

	metadataValues := url.Values{}
	for _, m := range metadata {
		metadataValues = m(metadataValues)
	}

	metadataEncoded := metadataValues.Encode()

	p := r.Pipeline()
	_ = p.SetNX(ctx, r.getKeyURLFull(shortID), fullURL, 0)
	_ = p.SetNX(ctx, r.getKeyURLMetadata(shortID), metadataEncoded, 0)

	_, err = p.Exec(ctx)

	return
}

// GetFull loads a saved full URL by a short ID
func (r *RedisSource) GetFull(ctx context.Context, shortID string) (fullURL string, err error) {

	stringCmd := r.Get(ctx, r.getKeyURLFull(shortID))
	fullURL, err = stringCmd.Result()

	return
}

// AddMetric modifies a metric data
func (r *RedisSource) AddMetric(ctx context.Context, shortID string, metricID int, dataMod int) (err error) {

	intCmd := r.IncrBy(ctx, r.getKeyURLMetrica(shortID, metricID), int64(dataMod))
	err = intCmd.Err()

	return
}

// GetMetric loads a metric data
func (r *RedisSource) GetMetric(ctx context.Context, shortID string, metricID int) (data int, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisSource) getKeyLastID() string {
	return redisKeyLastID
}
func (r *RedisSource) getKeyURLMetadata(id string) string {
	return fmt.Sprintf(redisKeyURLMetadata, id)
}
func (r *RedisSource) getKeyURLFull(id string) string {
	return fmt.Sprintf(redisKeyURLFull, id)
}
func (r *RedisSource) getKeyURLMetrica(id string, metricaID int) string {
	return fmt.Sprintf(redisKeyURLMetric, id, metricaID)
}
