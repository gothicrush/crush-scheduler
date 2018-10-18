package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {

	// 创建命令对象
	var cmd *exec.Cmd = exec.Command(`C:\cygwin\bin\bash.exe`,`-c`,`sleep 5; ls -l`)

	// 执行命令，并获取执行结果
	result,err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
		return
	}

	// 打印命令执行结果
	fmt.Println(string(result))
}
