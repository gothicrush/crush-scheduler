package main

import (
	"context"
	"log"
	"os/exec"
	"time"
)

func main() {

	ctx, cancelFuntion := context.WithCancel(context.TODO())

	dataChan := make(chan []byte, 1000)

	go func() {
		time.Sleep(5 * time.Second)
		var cmd *exec.Cmd = exec.CommandContext(ctx,`C:\cygwin\bin\bash.exe`,`-c`,`ls -al`)

		data, err := cmd.CombinedOutput()

		if err != nil {
			log.Println(err)
		}

		dataChan <- data
	}()

	time.Sleep(2 * time.Second)

	cancelFuntion()

	log.Println(string(<-dataChan))
}
