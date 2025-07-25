package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	gorilla "github.com/mtlynch/gorilla-handlers"

	"github.com/mtlynch/picoshare/v2/garbagecollect"
	"github.com/mtlynch/picoshare/v2/handlers"
	"github.com/mtlynch/picoshare/v2/handlers/auth"
	"github.com/mtlynch/picoshare/v2/handlers/auth/authentik"
	"github.com/mtlynch/picoshare/v2/handlers/auth/shared_secret"
	"github.com/mtlynch/picoshare/v2/space"
	"github.com/mtlynch/picoshare/v2/store/sqlite"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Print("starting picoshare server")

	dbPath := flag.String("db", "data/store.db", "path to database")
	flag.Parse()

	authenticator, err := createAuthenticator()
	if err != nil {
		log.Fatalf("failed to create authenticator: %v", err)
	}

	dbDir := filepath.Dir(*dbPath)

	ensureDirExists(dbDir)

	store := sqlite.New(*dbPath, isLitestreamEnabled())

	spaceChecker := space.NewChecker(*dbPath, &store)

	collector := garbagecollect.NewCollector(store)
	gc := garbagecollect.NewScheduler(&collector, 7*time.Hour)
	gc.StartAsync()

	clock := handlers.NewClock()

	server := handlers.New(authenticator, &store, spaceChecker, &collector, &clock)

	h := gorilla.LoggingHandler(os.Stdout, server.Router())
	if os.Getenv("PS_BEHIND_PROXY") != "" {
		h = gorilla.ProxyIPHeadersHandler(h)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4001"
	}

	stop := setupSignalHandler()
	httpSrv := http.Server{Addr: fmt.Sprintf(":%s", port), Handler: h}
	go func() {
		log.Printf("listening on %s", port)
		log.Printf("http server exit: %s", httpSrv.ListenAndServe())
	}()
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Fatal(httpSrv.Shutdown(ctx))
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("missing required environment variable: %s", key))
	}
	return val
}

func ensureDirExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func isLitestreamEnabled() bool {
	return os.Getenv("LITESTREAM_BUCKET") != ""
}

func createAuthenticator() (handlers.Authenticator, error) {
	authProvider := getEnvWithDefault("PS_AUTH_PROVIDER", "shared_secret")
	
	switch authProvider {
	case "shared_secret":
		return createSharedSecretAuth()
	case "authentik":
		return createAuthentikAuth()
	default:
		return createSharedSecretAuth() // fallback
	}
}

func createSharedSecretAuth() (handlers.Authenticator, error) {
	sharedSecret := os.Getenv("PS_SHARED_SECRET")
	if sharedSecret == "" {
		return nil, fmt.Errorf("PS_SHARED_SECRET is required for shared_secret auth")
	}

	adapter, err := shared_secret.NewAdapter(sharedSecret)
	if err != nil {
		return nil, err
	}

	multiAuth := auth.NewMultiProviderAuthenticator("shared_secret")
	multiAuth.RegisterProvider("shared_secret", adapter)
	
	return multiAuth, nil
}

func createAuthentikAuth() (handlers.Authenticator, error) {
	config := authentik.Config{
		AuthentikURL:   requireEnv("PS_AUTHENTIK_URL"),
		ClientID:       requireEnv("PS_AUTHENTIK_CLIENT_ID"),
		ClientSecret:   requireEnv("PS_AUTHENTIK_CLIENT_SECRET"),
		RedirectURI:    requireEnv("PS_AUTHENTIK_REDIRECT_URI"),
		ReverseProxy:   getEnvWithDefault("PS_AUTHENTIK_REVERSE_PROXY", "false") == "true",
	}

	// Parse trusted proxies if specified
	if trustedProxiesEnv := os.Getenv("PS_AUTHENTIK_TRUSTED_PROXIES"); trustedProxiesEnv != "" {
		config.TrustedProxies = strings.Split(trustedProxiesEnv, ",")
		for i, proxy := range config.TrustedProxies {
			config.TrustedProxies[i] = strings.TrimSpace(proxy)
		}
	}

	provider := authentik.New(config)
	
	multiAuth := auth.NewMultiProviderAuthenticator("authentik")
	multiAuth.RegisterProvider("authentik", provider)
	
	return multiAuth, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupSignalHandler() <-chan struct{} {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()
	return stop
}
