`internal/` is a special directory name recognised by the `go` tool which will prevent one package from being imported by another unless both share a common ancestor.

Packages within an `internal/` directory are therefore said to be **internal packages**.
