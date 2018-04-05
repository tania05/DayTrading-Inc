package main

import (
	"github.com/valyala/gorpc"
	"common/config"
	"fmt"
	"log"
	"common/logger"
	"os"
	"time"
	"os/exec"
)

func setupQuoteRpcs() {
	gorpc.RegisterType(&logger.UserCommandLog{})
	gorpc.RegisterType(&logger.QuoteServerLog{})
	gorpc.RegisterType(&logger.AccountTransactionLog{})
	gorpc.RegisterType(&logger.SystemEventLog{})
	gorpc.RegisterType(&logger.ErrorEventLog{})
	gorpc.RegisterType(&logger.DebugEventLog{})
	gorpc.RegisterType(&logger.DumpAdmin{})
	gorpc.RegisterType(&logger.DumplogCommand{})

	dispatcher := gorpc.NewDispatcher()
	dispatcher.AddFunc(logger.FUserCommandLog, func(v *logger.UserCommandLog) error {
		fmt.Println("Before timestamp")
		v.Timestamp = time.Now().UnixNano() / 1e6
		fmt.Println("After timestamp")
		return Log(*v)
	})
	dispatcher.AddFunc(logger.FAccountTransactionLog, func(v *logger.AccountTransactionLog) error {
		v.Timestamp = time.Now().UnixNano() / 1e6
		return Log(*v)
	})
	dispatcher.AddFunc(logger.FDebugEventLog, func(v *logger.DebugEventLog) error {
		v.Timestamp = time.Now().UnixNano() / 1e6
		return Log(*v)
	})
	dispatcher.AddFunc(logger.FErrorEventLog, func(v *logger.ErrorEventLog) error {
		v.Timestamp = time.Now().UnixNano() / 1e6
		return Log(*v)
	})
	dispatcher.AddFunc(logger.FQuoteServerLog, func(v *logger.QuoteServerLog) error {
		v.Timestamp = time.Now().UnixNano() / 1e6
		return Log(*v)
	})
	dispatcher.AddFunc(logger.FSystemEventLog, func(v *logger.SystemEventLog) error {
		v.Timestamp = time.Now().UnixNano() / 1e6
		return Log(*v)
	})

	dispatcher.AddFunc(logger.FAdminDumplog, func(v *logger.DumpAdmin) error {
		return AdminDumplog(*v)
	})

	dispatcher.AddFunc(logger.FDumplog, func(v *logger.DumplogCommand) error {
		return Dumplog(*v)
	})
	s := &gorpc.Server{
		Addr:    fmt.Sprintf(":%d", config.GlobalConfig.AuditServer.Port),
		Handler: dispatcher.NewHandlerFunc(),
	}

	if err := s.Serve(); err != nil {
		log.Fatalf("Can't start rpc server: %s", err)
		panic(err)
	}
}

var logChan chan []byte

func main() {


	logChan = make(chan []byte, 64)

	go func() {
		// fmt.Println("8")
		f, err := os.OpenFile("log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		// fmt.Println("7")
		if err != nil {
			panic(err)
		}
		for {
			// fmt.Println("5")
			bytes := <-logChan
			// fmt.Println("6")
			// fmt.Println(string(bytes))
			f.Write(bytes)
			f.Write([]byte("\n"))
		}
	}()

	setupQuoteRpcs()
}

func AdminDumplog(dump logger.DumpAdmin) error {
	fmt.Println("You have hit the admin dump log")
	cmd := exec.Command("cp", "log", dump.FileName)
	if err := cmd.Run(); err != nil {
		fmt.Println("Sorry couldn't make the dump log. Dying.")
		panic(err)
	}
	return nil
}

func Dumplog(dump logger.DumplogCommand) error {
	fmt.Println("You have normal dump log")
	cmd := exec.Command("/bin/sh", "-c", "grep \"<username>"+dump.UserId+"</username>\" log > \"" + dump.FileName + "\"")
	if err := cmd.Run(); err != nil {
		fmt.Println("Sorry couldn't make the dump log for the specified user. Dying.")
		panic(err)
	}
	return nil
}

func Log(msg logger.XmlLoggable) error {
	// fmt.Println("1")
	bytes, err := msg.AsXml()
	// fmt.Println("2")
	if err != nil {
		panic(err)
	}
	// fmt.Println("3")
	logChan <- bytes
	// fmt.Println("4")
	return nil
}
