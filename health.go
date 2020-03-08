// Package health provides an implementation of the upcoming RFC Health Check Response Format for HTTP APIs.
// To use it, ask health to create a new http.Handler and add it to your server.
// While doing so, you can add information about the service and optionally providers for Checks.
//
// Example:
//  package main
//
//  import (
//  	"github.com/nelkinda/health-go"
//  	"github.com/nelkinda/health-go/checks/uptime"
//  	"net/http"
//  )
//
//  func main() {
//  	h := health.New(
//  		health.Health{
//  			Version: "1",
//  			ReleaseId: "1.0.0-SNAPSHOT",
//  		},
//  		uptime.System(),
//  		uptime.Process(),
//  	)
//  	http.HandleFunc("/health", h.Handler)
//  	http.ListenAndServe(":80", nil)
//  }
//
// References
// * Official draft: https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html
// * Latest published draft: https://inadarei.github.io/rfc-healthcheck/
// * Git Repository of the RFC: https://github.com/inadarei/rfc-healthcheck
package health

import (
	"encoding/json"
	"github.com/nelkinda/http-go/header"
	"github.com/nelkinda/http-go/mimetype"
	"net/http"
)

// Status represents a health status.
// Possible values are pass, warn, and fail.
type Status string

// Health Check Response Format for HTTP APIs uses JSON format described in RFC 8259 and has the media type "application/health+json".
// Its content consists of a single mandatory root field ("status") and several optional fields:
// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-3
type Health struct {

	// status: (required) indicates whether the service status is acceptable or not.
	// API publishers SHOULD use following values for the field:
	// * "pass": healthy (acceptable aliases: "ok" to support Node's Terminus and "up" for Java's SpringBoot),
	// * "fail": unhealthy (acceptable aliases: "error" to support Node's Terminus and "down" for Java's SpringBoot),
	//   and
	// * "warn": healthy, with some concerns.
	//
	// The value of the status field is case-insensitive and tightly related with the HTTP response code returned by the health endpoint.
	// For "pass" and "warn" statuses, HTTP response code in the 2xx-3xx range MUST be used.
	// For "fail" status, HTTP response code in the 4xx-5xx range MUST be used.
	// In case of "warn" status, endpoints MUST return HTTP status in the 2xx-3xx range, and additional information SHOULD be provided, utilizing optional fields of the response.
	//
	// A health endpoint is only meaningful in the context of the component it indicates the health of.
	// It has no other meaning or purpose.
	// As such, its health is a conduit to the health of the component.
	// Clients SHOULD assume that the HTTP response code returned by the health endpoint is applicable to the entire component (e.g. a larger API or a microservice).
	// This is compatible with the behavior that current infrastructural tooling expects: load-balancers, service discoveries and others, utilizing health-checks.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-3.1
	Status Status `json:"status" example:"pass"`

	// version: (optional) public version of the service
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-3.2
	Version string `json:"version,omitempty" example:"1"`

	// releaseId: (optional) in well-designed APIs, backwards-compatible changes in the service should not update a version number.
	// APIs usually change their version number as infrequently as possible, to preserve stable interface.
	// However, implementation of an API may change much more frequently, which leads to the importance of having separate "release number" or "releaseID" that is different from the public version of the API.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-3.3
	// [Note: It is probably recommended to use Semantic Versioning for this field, see https://semver.org/]
	ReleaseID string `json:"releaseId,omitempty" example:"1.14.2-SNAPSHOT"`

	// notes: (optional) array of notes relevant to current state of health
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-3.4
	Notes []string `json:"notes,omitempty"`

	// output: (optional) raw error output, in case of "fail" or "warn" states.
	// This field SHOULD be omitted for "pass" state.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-3.5
	Output string `json:"output,omitempty"`

	// checks (optional) is an object that provides detailed health status of additional downstream systems and endpoints which can affect the overall health of the main API.
	// Please refer to the "The Checks Object" section for more information.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-3.6
	Checks map[string][]Checks `json:"checks,omitempty"`

	// links (optional) is an objects containing link relations and URIs [RFC3986] for external links that MAY contain more information about the health of the endpoint.
	// All values of this object SHALL be URIs.
	// Keys MAY also be URIs.
	// Per web-linking standards [RFC8288] a link relationship SHOULD either be a common/registered one or be indicated as a URI, to avoid name clashes.
	// If a "self" link is provided, it MAY be used by clients to check health via HTTP response code, as mentioned above.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-3.7
	Links map[string]string `json:"links,omitempty"`

	// serviceId (optional) is a unique identifier of the service, in the application scope.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-3.8
	ServiceID string `json:"serviceId,omitempty"`

	// description (optional) is a human-friendly description of the service.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-3.9
	Description string `json:"description,omitempty"`
}

// The Checks Object
// The "checks" object MAY have a number of unique keys, one for each logical downstream dependency or sub-component.
// Since each sub-component may be backed by several nodes with varying health statuses, these keys point to arrays of objects.
// In case of a single-node sub-component (or if presence of nodes is not relevant), a single-element array SHOULD be used as the value, for consistency.
//
// The key identifying an element in the object SHOULD be a unique string within the details section.
// It MAY have two parts: "{componentName}:{measurementName}", in which case the meaning of the parts SHOULD be as follows:
// o componentName: (optional) human-readable name for the component.
//   MUST not contain a colon, in the name, since colon is used as a separator.
// o measurementName: (optional) name of the measurement type (a data point type) that the status is reported for.
//   MUST not contain a colon, in the name, since colon is used as a separator.
//   The observation's name can be one of:
//   * A pre-defined value from this spec.
//     Pre-defined values include:
//     + utilization
//     + responseTime
//     + connections
//     + uptime
//   * A common and standard term from a well-known source such as schema.org, IANA or microformats.
//   * A URI that indicates extra semantics and processing rules that MAY be provided by a resource at the other end of the URI.
//     URIs do not have to be dereferenceable, however.
//     They are just a namespace, and the meaning of a namespace CAN be provided by any convenient means (e.g. publishing an RFC, Swagger document or a nicely printed book).
// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-4
type Checks struct {
	// componentId: (optional) is a unique identifier of an instance of a specific sub-component/dependency of a service.
	// Multiple objects with the same componentID MAY appear in the details, if they are from different nodes.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-4.1
	ComponentID string `json:"componentId,omitempty"`

	// componentType: (optional) SHOULD be present if componentName is present.
	// It's a type of the component and could be one of:
	// * Pre-defined value from this spec. Pre-defined values include:
	//     * component
	//     * datastore
	//     * system
	// * A common and standard term from a well-known source such as schema.org, IANA or microformats.
	// * A URI that indicates extra semantics and processing rules that MAY be provided by a resource at the other end of the URI.
	//   URIs do not have to be dereferenceable, however.
	//   They are just a namespace, and the meaning of a namespace CAN be provided by any convenient means (e.g. publishing an RFC, Swagger document or a nicely printed book).
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-4.2
	ComponentType string `json:"componentType,omitempty"`

	// observedValue: (optional) could be any valid JSON value, such as: string, number, object, array or literal.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-4.3
	ObservedValue interface{} `json:"observedValue,omitempty"`

	// observedUnit (optional) SHOULD be present if observedValue is present.
	// Clarifies the unit of measurement in which observedUnit is reported, e.g. for a time-based value it is important to know whether the time is reported in seconds, minutes, hours or something else.
	// To make sure unit is denoted by a well-understood name or an abbreviation, it should be one of:
	// * A common and standard term from a well-known source such as schema.org, IANA, microformats, or a standards document such as RFC 3339.
	// * A URI that indicates extra semantics and processing rules that MAY be provided by a resource at the other end of the URI.
	//   URIs do not have to be dereferencable, however.
	//   They are just a namespace, and the meaning of a namespace CAN be provided by any convenient means (e.g. publishing an RFC, Open API Spec document or a nicely printed book).
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-4.4
	ObservedUnit string `json:"observedUnit,omitempty"`

	// status (optional) has the exact same meaning as the top-level "output" element, but for the sub-component/downstream dependency represented by the details object.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-4.5
	Status Status `json:"status" example:"pass"`

	// affectedEndpoints (optional) is a JSON array containing URI Templates as defined by [RFC6570].
	// This field SHOULD be omitted if the "status" field is present and has a value equal to "pass".
	// A typical API has many URI endpoints.
	// Most of the time we are interested in the overall health of the API, without diving into details.
	// That said, sometimes operational and resilience middleware needs to know more details about the health of the API (which is why "checks" property provides details).
	// In such cases, we often need to indicate which particular endpoints are affected by a particular check's troubles vs. other endpoints that may be fine.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-4.6
	AffectedEndpoints []string `json:"affectedEndpoints,omitempty"`

	// time (optional) is the date-time, in ISO8601 format, at which the reading of the observedValue was recorded.
	// This assumes that the value can be cached and the reading typically doesn't happen in real time, for performance and scalability purposes.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-4.7
	Time string `json:"time,omitempty" example:"2019-02-20T22:01:44,654015561+00:00"`

	// output (optional) has the exact same meaning as the top-level "output" element, but for the sub-component/downstream dependency represented by the details object.
	// As is the case for the top-level element, this field SHOULD be omitted for "pass" state of a downstream dependency.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-4.8
	Output string `json:"output,omitempty"`

	// links (optional) has the exact same meaning as the top-level "links" element, but for the sub-component/downstream dependency represented by the details object.
	// See https://tools.ietf.org/id/draft-inadarei-api-health-check-04.html#section-4.9
	Links map[string]string `json:"links,omitempty"`
}

const (
	// Pass represents a healthy service "pass"
	Pass Status = "pass"

	// Fail represents an unhealthy service "fail"
	Fail Status = "fail"

	// Warn represents a healthy service with some minor problem "warn"
	Warn Status = "warn"
)

// ChecksProvider provides health checks, potentially with prior authorization.
type ChecksProvider interface {
	// HealthChecks asks the ChecksProvider for its current Health status.
	HealthChecks() map[string][]Checks

	// AuthorizeHealth asks whether the ChecksProvider authorizes Checks to be included in a Health response to this request.
	AuthorizeHealth(r *http.Request) bool
}

// Handler implements the health endpoint.
// @Summary Service health
// @Description Returns the service health according to the upcoming IETF RFC Health Check Response Format for HTTP APIs https://tools.ietf.org/id/draft-inadarei-api-health-check-02.html
// @Produce application/json
// @Success 200 {object} health.Health
// @Router /health [GET]
func (h *Service) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(header.ContentType, mimetype.ApplicationHealthJson)
	if r.Method == http.MethodOptions {
		w.Header().Set("Allow", "OPTIONS, GET, HEAD")
		w.Header().Set("Cache-Control", "max-age=604800")
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.template.Status = Pass
	h.template.Checks = make(map[string][]Checks)
	for _, checksProvider := range h.checksProviders {
		checksMap := checksProvider.HealthChecks()
		for checksKey, checks := range checksMap {
			h.template.Checks[checksKey] = append(h.template.Checks[checksKey], checks...)
		}
	}
	_ = json.NewEncoder(w).Encode(h.template)
}

// Service describes an instance of a health service.
type Service struct {
	// The providers for checks of this health service.
	checksProviders []ChecksProvider
	// The template for the outer health response.
	template         Health
}

// New creates a new health service.
func New(template Health, checksProviders ...ChecksProvider) *Service {
	return &Service{checksProviders: checksProviders, template: template}
}
