package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/maxrasky/crema/internal/inmemory"
	"github.com/maxrasky/crema/internal/option"
	"github.com/spf13/viper"

	"github.com/maxrasky/crema/internal/service"
	"github.com/maxrasky/crema/internal/service/proto"
	"google.golang.org/grpc"

	"github.com/maxrasky/crema/internal/memcached"
)

func main() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigType("toml")
	v.SetConfigName("conf")
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var conf option.Option
	if err := v.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	var (
		store service.Storager
		err   error
	)

	if conf.Cache.Address == nil {
		store = inmemory.New()
	} else {
		store, err = memcached.New(*conf.Cache.Address)
		if err != nil {
			log.Fatal(err)
		}
	}

	listener, err := net.Listen("tcp", conf.App.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = listener.Close()
	}()

	srv := service.New(store)

	server := grpc.NewServer()
	proto.RegisterServiceServer(server, srv)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	done := make(chan struct{})
	log.Println("starting app")

	go func() {
		if err = server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		<-ctx.Done()
		server.Stop()
		close(done)
	}()

	<-done
	log.Println("stopping app")
	time.Sleep(time.Second)
}
