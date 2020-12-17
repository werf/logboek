logboek — библиотека для организации структурированного и информативного вывода
_______________________________________________________________________________

## Логер, каналы логирования и потоки вывода

<img src="https://github.com/werf/logboek/blob/master/logger.png?raw=true">

При создании логера необходимо указать **стримы**, `OutStream` и `ErrStream`, которые должны подходить под интерфейс `io.Writer` (это может быть файл, стандартные потоки вывода, буфер или произвольная имплементация):

```go
import "github.com/werf/logboek"

// NewLogger(outStream, errStream io.Writer) *Logger
l := logboek.NewLogger(os.Stdout, os.Stderr)
```

Настройки стримов позволяют задать параметры оформления, такие как префикс и тег, а также различные режимы работы. Они являются общими для `OutStream` и `ErrStream`, и, соответственно, для всех каналов, о которых пойдёт речь далее.

```go
l.Streams()
```

Логер связан с **каналами логирования** `Error`, `Warn`, `Default`, `Info` и `Debug`. При использовании каналов `Error` и `Warn` все сообщения пишутся в `ErrStream`, а в случае с остальными в `OutStream`.

Каналы логирования позволяют организовать вывод для различных режимов работы приложения (подробный и дебаг режимы), ветвление и выполнение кода в случае того или иного активного канала (активация канала также включает вывод в нижестоящих каналах по приоритету):

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

Если каналы не требуются, то можно ограничиться использованием канала `Default`, методы которого доступны на верхнем уровне логера:

```go
l.LogLn() // l.Default().LogLn()
l.LogF()  // l.Default().LogF()
...
```

<!---
- Ширина терминала
- Наследование настроек
- Прокси
-->

## Default Logger

По умолчанию библиотека инициализирует `DefaultLogger` с предустановленными стримами `os.Stdout` и `os.Stderr`. Для работы с логером можно использовать сам экземпляр или верхнеуровневые функции библиотеки, которые соответствуют всем доступным методам логера:

```go
import "github.com/werf/logboek"

logboek.DefaultLogger()

logboek.Default() // logboek.DefaultLogger().Default()
logboek.LogLn()   // logboek.DefaultLogger().LogLn()
logboek.Streams() // logboek.DefaultLogger().Streams()
...
```

<!---
## Методы логирования

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
