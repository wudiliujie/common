package console

import (
	"bufio"
	"github.com/wudiliujie/common/log"
	"os"
	"strings"
)

var CloseTag chan int32
var IsStartConsole = true

func Init(closeTag chan int32) {
	CloseTag = closeTag
	if IsStartConsole {
		go run()
	}

}

func Destroy() {

}

func run() {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Error("console ReadString is error: %v", err)
			continue
		}
		line = strings.TrimSuffix(line[:len(line)-1], "\r")

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		name := args[0]
		c := getCommand(name)
		if c == nil {
			log.Error("command not found, try `help` for help\r\n")
			continue
		}
		output := c.run(args[1:])
		if output != "" {
			log.Release("%v cmd run result: %v", name, output)
		}
	}
}
