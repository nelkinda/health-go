package sendgrid

import (
	"encoding/json"
	"github.com/nelkinda/health-go"
	"net/http"
	"time"
)

type sendGrid struct{}

const sendGridUrl = "http://status.sendgrid.com/"

func getSendGridStatus() health.Details {
	client := &http.Client{Timeout: time.Second * 2}
	req, err := http.NewRequest(http.MethodGet, sendGridUrl, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return health.Details{Status: health.Fail, Output: err.Error()}
	}
	sendGridHealth := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&sendGridHealth)
	if err != nil {
		return health.Details{Status: health.Fail, Output: err.Error()}
	}
	sendGridStatus := sendGridHealth["status"]
	switch vv := sendGridStatus.(type) {
	case map[string]interface{}:
		indicator := vv["indicator"]
		switch indicator {
		case "none":
			return health.Details{Status: health.Pass}
		case "minor", "major":
			return health.Details{Status: health.Warn}
		default:
			description := vv["description"]
			switch descriptionText := description.(type) {
			case string:
				return health.Details{Status: health.Fail, Output: descriptionText}
			default:
				return health.Details{Status: health.Fail, Output: "Could not get description from SendGrid."}
			}
		}
	}
	return health.Details{Status: health.Fail, Output: "Could not parse response from SendGrid."}
}

func (s *sendGrid) HealthDetails() map[string][]health.Details {
	now := time.Now().UTC()
	details := getSendGridStatus()
	details.Time = now.Format(time.RFC3339Nano)
	return map[string][]health.Details{"SendGrid": {details}}
}

func (*sendGrid) AuthorizeHealth(r *http.Request) bool {
	return true
}

func Health() health.DetailsProvider {
	return &sendGrid{}
}
