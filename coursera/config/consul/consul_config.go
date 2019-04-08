package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

var (
	consulAddr = flag.String("addr", "192.168.99.100:32769", "consul addr (8500 in original consul)")

	consul          *consulapi.Client
	consulLastIndex uint64 = 0

	globalCfg   = make(map[string]string)
	globalCfgMu = &sync.RWMutex{}

	cfgPrefix      = "myapi/"
	prefixStripper = strings.NewReplacer(cfgPrefix, "")
)

// линтер ругается если используем базовые типы в Value контекста
// типа так безопаснее разграничивать
type key string

const configKey key = "configKey"

func configMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		globalCfgMu.RLock()
		localCfg := globalCfg
		globalCfgMu.RUnlock()

		ctx = context.WithValue(ctx,
			configKey,
			localCfg)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ConfigFromContext(ctx context.Context) (map[string]string, error) {
	cfg, ok := ctx.Value(configKey).(map[string]string)
	if !ok {
		return nil, fmt.Errorf("config not found")
	}
	return cfg, nil
}

func loadPostsHandle(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	localCfg, err := ConfigFromContext(ctx)

	if err != nil {
		http.Error(w, "no config!", 500)
		return
	}

	fmt.Fprintf(w, "localCfg version %v\n%#v\n", consulLastIndex, localCfg)
	fmt.Fprintln(w, "Request done")
}

func main() {
	flag.Parse()

	var err error
	config := consulapi.DefaultConfig()
	config.Address = *consulAddr
	consul, err = consulapi.NewClient(config)

	if err != nil {
		fmt.Println("consul error", err)
		return
	}

	loadConfig()
	go runConfigUpdater()

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", loadPostsHandle)

	siteHandler := configMiddleware(siteMux)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", siteHandler)
}

func runConfigUpdater() {
	ticker := time.Tick(3 * time.Second)
	for _ = range ticker {
		loadConfig()
	}
}

func loadConfig() {
	qo := &consulapi.QueryOptions{
		WaitIndex: consulLastIndex,
	}
	kvPairs, qm, err := consul.KV().List(cfgPrefix, qo)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("remote consulLastIndex", qm.LastIndex)
	if consulLastIndex == qm.LastIndex {
		fmt.Println("consulLastIndex not changed")
		return
	}

	newConfig := make(map[string]string)

	for idx, item := range kvPairs {
		if item.Key == cfgPrefix {
			continue
		}
		fmt.Printf("item[%d] %#v\n", idx, item)
		key := prefixStripper.Replace(item.Key)
		newConfig[key] = string(item.Value)
	}

	globalCfgMu.Lock()
	globalCfg = newConfig
	consulLastIndex = qm.LastIndex
	globalCfgMu.Unlock()

	fmt.Printf("config updated to version %v\n\t%#v\n\n", consulLastIndex, newConfig)
}
