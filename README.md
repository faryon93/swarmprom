# swarmprom
This library provides a secure way to expose the [Prometheus](https://prometheus.io) metrics http handler in a docker swarm environment.
The access to the handler is only allowed for containers beloging to the configured swarm service.
The relevant ip addresses are obtained from the swarm dns service discovery mechanism.

## Usage
Use the prometheus primitives to setup your monitoring environment. Just use `swarmprom.Handler()` as the handler for your HTTP server instance.

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
        http.Handle("/metrics", swarmprom.Handler())
        http.ListenAndServe(":8080", nil)
    }

## Customization
To define a custom handler function in case the access get rejected use `swarmprom.SetRejectHandler()`:

    swarmprom.SetRejectHandler(func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "Nice night for a walk!", http.StatusNotFound)
    })
