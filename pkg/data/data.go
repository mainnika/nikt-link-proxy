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
	MakeLink(fullURL string) (l Link, err error)
	// ResolveLink resolves a short link by ID
	ResolveLink(linkID string) (l Link, err error)
	// FireMetric increases a specific metric for a link
	FireMetric(linkID string, metric Metric, mod int) (err error)
}

type Data struct{}

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

// SourceByContext returns database session from the routing context
func (d Data) SourceByContext(c context.Context) (source datasource.DataSource) {

	source, hasSource := c.Value(datacontext.DataSourceKey).(datasource.DataSource)
	if !hasSource {
		panic(fmt.Errorf("no data source in the context"))
	}

	return
}

// MakeLink makes a short link data for a given full URL
func (d *Data) MakeLink(fullURL string) (l Link, err error) {
	//TODO implement me
	panic("implement me")
}

// ResolveLink resolves a short link by ID
func (d *Data) ResolveLink(linkID string) (l Link, err error) {
	//TODO implement me
	panic("implement me")
}

// FireMetric increases a specific metric for a link
func (d *Data) FireMetric(linkID string, metric Metric, mod int) (err error) {
	//TODO implement me
	panic("implement me")
}
