package httpsrv

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Server struct {
	httpSrv  *http.Server
	listener net.Listener
	out      io.Writer
}

func logRequest(out io.Writer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(out, "[%s] %s\n", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func New(opts ...Options) (*Server, error) {
	cfg := buildConfig(opts...)
	return newSrv(cfg)
}

func newSrv(cfg *config) (*Server, error) {
	root, err := os.OpenRoot(cfg.path)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", cfg.path, err)
	}

	protocols := &http.Protocols{}
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.port))
	if err != nil {
		return nil, fmt.Errorf("could not listen to port %d: %w", cfg.port, err)
	}

	srv := http.Server{
		Handler:   logRequest(os.Stdout, http.FileServerFS(root.FS())),
		Protocols: protocols,
	}

	return &Server{
		httpSrv:  &srv,
		listener: listener,
	}, nil
}

func (srv *Server) Open() error {
	addr := fmt.Sprintf("http://localhost:%d", srv.listener.Addr().(*net.TCPAddr).Port)
	go srv.httpSrv.Serve(srv.listener)
	openURL(addr)
	fmt.Printf("listening on %s\n", addr)
	return nil
}

func (srv *Server) Close(ctx context.Context) error {
	return srv.httpSrv.Shutdown(ctx)
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
