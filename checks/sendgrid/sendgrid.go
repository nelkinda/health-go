package sendgrid

import (
	"encoding/json"
	"github.com/nelkinda/health-go"
	"net/http"
	"time"
)

type sendGrid struct{}

const sendGridURL = "http://status.sendgrid.com/"

func getSendGridStatus() health.Checks {
	client := &http.Client{Timeout: time.Second * 2}
	req, err := http.NewRequest(http.MethodGet, sendGridURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return health.Checks{Status: health.Fail, Output: err.Error()}
	}
	sendGridHealth := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&sendGridHealth)
	if err != nil {
		return health.Checks{Status: health.Fail, Output: err.Error()}
	}
	sendGridStatus := sendGridHealth["status"]
	switch vv := sendGridStatus.(type) {
	case map[string]interface{}:
		indicator := vv["indicator"]
		switch indicator {
		case "none":
			return health.Checks{Status: health.Pass}
		case "minor", "major":
			return health.Checks{Status: health.Warn}
		default:
			description := vv["description"]
			switch descriptionText := description.(type) {
			case string:
				return health.Checks{Status: health.Fail, Output: descriptionText}
			default:
				return health.Checks{Status: health.Fail, Output: "Could not get description from SendGrid."}
			}
		}
	}
	return health.Checks{Status: health.Fail, Output: "Could not parse response from SendGrid."}
}

func (s *sendGrid) HealthChecks() map[string][]health.Checks {
	now := time.Now().UTC()
	checks := getSendGridStatus()
	checks.Time = now.Format(time.RFC3339Nano)
	return map[string][]health.Checks{"SendGrid": {checks}}
}

func (*sendGrid) AuthorizeHealth(r *http.Request) bool {
	return true
}

// Health returns a ChecksProvider that provides SendGrid health.
// SendGrid health is determined by a simple HTTP ping to SendGrid.
func Health() health.ChecksProvider {
	return &sendGrid{}
}
