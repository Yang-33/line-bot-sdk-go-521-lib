package gohttplib

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
)

var (
	DumpAllDeps bool

	cachedVersion string
	versionOnce   sync.Once
)

func libraryVersion() string {
	versionOnce.Do(func() {
		cachedVersion = fetchVersion()
	})
	return cachedVersion
}

func dumpAllDependencies() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("No build info found.")
		return
	}
	fmt.Println("=== All Dependencies ===")
	for _, dep := range info.Deps {
		fmt.Printf("- %s %s\n", dep.Path, dep.Version)
	}
	fmt.Println("========================")
}

func fetchVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown-ng"
	}
	for _, dep := range info.Deps {
		// Find the module path of your library by exact or partial match
		if strings.Contains(dep.Path, "github.com/Yang-33/line-bot-sdk-go-521-lib") {
			return dep.Version
		}
	}
	return "unknown-notfound"
}

func NewHttpClient() *http.Client {
	ua := "go-http-lib-example/" + libraryVersion()

	// Show all dependencies if a special parameter is true
	if DumpAllDeps {
		dumpAllDependencies()
	}

	return &http.Client{
		Transport: &userAgentTransport{
			base: http.DefaultTransport,
			ua:   ua,
		},
	}
}

type userAgentTransport struct {
	base http.RoundTripper
	ua   string
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.ua)
	return t.base.RoundTrip(req)
}
