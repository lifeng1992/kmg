package goCmd

import (
	"os"

	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgConsole"
)

// run go command in current project
// 1.go build -i github.com/xxx/xxx use to get fastest speed of build.
// 2.try remove pkg directory if you found you change is ignore.
func GoCommand() {
	kmgc, err := kmgConfig.LoadEnvFromWd()
	kmgConsole.ExitOnErr(err)
	err = kmgCmd.CmdSlice(append([]string{"go"}, os.Args[1:]...)).
		MustSetEnv("GOPATH", kmgc.GOPATHToString()).StdioRun()
	kmgConsole.ExitOnErr(err)
}
