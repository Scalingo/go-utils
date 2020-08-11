difflib
==========

Go-difflib is a partial port of
[Python 3 difflib package](https://docs.python.org/3/library/difflib.html).
Its main goal is to make unified and context diff available in pure Go,
mostly for testing purposes.

The following class and functions (and related tests) have been ported:

* `SequenceMatcher`: Struct containing some methods to compares sequence of strings.

* `unified_diff()`:
Compare A and B (array of strings) and return a generator generating the
output lines in `unified diff format`.
Unified diffs are a compact way of showing just the lines that have changed plus a few lines of context. The changes are shown in an inline style (instead of separate before/after blocks).

* `context_diff()`:
Compare A and B (array of strings) and return a generator generating the
output lines in `context diff format`.
Context diffs are a compact way of showing just the lines that have changed plus a few lines of context. The changes are shown in a before/after style.

> You can also have the output in colors by enabling the `WithColors` parameter (only for unified diff for the moment).

## Installation

```bash
$ go get github.com/Scalingo/go-utils/difflib
```

## Quick Start

Diffs are configured with Unified (or ContextDiff) structures, and can
be output to an io.Writer or returned as a string.

```Go
diff := difflib.UnifiedDiff{
    A:        difflib.SplitLines("foo\nbar\n"),
    B:        difflib.SplitLines("foo\nbaz\n"),
    FromFile: "Original",
    ToFile:   "Current",
    Context:  3,
    WithColors: false
}
text, _ := difflib.GetUnifiedDiffString(diff)
fmt.Printf(text)
```

Would output:

```shell
--- Original
+++ Current
@@ -1,3 +1,3 @@
 foo
-bar
+baz
```
