package datasource

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var _ DataSource = (*RedisSource)(nil)

// RedisSource uses redis database to handle link data
type RedisSource struct {
	redis.UniversalClient
}

// NewRedisSource creates a new redis source
func NewRedisSource(uc redis.UniversalClient) *RedisSource {
	return &RedisSource{UniversalClient: uc}
}

// Sync initiates redis source initial state
func (r *RedisSource) Sync(ctx context.Context) (err error) {
	//TODO implement me
	panic("implement me")
}

// CreateShortID returns the next unique shortID
func (r *RedisSource) CreateShortID(ctx context.Context) (shortID string, err error) {
	//TODO implement me
	panic("implement me")
}

// InsertURL saves a pair short→full
func (r *RedisSource) InsertURL(ctx context.Context, shortID, fullURL string, metadata ...MetadataOpts) (err error) {
	//TODO implement me
	panic("implement me")
}

// GetFull loads a saved full URL by a short ID
func (r *RedisSource) GetFull(ctx context.Context, shortID string) (fullURL string, err error) {
	//TODO implement me
	panic("implement me")
}

// AddMetric modifies a metric data
func (r *RedisSource) AddMetric(ctx context.Context, shortID string, metricID int, dataMod int) (err error) {
	//TODO implement me
	panic("implement me")
}

// GetMetric loads a metric data
func (r *RedisSource) GetMetric(ctx context.Context, shortID string, metricID int) (data int, err error) {
	//TODO implement me
	panic("implement me")
}
