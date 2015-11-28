package eye

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var Env = make(map[string]string)

func NewEnv(filename string) {
	if filename == "" {
		panic("config file is empty")
	}
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(fmt.Sprintf("open conf panic:%s", filename))
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) <= 0 {
			continue
		}
		env := strings.Split(line, "=")
		if len(env) != 2 {
			panic("conf env error")
		}
		Env[env[0]] = env[1]
	}
	return
}

func init() {
	NewEnv(fmt.Sprintf("%s/conf/env.conf", os.Getenv("GOPATH")))
}
