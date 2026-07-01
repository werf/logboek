## logboek — a library for structured and informative output

### Logger, channels, and streams

<img src="https://github.com/werf/logboek/blob/master/logger.png?raw=true">

When creating a logger, you need to specify the **streams**, `OutStream` and `ErrStream`, which must adhere to the `io.Writer` interface (this can be a file, standard output streams, a buffer, or any custom implementation):

```go
import "github.com/werf/logboek"

// NewLogger(outStream, errStream io.Writer) *Logger
l := logboek.NewLogger(os.Stdout, os.Stderr)
```

Stream settings allow you to define formatting parameters such as prefix and tag, as well as various modes of operation. These settings apply to both `OutStream` and `ErrStream`, and consequently to all channels that will be discussed later.

```go
l.Streams()
```

The logger is connected to the **log channels** `Error`, `Warn`, `Default`, `Info`, and `Debug`. When using the `Error` and `Warn` channels, all messages are written to `ErrStream`, while for the others, they go to `OutStream`.

Log channels allow you to organize the output for various application modes (verbose and debug modes), branch execution, and control flow depending on the active channel (activating a channel also triggers output to lower-priority channels):

```go
import (
    "github.com/werf/logboek"
    "github.com/werf/logboek/pkg/level"
)

switch mode {
case "verbose":
    l.SetAcceptedLevel(level.Info)
case "debug":
    l.SetAcceptedLevel(level.Debug)
case "quiet":
    l.SetAcceptedLevel(level.Error)  
}

...

if l.Debug().IsAccepted() {
  ... // do and print something special
}
```

If channels are not required, you can simply use the `Default` channel, whose methods are available at the top level of the logger:

```go
l.LogLn() // l.Default().LogLn()
l.LogF()  // l.Default().LogF()
...
```

<!---
- Terminal width
- Inherited settings
- Proxy
-->

### Default logger

By default, the library initializes the `DefaultLogger` with preset streams `os.Stdout` and `os.Stderr`. You can interact with the logger using the instance itself or the high-level library functions that correspond to all available logger methods:

```go
import "github.com/werf/logboek"

logboek.DefaultLogger()

logboek.Default() // logboek.DefaultLogger().Default()
logboek.LogLn()   // logboek.DefaultLogger().LogLn()
logboek.Streams() // logboek.DefaultLogger().Streams()
...
```

### Using logboek with context.Context

Logboek can propagate a logger through `context.Context`:

```go
ctx := logboek.NewContext(context.Background(), myLogger)

// Inside a helper function
logboek.Context(ctx).LogLn("hello from helper")
```

`logboek.Context(ctx)` is **safe**: if the provided context does not hold a logger (or is `nil`/`context.Background()`), it automatically falls back to `logboek.DefaultLogger()`.

If you need the original strict behaviour (panic when no logger found), call `logboek.MustContext(ctx)` instead.

```go
// Will panic if ctx has no logger
l := logboek.MustContext(ctx)
```
<!---
## Logging Methods

<img align="right" src="https://github.com/werf/logboek/blob/master/logboek.png?raw=true">
-->

<!---
## Processes and blocks
## Prefix and tag
## Modes
- isMuted                            
- isStyleEnabled                     
- isLineWrappingEnabled              
- isProxyStreamDataFormattingEnabled 
- isGitlabCollapsibleSectionsEnabled 
- isPrefixWithTimeEnabled            
- isLogProcessBorderEnabled 
## Using in external libraries
## Using in go-routines
-->
