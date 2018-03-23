package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/brian-armstrong/gpio"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const version string = "0.1"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("listen-address", ":9500", "Address on which to expose metrics.")
	metricsPath   = flag.String("path", "/metrics", "Path under which to expose metrics.")
	wattPerPulse  = flag.Int("watt", 800, "Watt per pulse")
	pin           = flag.Int("pin", 3, "Pin connected to s0 signal")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: s0_exporter [ ... ]\n\nParameters:")
		fmt.Println()
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	go startListener()
	startServer()
}

func printVersion() {
	fmt.Println("s0_exporter")
	fmt.Printf("Version: %s\n", version)
}

func startListener() {
	watcher := gpio.NewWatcher()
	watcher.AddPin(uint(*pin))
	defer watcher.Close()
	sum.count = 0
	for {
		_, value := watcher.Watch()
		if value != 1 {
			continue
		}
		CounterUp()
	}
}

func startServer() {
	log.Infof("Starting s0 exporter (Version: %s)\n", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>s0 Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>s0 Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			<h2>More information:</h2>
			<p><a href="https://github.com/l3akage/s0_exporter">github.com/l3akage/s0_exporter</a></p>
			</body>
			</html>`))
	})
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	log.Infof("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(&s0Collector{})

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)
}
