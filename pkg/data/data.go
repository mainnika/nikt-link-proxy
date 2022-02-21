package data

import (
	"context"
	"fmt"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data/datacontext"
	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/data/datasource"
)

var _ DataInterface = Data{}

// DataInterface contains function to work with data
type DataInterface interface {
}

type Data struct{}

// DataOpt is a functor to modify data opts
type DataOpt func(d Data) Data

// NewData creates a new data provider
func NewData(opts ...DataOpt) (d Data) {

	d = Data{}

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
