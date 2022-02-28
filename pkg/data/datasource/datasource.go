package datasource

import (
	"bytes"
	"context"
	"encoding/base64"
)

const (
	shortIDbytesLength   = 8
	shortIDencodedLength = (shortIDbytesLength*8 + 5) / 6 // EncodedLen of shortIDbytesLength with no padding
)

// DataSource is the interface to retrieve data from the source
type DataSource interface {
	// Sync initiates data source initial state
	Sync(ctx context.Context) (err error)
	// CreateShortID returns the next unique shortID
	CreateShortID(ctx context.Context) (shortID string, err error)
	// InsertURL saves a pair short→full
	InsertURL(ctx context.Context, shortID, fullURL string, metadata ...MetadataOpts) (err error)
	// GetFull loads a saved full URL by a short ID
	GetFull(ctx context.Context, shortID string) (fullURL string, err error)
	// AddMetric modifies a metric data
	AddMetric(ctx context.Context, shortID string, metricID int, dataMod int) (err error)
	// GetMetric loads a metric data
	GetMetric(ctx context.Context, shortID string, metricID int) (data int, err error)
}

// MetadataOpts functor to add metadata to the url
type MetadataOpts func(map[string][]string) map[string][]string

func makeEncoded(shortIDbytes [shortIDbytesLength]byte) (shortIDencoded [shortIDencodedLength]byte) {
	base64.RawURLEncoding.Encode(shortIDencoded[:], shortIDbytes[:])
	return // must always fill whole buffer of shortIDencoded
}
func makeTrimmed(shortIDencoded [shortIDencodedLength]byte) (shortIDtrimmed []byte) {
	shortIDtrimmed = bytes.TrimLeftFunc(shortIDencoded[:], func(r rune) bool { return r == 'A' })
	return // returns slice of memory shortIDencoded, no alloc here
}
