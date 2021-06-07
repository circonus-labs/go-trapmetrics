# Package `go-trapmetrics`

A simplified, streamlined, and updated replacement for `circonus-gometrics`. Covers basic functionality for collecting and submitting
metrics to Circonus.

```go
package main

import (
    "fmt"

    apiclient "github.com/circonus-labs/go-apiclient"
    trapcheck "github.com/circonus-labs/go-trapcheck"
    trapmetrics "github.com/circonus-labs/go-trapmetrics"
)

func main() {

    client, err := apiclient.New(&apiclient.Config{/* ...settings... */})
    if err != nil {
        panic(err)
    }

    check, err := trapcheck.New(&trapcheck.Config{Client: client, /* ...other settings... */})
    if err != nil {
        panic(err)
    }

    metrics, err := trapmetrics.New(&trapmetrics.Config{Check: check, /* ...other settings...*/})
    if err != nil {
        panic(err)
    }

    // NOTE: gauges and text take an optional timestamp for the sample, pass nil to use current time
    ts := time.Now()    
    metrics.GaugeSet("gauge",trapmetrics.Tags{{Category:"a",Value:"b"}},123,&ts)
    metrics.TextSet("text",nil,"some text",nil)

    // Counters and histograms will apply a timestamp at the time of flushing as they only contain
    // one stample and it is mutable until flush time.
    metrics.CounterIncrement("counter",trapmetrics.Tags{{Cateogry:"a",Value:"b"}})
    metrics.HistogramRecordValue("histogram",nil,27)
    metrics.CumulativeHistogramRecordCountForValue("cumulative_histogram",nil,128,3.6)

    result, err := metrics.Flush()
    if err != nil {
        panic(err)
    }

    fmt.Println(result)
}
```
