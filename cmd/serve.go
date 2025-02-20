package cmd

import(
	"fmt"

	"cross_tech/server"
)

func serve() {
	srv, err := server.New()
	if err != nil {
		fmt.Printf("failed, creating a server: %v", err)
		return
	}
	err = srv.Start(srvPort)
	if err != nil {
		fmt.Printf("failed, starting a server: %v", err)
		return
	}
}