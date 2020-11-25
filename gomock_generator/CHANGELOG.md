## To Release

## v1.2.0

* Add notion of `BaseDirectory` when go modules are used, which is sometimes different from `BasePackage`:

    ```
    {
      "BaseDirectory": "github.com/Scalingo/go-scalingo",
      "BasePackage": "github.com/Scalingo/go-scalingo/v4"
    }
    ```

## v1.1.0

* Pretty print JSON for mocks_sig.json
* Correctly handle go modules to read data at the right location
* Fix regression: when SrcPackage was not defined, an error was spawn

## v1.0.0 - v1.0.1

* Initial breakdown of go-utils into subpackages
