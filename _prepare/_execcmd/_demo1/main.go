package main

import (
	"log"
	"os/exec"
)

func main() {

	/*
	通过 exec.Command 函数创建 Command对象 *exec.Cmd
	调用 *exec.Cmd 的 Run 方法执行命令
	通过 Run 方法可以执行，但是无法获取结果
	 */

	// Linux下
	// var cmd *exec.Cmd = exec.Command("/bin/bash", "-c", "echo 1;echo 2;")

    // Windows下
    var cmd *exec.Cmd = exec.Command(`C:\cygwin\bin\bash.exe`, `-c`, `echo 1`)

	err := cmd.Run()

	if err != nil {
		log.Println(err)
	}
}