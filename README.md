# Various Go Utils and Helpers

## Structure of This Repository

This repository is hosting modules, each of these modules are independant, they should all have their own:

* Dependencies (handled with go modules)
* `README.md`
* `CHANGELOG.md`
* Versioning through git tags. (Example for `etcd` → tag will look like `etcd/v1.0.0`)

## Adding a New Module

- Follow the defined structure above
- Add the module on the dependabot config located in `.github` directory

## Release a New Version of a Module

Bump new version number in:

- `module/CHANGELOG.md`
- `module/README.md`

Commit, tag and create a new release:

```sh
module="XXX"
version="X.Y.Z"

git switch --create release/${module}/${version}
git add ${module}/CHANGELOG.md ${module}/README.md
git commit --message="[${module}] Bump v${version}"
git push --set-upstream origin release/${module}/${version}
gh pr create --reviewer=Scalingo/team-ist --title "$(git log -1 --pretty=%B)"
```

Once the pull request merged, you can tag the new release.

### Tag the New Release

```bash
git tag ${module}/v${version}
git push origin master ${module}/v${version}
```

## Use One Module in Your Project

With go modules, it's as easy as `go get github.com/Scalingo/go-utils/module`

For instance:

```
go get github.com/Scalingo/go-utils/logger
go get github.com/Scalingo/go-utils/logger@v1.0.0
go get github.com/Scalingo/go-utils/logger@<branch name>
```
