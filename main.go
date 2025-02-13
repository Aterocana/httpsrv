package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

const dateLayout = "2006/01/02-15:04:05"

var (
	version   = "unknown"
	buildDate = dateLayout
	gitCommit = "unknown"
)

// flags parses cmd line flags and it returns the path and the port to use.
func flags() (string, int) {
	var path string
	var port int
	var askHelp bool
	var askVersion bool

	flag.StringVar(&path, "path", "./", "the folder you want to serve. Default is current (./)")
	flag.IntVar(&port, "port", 0, "the port the server should listen on. Default is a random one")
	flag.BoolVar(&askHelp, "help", false, "show help")
	flag.BoolVar(&askVersion, "version", false, "show version")

	flag.Parse()
	if askHelp {
		flag.Usage()
		os.Exit(0)
	}
	if askVersion {
		buildDateTime, err := time.Parse(dateLayout, buildDate)
		if err != nil {
			buildDateTime = time.Time{}
		}
		fmt.Printf("%s (#%s) build on %s\n", version, gitCommit, buildDateTime)
		os.Exit(0)
	}
	return path, port
}

// server creates a [http.Server], which is returned, and it starts to accept requests in a goroutine.
func server(path string, port int) *http.Server {
	root := must(os.OpenRoot(path))
	listener := must(net.Listen("tcp", fmt.Sprintf(":%d", port)))
	protocols := &http.Protocols{}
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	srv := http.Server{
		Handler:   http.FileServerFS(root.FS()),
		Protocols: protocols,
	}
	addr := fmt.Sprintf("http://localhost:%d", listener.Addr().(*net.TCPAddr).Port)
	fmt.Printf("listening on %s\n", addr)
	openURL(addr)
	go must[*uint8](nil, srv.Serve(listener))
	return &srv
}

// must forces function to have no errors, otherwise it panics.
func must[T any](arg T, err error) T {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(0x01)
	}
	return arg
}

// openURL opens the url in the broswer.
// lurked from: https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		// Check if running under WSL
		if isWSL() {
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			cmd = "xdg-open"
			args = []string{url}
		}
	}
	if len(args) > 1 {
		// args[0] is used for 'start' command argument, to prevent issues with URLs starting with a quote
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}
	return exec.Command(cmd, args...).Start()
}

// isWSL checks if the program is running on Windows Subsystem for Linux
func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}

func main() {
	srv := server(flags())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc
	fmt.Println(" shutting down server")

	ctx, cancel := context.WithTimeoutCause(context.Background(), 5*time.Second, fmt.Errorf("timeout reached"))
	defer cancel()
	must[*uint8](nil, srv.Shutdown(ctx))
}
