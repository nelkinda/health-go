package health

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/christianhujer/assert"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var variables map[string]string
var handler http.HandlerFunc

type sampleCheck struct{}

func (s *sampleCheck) HealthChecks(ctx context.Context) map[string][]Checks {
	return map[string][]Checks{
		"sampleCheck": {
			{
				ComponentType: "sampleCheck",
				Status:        Pass,
			},
		},
	}
}

func (*sampleCheck) AuthorizeHealth(r *http.Request) bool {
	return true
}

func SampleCheck() ChecksProvider {
	return &sampleCheck{}
}

func initHandler() {
	if handler != nil {
		return
	}
	h := New(
		Health{
			Version:   "1",
			ReleaseID: "1.0.0-SNAPSHOT",
		},
		WithChecksProviders(SampleCheck()),
		WithTracer(&opentracing.NoopTracer{}, opentracing.HTTPHeaders),
	)
	handler = h.Handler
}

func assertHealthResponseCode(t *testing.T, method string, expectedStatusCode int) {
	r, err := http.NewRequest(method, "/health", nil)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	handler(w, r)

	_ = assert.Equals(t, expectedStatusCode, w.Code)
}

func TestHandlerResponseCodes(t *testing.T) {
	initHandler()
	assertHealthResponseCode(t, http.MethodConnect, http.StatusMethodNotAllowed)
	assertHealthResponseCode(t, http.MethodDelete, http.StatusMethodNotAllowed)
	assertHealthResponseCode(t, http.MethodGet, http.StatusOK)
	assertHealthResponseCode(t, http.MethodHead, http.StatusOK)
	assertHealthResponseCode(t, http.MethodOptions, http.StatusOK)
	assertHealthResponseCode(t, http.MethodPatch, http.StatusMethodNotAllowed)
	assertHealthResponseCode(t, http.MethodPost, http.StatusMethodNotAllowed)
	assertHealthResponseCode(t, http.MethodPut, http.StatusMethodNotAllowed)
	assertHealthResponseCode(t, http.MethodTrace, http.StatusMethodNotAllowed)
}

func TestHandlerResponse(t *testing.T) {
	initHandler()
	ResetVariables(nil)
	r, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	handler(w, r)

	err = AssertJSONBytes(t, []byte(`
        {
            "status": "pass",
            "version": "(?P<version>\\d+)",
            "releaseId": "(?P<releaseId>\\d+\\.\\d+\\.\\d+(?:-\\w+)?)",
            "checks": {
                "sampleCheck": [
                    {
                        "componentType": "sampleCheck",
                        "Status": "PASS"
                    }
                ]
            }
        }
    `), w.Body.Bytes())
	if err != nil {
		t.Error(err)
	}
	_ = assert.Nil(t, err)
}

// AssertJSONBytes asserts that two JSON structures given as binary data are equal.
// Leaf values are compared using EqualsWithCaptureAndReplace.
func AssertJSONBytes(t *testing.T, expected []byte, actual []byte) error {
	var expectedMap map[string]interface{}
	var actualMap map[string]interface{}
	if err := json.Unmarshal(expected, &expectedMap); err != nil {
		return err
	}
	if err := json.Unmarshal(actual, &actualMap); err != nil {
		return err
	}
	return AssertJSON(t, expectedMap, actualMap)
}

// AssertJSON asserts that two JSON structures given as maps are equal.
// Leaf values are compared using EqualsWithCaptureAndReplace.
func AssertJSON(t *testing.T, expected map[string]interface{}, actual map[string]interface{}) error {
	for key, value := range expected {
		switch actual[key].(type) {
		case map[string]interface{}:
			err := AssertJSON(t, expected[key].(map[string]interface{}), actual[key].(map[string]interface{}))
			if err != nil {
				return err
			}
		case string:
			if EqualsWithCaptureAndReplace(t, actual[key].(string), value.(string)) != nil {
				return fmt.Errorf("expected JSON key/value pair: \"%s\": \"%s\", actual JSON key/value pair: \"%s\": \"%s\"", key, expected[key].(string), key, actual[key].(string))
			}
		case nil:
			return fmt.Errorf("expected JSON key/value pair: \"%s\": \"%s\", actual JSON key/value pair: missing", key, expected[key].(string))
		}
	}
	return nil
}

// EqualsWithCaptureAndReplace compares an input string against a pattern.
// It supports capturing using the following syntaxes:
// * generic regular expression capture `(?<name>regex)`
// * Golang regular expression capture `(?P<name>regex)`
// * placeholder capture `(>name)` (uses `.*` as regex)
func EqualsWithCaptureAndReplace(t *testing.T, input string, pattern string) error {
	pattern = regexp.MustCompile(`\(\?<`).ReplaceAllString(pattern, `(?P<`)
	pattern = regexp.MustCompile(`\$\(>([^()]+)\)`).ReplaceAllString(pattern, `(?P<$1>.*?)`)
	regex := regexp.MustCompile(`^` + pattern + `$`)
	match := regex.FindStringSubmatch(input)
	if match != nil {
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				variables[name] = match[i]
			}
		}
	}
	pattern = replace(pattern)
	return assert.True(t, regexp.MustCompile(pattern).MatchString(input))
}

func replace(pattern string) string {
	for name, value := range variables {
		pattern = regexp.MustCompile(`\$\([<]`+name+`\)`).ReplaceAllString(pattern, value)
	}
	return pattern
}

// ResetVariables clears the internal variables map.
// Call this before every test case.
func ResetVariables(_ interface{}) {
	variables = make(map[string]string)
}
