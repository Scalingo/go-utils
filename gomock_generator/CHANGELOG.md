# Changelog

## To be Released

* Bump github.com/sirupsen/logrus from 1.7.0 to 1.8.0

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
