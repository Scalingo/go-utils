# Various Go Utils and Helpers

[ ![Codeship Status for
Scalingo/go-utils](https://app.codeship.com/projects/af479f60-02c1-0136-d485-6637162e76f3/status?branch=master)](https://app.codeship.com/projects/280142)

## Release a New Version

Bump new version number in:

- `CHANGELOG.md`
- `README.md`

Commit, tag and create a new release:

```sh
git add CHANGELOG.md README.md
git commit -m "Bump v7.1.0"
git tag v7.1.0
git push origin master
git push --tags
hub release create v7.1.0
```

The title of the release should be the version number and the text of the
release is the same as the changelog.
