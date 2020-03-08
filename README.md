# `health-go`

![Go](https://github.com/nelkinda/health-go/workflows/Go/badge.svg)

Golang implementation of the upcoming [IETF RFC Health Check Response Format](https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html) for HTTP APIs.

## Usage
In your go program,

1. Create the health Handler.
1. Add the handler to your mux/server.

```go
package main

import (
	"github.com/nelkinda/health-go"
	"net/http"
)

func main() {
	// 1. Create the health Handler.
	h := health.New(health.Health{Version: "1", ReleaseID: "1.0.0-SNAPSHOT"}) 

	// 2. Add the handler to your mux/server.
	http.HandleFunc("/health", h.Handler)
	
	// 3. Start your server.
	http.ListenAndServe(":80", nil)
}
```

## Providing Checks
If is possible to provide checks.
This library comes with the following checks predefined:
- system uptime
- process uptime
- mongodb health
- SendGrid health
- sysinfo information (CPU Utilization, RAM, uptime, number of processes)

You can add any implementation of `ChecksProvider` to the varargs list of `health.New()`.

```go
package main

import (
	"context"
	"github.com/nelkinda/health-go"
	"github.com/nelkinda/health-go/checks/uptime"
	"github.com/nelkinda/health-go/checks/sysinfo"
	"github.com/nelkinda/health-go/checks/mongodb"
	"github.com/nelkinda/health-go/checks/sendgrid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func main() {
	url := "mongodb://127.0.0.1:27017"
	client, _ := mongo.NewClient(options.Client().ApplyURI(url))
	_ = client.Connect(context.Background())
	h := health.New(
		health.Health{
			Version: "1",
			ReleaseID: "1.0.0-SNAPSHOT",
		},
		uptime.System(),
		uptime.Process(),
		mongodb.Health(url, client, time.Duration(10)*time.Second, time.Duration(40)*time.Microsecond),
		sendgrid.Health(),
		sysinfo.Health(),
	)
	http.HandleFunc("/health", h.Handler)
	http.ListenAndServe(":80", nil)
}
```

## Sample Output (no configured checks)
```json
{
   "releaseId" : "1.0.0-SNAPSHOT",
   "status" : "pass",
   "version" : "1"
}
```

## Sample Output: `mongodb`
```json
{
   "releaseId" : "1.0.0-SNAPSHOT",
   "status" : "pass",
   "version" : "1",
   "checks" : {
      "mongodb:responseTime" : [
         {
            "componentId" : "mongodb://127.0.0.1:27017",
            "observedUnit" : "ns",
            "time" : "2020-03-08T16:48:01.594380018Z",
            "observedValue" : 147640,
            "status" : "pass"
         }
      ]
   }
}
```

## Sample Output: `sendgrid`
```json
{
   "status" : "pass",
   "version" : "1",
   "releaseId" : "1.0.0-SNAPSHOT",
   "checks" : {
      "SendGrid" : [
         {
            "status" : "pass",
            "time" : "2020-03-08T16:45:34.427704957Z"
         }
      ]
   }
}
```

## Sample Output: `uptime`
```json
{
   "status" : "pass",
   "releaseId" : "1.0.0-SNAPSHOT",
   "version" : "1",
   "checks" : {
      "uptime" : [
         {
            "time" : "2020-03-08T16:39:36.409862824Z",
            "observedValue" : 15312,
            "status" : "pass",
            "componentType" : "system",
            "observedUnit" : "s"
         },
         {
            "observedValue" : 6.365804997,
            "time" : "2020-03-08T16:39:36.409871632Z",
            "observedUnit" : "s",
            "componentType" : "process",
            "status" : "pass"
         }
      ]
   }
}
```

## Sample Output: `sysinfo`
```json
{
   "checks" : {
      "memory:utilization" : [
         {
            "componentType" : "system",
            "componentId" : "Total Ram",
            "observedValue" : 16694185984,
            "status" : "pass",
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedUnit" : "1 bytes"
         },
         {
            "componentId" : "Free Ram",
            "componentType" : "system",
            "observedValue" : 672645120,
            "status" : "pass",
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedUnit" : "1 bytes"
         },
         {
            "observedUnit" : "1 bytes",
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedValue" : 190525440,
            "status" : "pass",
            "componentType" : "system",
            "componentId" : "Shared Ram"
         },
         {
            "componentType" : "system",
            "componentId" : "Buffer Ram",
            "observedValue" : 660090880,
            "status" : "pass",
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedUnit" : "1 bytes"
         },
         {
            "componentType" : "system",
            "componentId" : "Total Swap",
            "status" : "pass",
            "observedValue" : 18207465472,
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedUnit" : "1 bytes"
         },
         {
            "observedUnit" : "1 bytes",
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedValue" : 18204581888,
            "status" : "pass",
            "componentId" : "Free Swap",
            "componentType" : "system"
         },
         {
            "componentType" : "system",
            "componentId" : "Total High",
            "status" : "pass",
            "observedValue" : 0,
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedUnit" : "1 bytes"
         },
         {
            "status" : "pass",
            "observedValue" : 0,
            "componentId" : "Free High",
            "componentType" : "system",
            "observedUnit" : "1 bytes",
            "time" : "2020-03-08T16:37:37.559642943Z"
         }
      ],
      "uptime" : [
         {
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedUnit" : "s",
            "componentType" : "system",
            "observedValue" : 15193,
            "status" : "pass"
         }
      ],
      "cpu:utilization" : [
         {
            "componentType" : "system",
            "componentId" : "1 minute",
            "status" : "pass",
            "observedValue" : 0,
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedUnit" : "%"
         },
         {
            "componentId" : "5 minutes",
            "componentType" : "system",
            "observedValue" : 0,
            "status" : "pass",
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedUnit" : "%"
         },
         {
            "componentType" : "system",
            "componentId" : "15 minutes",
            "observedValue" : 0,
            "status" : "pass",
            "time" : "2020-03-08T16:37:37.559642943Z",
            "observedUnit" : "%"
         },
         {
            "status" : "pass",
            "observedValue" : 1449,
            "componentId" : "Processes",
            "componentType" : "system",
            "time" : "2020-03-08T16:37:37.559642943Z"
         }
      ],
      "hostname" : [
         {
            "observedValue" : "Nelkinda-Blade-Stealth-2",
            "status" : "pass",
            "componentId" : "hostname",
            "componentType" : "system",
            "time" : "2020-03-08T16:37:37.559642943Z"
         }
      ]
   },
   "version" : "1",
   "releaseId" : "1.0.0-SNAPSHOT",
   "status" : "pass"
}
```

## References
* Official draft: https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html
* Latest published draft: https://inadarei.github.io/rfc-healthcheck/
* Git Repository of the RFC: https://github.com/inadarei/rfc-healthcheck
