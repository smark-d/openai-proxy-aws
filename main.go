package main

import (
	"github.com/smark-d/openai-proxy-aws/cmd"
	"github.com/smark-d/openai-proxy-aws/server/comm"
)

func main() {
	comm.InitConfig()
	cmd.RootCommand.Execute()
}
