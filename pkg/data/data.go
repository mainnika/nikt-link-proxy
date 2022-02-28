package data

import (
	"context"
	"fmt"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data/datacontext"
	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data/datasource"
)

var _ DataInterface = (*Data)(nil)

// DataInterface is a data handler that manages links
type DataInterface interface {
	// MakeLink makes a short link data for a given full URL
	MakeLink(ctx context.Context, fullURL string) (l Link, err error)
	// ResolveLink resolves a short link by ID
	ResolveLink(ctx context.Context, linkID string) (l Link, err error)
}

// Data implements data handler that works with the source to manage links
type Data struct {
	dsource datasource.DataSource
}

// DataOpt is a functor to modify data opts
type DataOpt func(d *Data) *Data

// NewData creates a new data provider
func NewData(opts ...DataOpt) (d *Data) {

	d = &Data{}

	for _, opt := range opts {
		d = opt(d)
	}

	return
}

// WithDataSource configures a data source
func WithDataSource(dsource datasource.DataSource) DataOpt {
	return func(d *Data) *Data {
		d.dsource = dsource
		return d
	}
}

// SourceByContext returns database session from the routing context
func (d Data) SourceByContext(c context.Context) (source datasource.DataSource) {

	source, hasSource := c.Value(datacontext.DataSourceKey).(datasource.DataSource)
	if !hasSource {
		panic(fmt.Errorf("no data source in the context"))
	}

	return
}

// MakeLink makes a short link data for a given full URL
func (d *Data) MakeLink(ctx context.Context, fullURL string) (l Link, err error) {

	shortID, err := d.dsource.CreateShortID(ctx)
	if err != nil {
		return
	}

	err = d.dsource.InsertURL(ctx, shortID, fullURL)
	if err != nil {
		return
	}

	l = Link{ID: shortID, FullURL: fullURL}

	return
}

// ResolveLink resolves a short link by ID
func (d *Data) ResolveLink(ctx context.Context, linkID string) (l Link, err error) {

	fullURL, err := d.dsource.GetFull(ctx, linkID)
	if err != nil {
		return
	}

	l = Link{ID: linkID, FullURL: fullURL}

	return
}
