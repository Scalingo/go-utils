package influx

import (
	"context"
	"net/url"

	influx "github.com/influxdata/influxdb/client/v2"

	"github.com/Scalingo/go-utils/errors/v3"
)

// Do actually executes the query to the specified InfluxDB instance hosted at the url argument.
func (q Query) Do(ctx context.Context, url string) (*influx.Response, error) {
	query := q.Build()
	response, err := executeQuery(ctx, url, query)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "executing query "+query)
	}
	return response, nil
}

func executeQuery(ctx context.Context, url, queryString string) (*influx.Response, error) {
	client, dbName, err := newClient(ctx, url)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "creating InfluxDB client")
	}
	defer client.Close()

	response, err := client.Query(influx.Query{
		Command:  queryString,
		Database: dbName,
	})
	if err != nil {
		return nil, errors.Wrap(ctx, err, "querying InfluxDB")
	}
	if response.Error() != nil {
		return nil, errors.Wrap(ctx, response.Error(), "fetching data from InfluxDB")
	}

	return response, nil
}

type influxInfo struct {
	host             string
	user             string
	password         string
	database         string
	connectionString string
}

func newClient(ctx context.Context, url string) (influx.Client, string, error) {
	infos, err := parseConnectionString(ctx, url)
	if err != nil {
		return nil, "", errors.Wrap(ctx, err, "parsing connection string")
	}
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:      infos.host,
		Username:  infos.user,
		Password:  infos.password,
		UserAgent: "Scalingo Utils",
	})

	if err != nil {
		return nil, "", errors.Wrap(ctx, err, "creating HTTP InfluxDB client")
	}

	return client, infos.database, nil
}

func parseConnectionString(ctx context.Context, con string) (*influxInfo, error) {
	url, err := url.Parse(con)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "parsing connection string")
	}

	var user, password string
	if url.User != nil {
		password, _ = url.User.Password()
		user = url.User.Username()
	}

	return &influxInfo{
		host:             url.Scheme + "://" + url.Host,
		user:             user,
		password:         password,
		database:         url.Path[1:],
		connectionString: con,
	}, nil
}
