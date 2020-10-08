# Changelog

## To Be Released

## v7.1.0 (August 13 2020)

* [difflib] Add a way to activate/desactivate the colors formatting
* [difflib] Add colors to have a better display in the shell (only for unified diff)
* [difflib] New package to print diff between two content in the shell (git diff style)
* [document] Do not add sensitive information in the log

## v7.0.1 (June 12 2020)

* [storage] Fix Swift authentication

## v7.0.0 (June 09 2020)

* [storage] Get swift configuration from the environment

## v6.7.1 (May 19 2020)

* [tarball] Fix unit: use a buffer of 512kb not MB...

## v6.7.0 (May 12 2020)

* [tarball] Use `io` package and disable disk cache when creating/extracting tarball archives

## v6.6.0 (May 4 2020)

* [io] New package with a configurable io.Copier

    ```
    io.NewCopier(...opts io.CopierOpt)
    io.WithBufferSize(int64)
    io.WithNoDiskCacheRead
    io.WithNoDiskCacheWrite
    io.WithNoDiskCache
    ```

## v6.5.0 (Apr 28 2020)

* [tarball] New package to manipulate tar and gzipped tarballs

    ```
    tarball.Create(context.Context, file string, in io.Reader, opts tarball.CreateOpts) error
    tarball.Extract(context.Context, file string, out io.Writer, opts tarball.ExtractOpts) error
    tarball.Tar(io.Writer, map[string]tarball.TarFileHeader) error
    ```

## v6.4.0 (Apr 27 2020)

* [fs] New package to manipulate the filesystem through a mockable interface

    ```
    // Mockable interface
    fs.Fs

    // Real OS filesystem manipulation
    fs.NewOsFs()
    // In-memory empty filesystem
    fs.NewMemFs()
    ```

## v6.3.0 (Apr 21 2020)

* [retry] Add WithErrorCallback option to add callbacks when a try fails

## v6.2.1 (Mar 27 2020)

* [retry] Add a constructor for RetryCancelError

## v6.2.0 (Mar 27 2020)

* [retry] Add ability to get the last error (LastErr) in case of timeout
* [retry] Add a custom error type to cancel a Retry `RetryCancelError`

## v6.1.0 (Mar 27 2020)

* [retry] Add WithMaxDuration option to set a timeout

## v6.0.0 (Mar 13 2020)

* [logger] Add `region_name` field by default if `REGION_NAME` environment variable

## v5.7.1 (Mar 19 2020)

* Fix `ValidationErrors` `Error` output [#109](https://github.com/Scalingo/go-utils/pull/109)

## v5.7.0 (Jan 29 2020)

* Add Retry interface

## v5.6.2 (Nov 14 2019)

* Update errgo, pkg/errors, logrus, and unfreeze their version

## v5.6.1 (Nov 12 2019)

* Update go-etcd-cron dependency with new etcd deps

## v5.6.0 (Nov 12 2019)

* Update deps Related to etcd (new location github.com/coreos/etcd -> go.etcd.io)

## v5.5.15 (Oct 22 2019)

* Update deps related to Rollbar/logging/errors

## v5.5.14 (August 30 2019)

* Add EnsureParanoidIndices method [#49](https://github.com/Scalingo/go-utils/pull/49)

## v5.5.13 (August 27 2019)

* [S3] Trim left slashes [#92](https://github.com/Scalingo/go-utils/pull/92)

## v5.5.11 (August 26 2019)

* Add simple object storage package for S3 and swift
  [#89](https://github.com/Scalingo/go-utils/pull/89)

## v5.5.10 (May 14 2019)

* Migration from `github.com/satori/go.uuid` to its fork `github.com/gofrs/uuid`, Reason:

> This project was originally forked from the github.com/satori/go.uuid
> repository after it appeared to be no longer maintained, while exhibiting
> critical flaws. We have decided to take over this project to ensure it receives
> regular maintenance for the benefit of the larger Go community.

## v5.5.9 (Apr 18 2019)

* [errors] Ability to merge one ValidationErrors in another

```
b := NewValidationErrorsBuilder()

if doc.Field == "" {
  b.Set("field", "is blank")
}

// { errors: { field: ["is blank"] } }
b.Merge(doc.SubDocument.Validate())

// or

// { errors: { subdoc_field: ["is blank"] } }
b.MergeWithPrefix("subdoc", doc.SubDocument.Validate())


return b.Build()
```

## v5.5.8 (Apr 16 2019)

* [paranoid] IsDeleted returns true if DeletedAt is not zero

## v5.5.7 (Apr 12 2019)

* [errors] Better handling of ValidationErrors

## v5.5.6 (Apr 11 2019)

* [mongo] Remove trailing comma in validation errors

## v5.5.5 (Apr 3 2019)

* [influx] Add support for subqueries

## v5.5.4 (Mar 28 2019)

* [influx] Always write measurement between double quotes for InfluxDB query builder

## v5.5.3 (Mar 11 2019)

* [logger] Rollbar plugin is reading ROLLBAR_ENV first, then GO_ENV

## v5.5.2 (Mar 05 2019)

* [mongo] Handle `timeout` parameter in the URL, default at 10s
* [mongo] Do not destroy URL when `ssl` or `timeout` parameter is first in the list

## v5.5.1 (Nov 26 2018)

* [nsqproducer] Fix Ping() method recursivity issue

## v5.5.0 (Nov 25 2018)

* [nsqproducer] Producers have a Ping() method for monitoring

## v5.4.1 (Oct 11 2018)

* [mongo] Make BuildSession public and usable with custom URLs

## v5.4.0 (Oct 10 2018)

* [httpclient] New library to forward request_id + a few helpers (TLS and authentication)

## v5.3.2 (Oct 3 2018)

* [mongo] Fix Iter error, return iterator error no callback one, #56
* [logger-rollbar] Switch rollbar library to use official lib

## v5.3.1 (Sep 27 2018)

* [graceful] Fix ListenAndServeTLS to correctly use a tls.Listener

## v5.3.0 (Sep 26 2018)

* Add `graceful` package for graceful restart/shutdown of HTTP servers

## v5.2.1 (Sep 25 2018)

* Add `SkipLogSet` options to NSQ LB producer env init

## v5.2.0 (Jul 20 2018)

* Add `WhereQuery` in `mongo/document`

## v4.2.3 (Jun 14 2018)

* Add ability to choose nsq producing strategy from env NSQ_PRODUCER_STRATEGY 'fallback' or 'random', default is 'fallback'

## v4.2.2 (Jun 14 2018)

* New package `nsqlbproducer` to produce among multiple nsqd instances
  Several strategies:
  * Random (default)
  * Fallback: always produce on the first host of the list, and fallback on other if it fails

## v4.0.0 (Mar 13 2018)

* `mongo/document` Add SortFields on all Find/Where methods

  Example of usage:

```
document.Where(ctx, col, query, data, document.SortField("-created_at"))
document.Find(ctx, col, query, id, data, document.SortField("-created_at"))
```

## v3.0.0 / v3.1.0 (Mar 8 2018)

* Uniformize API for Find/Sort with Unscoped versions and Sort versions

## v2.1.0 (Mar 5 2018)

* Add `mongo/document` package for CRUD operations against the database and
  base structure to embed in models

  ```
  type Model struct {
    document.Base `bson:",inline"`
  }

  type ParanoidModel struct {
    document.Paranoid `bson:",inline"`
  }
  ```

## v2.0.1

* Update rollbar plugin configuration, only skip 4 levels of stack

## v2.0.0

* Change interface of nsqproducer.Producer, add Stop method

## v1.1.2

### global

* No more reference to go-internal-tools

## v1.1.1

### nsqconsumer

* Create nsqconsumer.Error to handle when no retry should be done

## v1.1.0

### errors

* New errors.Notef and errors.Wrapf with context to let the error handling system
  read the context and its logger

## v1.0.2

### logger

* New API of logrus removed `logger.AddHook(hook)`, it is now `logger.Hooks.Add(hook)`

## v1.0.1

### logger

* Fix plugin system, use pointer instead of value for manager

## v1.0.0

### mongo (api change)

* Add a logger in initialization
