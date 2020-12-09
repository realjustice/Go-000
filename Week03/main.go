package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
import _ "net/http/pprof"

func main() {
	fmt.Printf("PID :%d\n", os.Getpid())
	done := make(chan error, 2)
	stop := make(chan struct{})
	sig := make(chan os.Signal)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		done <- Server(":8080", nil, stop, sig)
	}()

	go func() {
		done <- ServerDebug(":8081", nil, stop)
	}()

	var isStop = false
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			fmt.Println(err)
			if !isStop {
				isStop = true
				close(stop)
			}
		}
	}
}

func Server(addr string, handler http.Handler, stop <-chan struct{}, signal chan os.Signal) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		select {
		case sig := <-signal:
			if sig == syscall.SIGTERM || sig == syscall.SIGINT {
				fmt.Println("Program Exit...", sig)
				_ = s.Shutdown(context.Background())
			}
		case <-stop:
			_ = s.Shutdown(context.Background())
			fmt.Println("Server shutdown!")
		}
	}()
	return s.ListenAndServe()
}

func ServerDebug(addr string, handler http.Handler, stop <-chan struct{}) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-stop
		fmt.Println("Server debug shutdown!")
		_ = s.Shutdown(context.Background())
	}()
	return s.ListenAndServe()
}
