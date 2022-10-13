# Various Go Utils and Helpers

[ ![Travis Status for
Scalingo/go-utils](https://travis-ci.com/Scalingo/go-utils.svg?branch=master)](https://travis-ci.com/github/Scalingo/go-utils)

## Structure of This Repository

This repository is hosting modules, each of these modules are independant, they should all have their own:

* Dependencies (handled with go modules)
* `README.md`
* `CHANGELOG.md`
* Versioning through git tags. (Example for `etcd` → tag will look like `etcd/v1.0.0`)

## Release a New Version of a Module

Bump new version number in:

- `module/CHANGELOG.md`
- `module/README.md`

Commit, tag and create a new release:

```sh
module="XXX"
version="X.Y.Z"

git add ${module}/CHANGELOG.md ${module}/README.md
git commit -m "[${module}] Bump v${version}"
git tag ${module}/v${version}
git push origin master
git push --tags
```

## Use One Module in Your Project

With go modules, it's as easy as `go get github.com/Scalingo/go-utils/module`

For instance:

```
go get github.com/Scalingo/go-utils/logger
go get github.com/Scalingo/go-utils/logger@v1.0.0
go get github.com/Scalingo/go-utils/logger@<branch name>
```
