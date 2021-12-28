# Package `io` v1.0.0

This package aims at proposing a customizable `io.Copy` method. It introduces the `Copier` struct.

## Usage

Constructor:

```
NewCopier(...CopierOpt)
```

List of Options:

* `WithBufferSize(int64)`
* `WithNoDiskCacheRead`
* `WithNoDiskCacheWrite`
* `WithNoDiskCache` (helper wrapping the two above)
