# Changelog

## To be Released

* chore(go): use go 1.17
* build(deps): bump github.com/urfave/cli from 1.22.5 to 1.22.9

## v1.3.0

* Bump github.com/sirupsen/logrus from 1.7.0 to 1.8.1
* Bump go version to 1.16 and replace ioutil by io/os [#198](https://github.com/Scalingo/go-utils/pull/198)
* Bump github.com/go-utils/logger from v1.0.0 to v1.1.0

## v1.2.2

* Set the version in the final binary
* Update installation help in README.md

## v1.2.1

* Add notion of `BaseDirectory` when go modules are used, which is sometimes different from `BasePackage`:

    ```
    {
      "BaseDirectory": "github.com/Scalingo/go-scalingo",
      "BasePackage": "github.com/Scalingo/go-scalingo/v4"
    }
    ```

## v1.2.0

* Embedded interfaces are not throwing an error anymore.
  Known caveat: mocks with embedded interfaces will be regenerate at each execution.

## v1.1.0

* Pretty print JSON for mocks_sig.json
* Correctly handle go modules to read data at the right location
* Fix regression: when SrcPackage was not defined, an error was spawn

## v1.0.0 - v1.0.1

* Initial breakdown of go-utils into subpackages
