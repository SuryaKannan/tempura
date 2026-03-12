package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var port int = 8732

type BatterRequest struct {
	FunctionName string         `json:"function_name"`
	InputHash    string         `json:"input_hash"`
	Input        map[string]any `json:"input"`
}

type BatterResponse struct {
	Status string `json:"status"`
	Output any    `json:"output"`
}

func init() {
	serve.Flags().IntVarP(&port, "port", "p", port, "port to serve tempura")
	rootCmd.AddCommand(serve)
}

func handleBatter(w http.ResponseWriter, r *http.Request) {

	var req BatterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error: cannot decode batter request")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(BatterResponse{Status: "new", Output: nil})

}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Healthy!")
}

func runServer(server *http.Server) {
	err := server.ListenAndServe()

	if err != nil {
		fmt.Fprintf(os.Stderr, "server did not start up: %v\n", err)
	}

}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "start your tempura server",
	Run: func(cmd *cobra.Command, args []string) {

		mux := http.NewServeMux()

		mux.HandleFunc("/batter", handleBatter)

		mux.HandleFunc("/health", handleHealth)

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

		err := server.Shutdown(ctx)

		if err != nil {
			fmt.Fprintf(os.Stderr, "server not shut down cleanly:%v\n", err)
		}
	},
}
