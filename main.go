package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/didip/shawty/handlers"
	"github.com/didip/shawty/storages"
	"github.com/eyedeekay/onramp"
)

var (
	i2p  = flag.Bool("i2p", false, "Use I2P")
	host = flag.String("host", "127.0.0.1", "Host(ignored if using I2P)")
	port = flag.String("port", defport(), "Port(ignored if using I2P)")
)

func defport() string{
	p:= os.Getenv("PORT")
	if p != "" {
		return p
	}
	return "8080"
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	dir, _ := os.UserHomeDir()
	storage := &storages.Filesystem{}
	err := storage.Init(filepath.Join(dir, "shawty"))
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", handlers.EncodeHandler(storage))
	http.Handle("/dec/", handlers.DecodeHandler(storage))
	http.Handle("/red/", handlers.RedirectHandler(storage))

	if *i2p {
		garlic, err := onramp.NewGarlic()
		if err != nil {
			log.Fatal(err)
		}
		listener, err := garlic.Listen()
		if err != nil {
			log.Fatal(err)
		}
		err = http.Serve(listener, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = http.ListenAndServe(*host+":"+*port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

}
