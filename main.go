package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

func raise(err error) {
	if err != nil {
		panic(err)
	}
}

func must[T any](v T, err error) T {
	raise(err)
	return v
}

func Filter[T []U, U any](values T, f func(u U) bool) T {
	ret := make(T, 0, len(values))
	for _, value := range values {
		if f(value) {
			ret = append(ret, value)
		}
	}
	return ret
}

func main() {
	user := flag.String("user", "", "process's user")
	name := flag.String("name", "", "process's name like")
	flag.Parse()

	processes, err := process.Processes()
	raise(err)
	if *user != "" {
		processes = Filter(processes, func(p *process.Process) bool {
			userName := must(p.Username())
			return userName == *user
		})

	}
	if *name != "" {
		processes = Filter(processes, func(p *process.Process) bool {
			pn := must(p.Name())
			return strings.Contains(strings.ToLower(pn), *name)
		})

	}
	for _, p := range processes {
		fmt.Println(p.Pid, must(p.Name()), must(p.Username()), must(p.Cmdline()))
	}
}
