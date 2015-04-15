package kmgConsole

import (
	"fmt"
	"os"
)

//仅限于命令使用
func ExitOnErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

//仅限于命令使用
func ExitOnStderr(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}