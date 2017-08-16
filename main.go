package main

import (
	"flag"
	"net"
	"os"

	"go.uber.org/zap"

	"github.com/ChoTotOSS/fluent2gelf/agent"
	"github.com/ChoTotOSS/fluent2gelf/fluentd"
	"github.com/duythinht/zaptor"
)

var logger = zaptor.Default()

func main() {

	var c string

	flag.StringVar(&c, "c", "sample.yml", "config file")
	flag.Parse()

	f, err := os.Open(c)

	if err != nil {
		logger.Error("Open config", zap.Error(err))
	}

	agentStore := agent.AgentStoreLoad(f)

	doneList := make([](chan bool), len(agentStore.AgentList))

	for i, agent := range agentStore.AgentList {
		done := make(chan bool)
		doneList[i] = done
		go agent.Run(done)
	}

	serv, err := net.Listen("tcp", ":24224")
	checkError(err)

	logger.Info("Forward server was started 24224")
	for {
		conn, err := serv.Accept()
		logger.Debug("New client connected", zap.Any("conn", conn))
		checkError(err)
		go fluentd.ForwardHandle(conn, agentStore)
	}
}

func checkError(err error) {
	if err != nil {
		logger.Fatal("main#main", zap.Error(err))
		os.Exit(1)
	}
}
