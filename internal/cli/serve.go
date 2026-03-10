package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var port int = 8732

func init() {
	serve.Flags().IntVarP(&port, "port", "p", port, "port to serve tempura")
	rootCmd.AddCommand(serve)
}

func handleFry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "tempura is running!")
}

func runServer(server *http.Server) {
	server.ListenAndServe()
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "start your tempura server",
	Run: func(cmd *cobra.Command, args []string) {

		mux := http.NewServeMux()

		mux.HandleFunc("/", handleFry)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		}

		color.Green("\nserving Tempura on port %d!\n", port)

		go runServer(server)

		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		<-sigChan

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		server.Shutdown(ctx)
	},
}
