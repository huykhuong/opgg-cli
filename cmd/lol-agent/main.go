package main

import (
	"fmt"

	"github.com/huykhuong/lol/cmd/lol-agent/commands"
)

func main() {
	commands.Execute()
	fmt.Println("🚀 LoL Insight started")
}