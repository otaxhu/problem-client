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

// ProblemResponse reads the header Content-Type, and if it is "application/problem+json", reads
// the body from the response, parses the JSON as specified in RFC 9457 "Problem Details
// Specification", closes the body and returns the structure representing the problem details
//
// Also if there was some members that are not part of the specification (known as "Extension
// Members"), those are gonna be returned in the ExtensionMembers map
//
// Notes:
//
// This function doesn't reads the body if the Content-Type header doesn't match the content type
// specified before, nor closes the body
//
// On the other hand, if that condition matches, this function reads the body and closes the
// body, even if there was an error parsing the JSON
func ProblemResponse(res *http.Response) (problem *Problem, extensionMembers ExtensionMembers, err error) {
	if res.Header.Get("Content-Type") != "application/problem+json" {
		return
	}

	defer res.Body.Close()

	var x any
	if err = json.NewDecoder(res.Body).Decode(&x); err != nil {
		return
	}

	problem = &Problem{
		// "The "status" member, if present, is only advisory; it conveys the HTTP status code
		// used for the convenience of the consumer."
		//
		// https://www.rfc-editor.org/rfc/rfc9457.html#section-3.1.2-2
		//
		// Taking the status code from the http response
		Status: res.StatusCode,
	}

	extensionMembers = map[string]any{}

	m, ok := x.(map[string]any)
	if !ok {
		// "If a member's value type does not match the specified type, the member MUST be ignored"
		//
		// https://www.rfc-editor.org/rfc/rfc9457.html#section-3.1-1
		//
		// Here the entire JSON object doesn't match the specified schema, so we return a zero-
		// valued Problem struct, with Status set to status code from response and empty extension
		// members
		return
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
			problem.Type = val

		case "title":
			val, ok := v.(string)
			if !ok {
				continue
			}
			problem.Title = val

		case "detail":
			val, ok := v.(string)
			if !ok {
				continue
			}
			problem.Detail = val

		case "instance":
			val, ok := v.(string)
			if !ok {
				continue
			}
			problem.Instance = val

		default:
			extensionMembers[k] = v
		}
	}

	return
}
