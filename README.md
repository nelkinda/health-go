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

You can add any implementation of `DetailsProvider` to the varargs list of `health.New()`.

```go
package main

import (
	"github.com/nelkinda/health-go"
	"github.com/nelkinda/health-go/details/uptime"
	"net/http"
)

func main() {
	h := health.New(health.Health{Version: "1", ReleaseId: "1.0.0-SNAPSHOT"}, uptime.System(), uptime.Process())
	http.HandleFunc("/health", h.Handler)
	http.ListenAndServe(":80", nil)
}
```

## References
* Official draft: https://tools.ietf.org/id/draft-inadarei-api-health-check-02.html
* Latest published draft: https://inadarei.github.io/rfc-healthcheck/
* Git Repository of the RFC: https://github.com/inadarei/rfc-healthcheck
