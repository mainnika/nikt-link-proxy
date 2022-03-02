package datasource

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/utils"
)

const domainIDdefault = "_"

const (
	redisKeyLastID   = "lastID:%s"
	redisKeyLink     = "link:%s:%s"
	redisKeyMetadata = "metadata:%s:%s"
	redisKeyReversed = "reversed:%s:%s"
)

const (
	RedisScriptIDGetOrCreate = 0
)

var redisScriptsData = map[int]string{
	RedisScriptIDGetOrCreate: `
-- KEYS[1], ARGV[1] → redisKeyReversed, linkID
-- KEYS[2], ARGV[2] → redisKeyLink, fullURL
-- KEYS[3], ARGV[3] → redisKeyMetadata, metadata

local newlyCreated = redis.call("MSETNX", KEYS[1], ARGV[1], KEYS[2], ARGV[2], KEYS[3], ARGV[3])
if newlyCreated == 1 then
  return ARGV[1]
end

return redis.call("GET", KEYS[1])
`,
}

var _ DataSource = (*RedisSource)(nil)

// RedisSource uses redis database to handle link data
type RedisSource struct {
	redis.UniversalClient

	p big.Int
	q big.Int

	scriptSHAs map[int]string
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

// Sync initializes the redis scripts
func (r *RedisSource) Sync(ctx context.Context) (err error) {

	r.scriptSHAs = map[int]string{}

	for _, scriptID := range []int{
		RedisScriptIDGetOrCreate,
	} {
		script, hasScript := redisScriptsData[scriptID]
		if !hasScript {
			panic(fmt.Errorf("no script with id %d", scriptID))
		}

		stringCmd := r.ScriptLoad(ctx, script)
		scriptSHA := ""
		scriptSHA, err = stringCmd.Result()
		if err != nil {
			return
		}

		logrus.Debugf("Redis script load success, sha:%s", scriptSHA)

		r.scriptSHAs[scriptID] = scriptSHA
	}

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

// InsertURL saves a pair short→full, in case of dup of fullURL+metadata returns the old value
func (r *RedisSource) InsertURL(ctx context.Context, shortID, fullURL string, metadata ...MetadataOpts) (insertedID string, err error) {

	metadataValues := url.Values{}
	for _, m := range metadata {
		metadataValues = m(metadataValues)
	}

	metadataEncoded := metadataValues.Encode()
	hashed, err := utils.HashStrings(fullURL, metadataEncoded)
	if err != nil {
		return
	}

	getOrCreate, hasScript := r.scriptSHAs[RedisScriptIDGetOrCreate]
	if !hasScript {
		panic(fmt.Errorf("no script sha in cache"))
	}

	cmd := r.EvalSha(ctx,
		getOrCreate,
		[]string{
			r.getKeyReversed(hashed[:]),
			r.getKeyLink(shortID),
			r.getKeyMetadata(shortID),
		},
		shortID,
		fullURL,
		metadataEncoded,
	)

	insertedID, err = cmd.Text()

	return
}

// GetFull loads a saved full URL by a link ID
func (r *RedisSource) GetFull(ctx context.Context, linkID string) (fullURL string, err error) {

	stringCmd := r.Get(ctx, r.getKeyLink(linkID))
	fullURL, err = stringCmd.Result()

	return
}

func (r *RedisSource) getKeyLastID() string {
	return fmt.Sprintf(redisKeyLastID, domainIDdefault)
}
func (r *RedisSource) getKeyLink(id string) string {
	return fmt.Sprintf(redisKeyLink, domainIDdefault, id)
}
func (r *RedisSource) getKeyReversed(hash []byte) string {
	return fmt.Sprintf(redisKeyReversed, domainIDdefault, hash)
}
func (r *RedisSource) getKeyMetadata(metadata string) string {
	return fmt.Sprintf(redisKeyMetadata, domainIDdefault, metadata)
}
