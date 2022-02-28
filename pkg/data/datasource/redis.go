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

const (
	redisKeyLastID      = "lastID"
	redisKeyURLMetadata = "%s:metadata"
	redisKeyURLReversed = "%s:reversed"
	redisKeyURLFull     = "%s:full"
)

const (
	RedisScriptIDGetOrCreate = 0
)

var redisScriptsData = map[int]string{
	RedisScriptIDGetOrCreate: `
local existedShort = redis.call("get", KEYS[1]) or ""
if existedShort ~= "" then
  return existedShort
end

local isSet = redis.call("setnx", KEYS[2], ARGV[2])
if isSet == 0 then
  return error("dup")
end

redis.call("set", KEYS[1], ARGV[1])
redis.call("set", KEYS[3], ARGV[3])

return ARGV[1]
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
			r.getKeyURLReversed(hashed[:]),
			r.getKeyURLFull(shortID),
			r.getKeyURLMetadata(shortID),
		},
		shortID,
		fullURL,
		metadataEncoded,
	)

	insertedID, err = cmd.Text()

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
