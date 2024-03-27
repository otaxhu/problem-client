package problem_client

import (
	"encoding/json"
	"net/http"
)

type Problem struct {
	Type     string `json:"type,omitempty"`
	Status   int    `json:"status,omitempty"`
	Title    string `json:"title,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

type ExtensionMembers map[string]any

func ProblemResponse(res http.Response) (*Problem, ExtensionMembers, error) {
	if res.Header.Get("Content-Type") != "application/problem+json" {
		return nil, nil, nil
	}

	var x any
	if err := json.NewDecoder(res.Body).Decode(&x); err != nil {
		return nil, nil, err
	}

	p := &Problem{
		// "The "status" member, if present, is only advisory; it conveys the HTTP status code
		// used for the convenience of the consumer."
		//
		// https://www.rfc-editor.org/rfc/rfc9457.html#section-3.1.2-2
		//
		// Taking the status code from the http response
		Status: res.StatusCode,
	}

	extensionMembers := map[string]any{}

	m, ok := x.(map[string]any)
	if !ok {
		// "If a member's value type does not match the specified type, the member MUST be ignored"
		//
		// https://www.rfc-editor.org/rfc/rfc9457.html#section-3.1-1
		//
		// Here the entire JSON object doesn't match the specified schema, so we return a zero-
		// valued Problem struct, with Status set to status code from response and empty extension
		// members
		return p, extensionMembers, nil
	}

	for k, v := range m {

		switch k {
		case "status":
			/* No-op */

		case "type":
			val, ok := v.(string)
			if !ok {
				continue
			}
			p.Type = val

		case "title":
			val, ok := v.(string)
			if !ok {
				continue
			}
			p.Title = val

		case "detail":
			val, ok := v.(string)
			if !ok {
				continue
			}
			p.Detail = val

		case "instance":
			val, ok := v.(string)
			if !ok {
				continue
			}
			p.Instance = val

		default:
			extensionMembers[k] = v
		}
	}

	return p, extensionMembers, nil
}
