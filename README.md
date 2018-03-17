# swarmprom
[![Documentation](https://godoc.org/github.com/faryon93/swarmprom?status.svg)](http://godoc.org/godoc.org/github.com/faryon93/swarmprom)
[![Go Report Card](https://goreportcard.com/badge/github.com/faryon93/swarmprom)](https://goreportcard.com/report/github.com/faryon93/swarmprom)

This library provides a secure way to expose the [Prometheus](https://prometheus.io) metrics http handler in a docker swarm environment.
The access to the handler is only allowed for containers beloging to the configured swarm service.
The relevant ip addresses are obtained from the swarm dns service discovery mechanism.

## Usage
Use the prometheus primitives to setup your monitoring environment. Just use `swarmprom.Handler()` as the handler for your HTTP server instance.
As argument you have to supply the name of the swarm service prometheus is reachable by.

    import (
        "github.com/faryon93/swarmprom"
        "github.com/prometheus/client_golang/prometheus"
    )

    func main() {
        // register prometheus metrics
        var counter = prometheus.NewCounter(prometheus.CounterOpts{
            Namespace: "namespace",
            Subsystem: "subsystem",
            Name:      "test_total",
            Help:      "Total number of tests.",
        })
        prometheus.MustRegister(counter)

        // expose http endpoint
        http.Handle("/metrics", swarmprom.Handler("prometheus"))
        http.ListenAndServe(":8080", nil)
    }

## Customization
To define a custom handler function in case the access get rejected use `swarmprom.SetRejectHandler()`:

    swarmprom.SetRejectHandler(func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "Nice night for a walk!", http.StatusNotFound)
    })

To use a custom logging entry use `swarmprom.SetLogger()`:

    swarmprom.SetLogger(logrus.WithField("handler", "prometheus"))