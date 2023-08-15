# slogevent

slogevent is a simplified API for integrating with slog.

Add `slogevent` as a dependency by running:

```sh
go get github.com/vikstrous/slogevent
```

Writing a handler for Go's standard library [log/slog](https://pkg.go.dev/log/slog) package is hard for good reasons. See the details for how to do that [here](https://golang.org/s/slog-handler-guide). However, sometimes there's a need to simply trigger an event when certain logs are written and this should be easy.

This package provides the missing glue to allow you to write log handlers like this:

```go
func EventHandler(ctx context.Context, e slogevent.Event) {
    if e.Level >= slog.LevelError {
        attrs, _ := json.Marshal(e.Attrs)
        SoundTheAlarm(e.Message, stirng(e.attrs))
    }
}
```

To use this custom even handler with `slog`, pass it into slog as a handler when creating your logger, like this:
```go
slogLogger := slog.New(slogevent.NewHandler(EventHandler, slog.NewTextHandler(os.Stderr, nil)))
```

slogevent acts as a wrapper for whatever handler you would normally use and fires your event before the next handler runs. It takes care of the annoying glue code needed to collect attrs and groups when `With()` and `WithGroup()` are used and converts all groups into a flag list of `slog.Attr`s that you can use without much effort.

This package is pre-v1 and the API may change. Please provide feedback about the API in the issues and let me know how I can make your slogging easier :muscle: