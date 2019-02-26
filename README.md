# health-go

Golang implementation of the upcoming IETF RFC Health Check Response Format for HTTP APIs.

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
	h := health.New(health.Health{Version: "1", ReleaseId: "1.0.0-SNAPSHOT"}) 

	// 2. Add the handler to your mux/server.
	http.HandleFunc("/health", h.Handler)
	
	// 3. Start your server.
	http.ListenAndServe(":80", nil)
}
```

## Providing Details
If is possible to provide details.
This library comes with the following details predefined:
- system uptime
- process uptime
- mongodb health
- SendGrid health

You can add any implementation of `DetailsProvider` to the varargs list of `health.New()`.

```go
package main

import (
	"context"
	"github.com/nelkinda/health-go"
	"github.com/nelkinda/health-go/details/uptime"
	"github.com/nelkinda/health-go/details/mongodb"
	"github.com/nelkinda/health-go/details/sendgrid"
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
			ReleaseId: "1.0.0-SNAPSHOT",
		},
		uptime.System(),
		uptime.Process(),
		mongodb.Health(url, client, time.Duration(10)*time.Second, time.Duration(40)*time.Microsecond),
		sendgrid.Health(),
	)
	http.HandleFunc("/health", h.Handler)
	http.ListenAndServe(":80", nil)
}
```

## References
* Official draft: https://tools.ietf.org/id/draft-inadarei-api-health-check-02.html
* Latest published draft: https://inadarei.github.io/rfc-healthcheck/
* Git Repository of the RFC: https://github.com/inadarei/rfc-healthcheck
