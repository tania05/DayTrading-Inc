package main

import (
	"github.com/valyala/gorpc"
	. "common/rpc/quotestructs"
	"common/database"
	"common/config"
	"fmt"
	"log"
	"sync"
	"time"
	"strings"
	"common/quote"
	"common/context"
	"common/logger"
	"database/sql"
)


func setupQuoteRpcs() {
	auditCount = 1 

	gorpc.RegisterType(&logger.UserCommandLog{})
	gorpc.RegisterType(&logger.QuoteServerLog{})
	gorpc.RegisterType(&logger.AccountTransactionLog{})
	gorpc.RegisterType(&logger.SystemEventLog{})
	gorpc.RegisterType(&logger.ErrorEventLog{})
	gorpc.RegisterType(&logger.DebugEventLog{})

	dispatcher = gorpc.NewDispatcher()
	dispatcher.AddFunc(logger.FLog, Log)

	s := &gorpc.Server{
		Addr: fmt.Sprintf(":%d", config.GlobalConfig.AuditServer.Port),
		Handler: d.NewHandlerFunc(),
	}

	if err := s.Serve(); err != nil {
		log.Fatalf("Can't start rpc server: %s", err)
		panic(err)
	}
}

func main() {

	setupQuoteRpcs()

	logChan = make(chan []byte, 64)
	
	go func() {
		f, err := os.OpenFile("log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		for {
			bytes := <-logChan
			f.Write(bytes)
			f.Write([]byte("\n"))
		}
	}()
}

func Log(msg *XmlLoggable) {
	bytes, err := msg.asXml()
	if err != nil {
		panic(err)
	}
	logChan <- bytes
}
