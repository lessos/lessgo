## lessgo/logger
lessgo/logger is a Simplified implementation of Leveled logs ([google glog](https://github.com/google/glog)).

## Example

create a new test file main.go, and input the following codes

```go
package main

import (
    "github.com/lessos/lessgo/logger"
)

func main() {

    // API:Print
    logger.Print("info", "started")
    logger.Print("error", "the error code/message: ", 400, "/", "bad request")

    // API::Printf
    logger.Printf("error", "the error code/message: %d/%s", 400, "bad request")

    select {}
}
```

build the main.go file, run it and output the log into stderr console

```shell
go build main.go
./main -logtostderr=true
I 2016-04-30 22:24:26.463605 main.go:10] started
E 2016-04-30 22:24:26.463620 main.go:11] the error code/message: 400/bad request
E 2016-04-30 22:24:26.463628 main.go:14] the error code/message: 400/bad request
```

or run it and output the log into file
```shell
./main -log_dir="/var/log/"
```

the output file name will formated like this

```
/var/log/{program name}.{hostname}.{user name}.log.{level tag}.{date}-{time}.{pid}
```

## log levels
by default, the log levels ware defined as

<table>
<tr>
    <td>Level</td>
    <td>Tag</td>
    <td></td>
</tr>
<tr>
    <td>0</td>
    <td>debug</td>
    <td>Designates fine-grained informational events that are most useful to debug an application.</td>
</tr>
<tr>
    <td>1</td>
    <td>info</td>
    <td>Designates informational messages that highlight the progress of the application at coarse-grained level</td>
</tr>
<tr>
    <td>2</td>
    <td>warn</td>
    <td>Designates potentially harmful situations</td>
</tr>
<tr>
    <td>3</td>
    <td>error</td>
    <td>Designates error events that might still allow the application to continue running</td>
</tr>
<tr>
    <td>4</td>
    <td>fatal</td>
    <td>Designates very severe error events that will presumably lead the application to abort</td>
</tr>
</table>

You can also define your custom levels:
```go
logger.LevelConfig([]string{"warn", "error", "fatal"})
```

## Setting Flags

The flags influence lessos/logger's output behavior by passing on the command line. For example, if you want to turn the flag --logtostderr on, you can start your application with the following command line:

``` shell
./your_application --logtostderr=true
```

The following flags are most commonly used:
<table>
<tr>
    <td>flag</td>
    <td>type,default</td>
    <td></td>
</tr>
<tr>
    <td>log_dir</td>
    <td>string, default=""</td>
    <td></td>
</tr>
<tr>
    <td>logtostderr</td>
    <td>bool, default=false</td>
    <td>If output log messages to stderr</td>
</tr>
<tr>
    <td>minloglevel</td>
    <td>int, default=1(which is INFO)</td>
    <td>Log messages at or above this level. Again, the numbers of severity levels DEBUG, INFO, WARN, ERROR, and FATAL are 0, 1, 2, 3, and 4, respectively.</td>
</tr>
<tr>
    <td>logtolevels</td>
    <td>bool, default=false</td>
    <td>Output messages to multi leveled logfiles from minloglevel to the max, or output messages to the minloglevel logfile.</td>
</tr>
</table>
