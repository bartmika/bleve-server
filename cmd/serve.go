package cmd

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/bartmika/bleve-server/internal/controller"
	bleve_rpc "github.com/bartmika/bleve-server/internal/rpc_server"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the rpc server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doServe(cmd, args)
	},
}

func doServe(cmd *cobra.Command, args []string) {
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing address for", applicationAddress)
	log.Println("Executing from", path)
	log.Println("Opening indices at", applicationHomeDirectoryPath)

	tcpAddr, err := net.ResolveTCPAddr("tcp", applicationAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize our application controller.
	c, err := controller.New(applicationHomeDirectoryPath)
	if err != nil {
		log.Fatal(err)
	}

	// Integrate our controller with RPC server.
	r := bleve_rpc.RPC{
		Controller: c,
	}

	rpc.Register(&r)
	rpc.HandleHTTP()

	log.Println("RPC API was initialized.")
	l, e := net.ListenTCP("tcp", tcpAddr)
	if e != nil {
		l.Close()
		log.Fatal("RPC API failed to initialize:", e.Error())
	}

	log.Println("Started rpc service.")
	runMainRuntimeLoop(&r, l)
}

func runMainRuntimeLoop(r *bleve_rpc.RPC, l *net.TCPListener) {
	// The following code will attach a background handler to run when the
	// application detects a shutdown signal.
	// Special thanks via https://guzalexander.com/2017/05/31/gracefully-exit-server-in-go.html
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs // Block execution until signal from terminal gets triggered here.

		// Finish any RPC communication taking place at the moment before
		// shutting down the RPC server.
		l.Close()
	}()

	// Attach the following anonymous function to run on all cases (ex: panic,
	// termination signal, etc) so we can gracefully shutdown the service.
	defer func() {
		stopMainRuntimeLoop(r, l)
	}()

	// Safety net for 'too many open files' issue on legacy code.
	// Set a sane timeout duration for the http.DefaultClient, to ensure idle connections are terminated.
	// Reference: https://stackoverflow.com/questions/37454236/net-http-server-too-many-open-files-error
	http.DefaultClient.Timeout = time.Minute * 10

	// DEVELOPER NOTES:
	// If you get "too many open files" then please read the following article
	// http://publib.boulder.ibm.com/httpserv/ihsdiag/too_many_open_files.html
	// so you can run in your console:
	// $ ulimit -H -n 4096
	// $ ulimit -n 4096

	// Run the main loop blocking code.
	http.Serve(l, nil)
}

func stopMainRuntimeLoop(r *bleve_rpc.RPC, l *net.TCPListener) {
	log.Printf("Starting graceful shutdown now...")
	r.Controller.Close()
	l.Close()
	log.Printf("Terminated TCPListener.")
	log.Printf("Graceful shutdown finished.")
}
