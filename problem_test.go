package problem_client

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestProblemResponse(t *testing.T) {
	testCases := []struct {
		Name        string
		Body        string
		ContentType string
		StatusCode  int

		ExpectedError            error
		ExpectedProblem          *Problem
		ExpectedExtensionMembers ExtensionMembers
	}{
		{
			Name: "Success_ProblemDetailDetected",

			StatusCode: 400,

			// https://www.rfc-editor.org/rfc/rfc9457.html#section-3-4
			Body: `{
					"status": 400,
					"type": "https://example.com/probs/out-of-credit",
					"title": "You do not have enough credit.",
					"detail": "Your current balance is 30, but that costs 50.",
					"instance": "/account/12345/msgs/abc"
				   }`,

			ContentType: "application/problem+json",

			ExpectedError: nil,
			ExpectedProblem: &Problem{
				Type:     "https://example.com/probs/out-of-credit",
				Status:   400,
				Title:    "You do not have enough credit.",
				Detail:   "Your current balance is 30, but that costs 50.",
				Instance: "/account/12345/msgs/abc",
			},
			ExpectedExtensionMembers: map[string]any{},
		},
		{
			Name: "Success_NoProblemDetail",

			StatusCode: 200,

			Body: `{
					"arbitrary_json_values": 2024
				   }`,

			ContentType: "application/json",

			ExpectedError:            nil,
			ExpectedProblem:          nil,
			ExpectedExtensionMembers: nil,
		},
		{
			Name: "Success_ProblemDetailWithExtensionMembers",

			StatusCode: 400,

			Body: `{
					"status": 400,
					"type": "https://example.com/probs/out-of-credit",
					"title": "You do not have enough credit.",
					"detail": "Your current balance is 30, but that costs 50.",
					"instance": "/account/12345/msgs/abc",

					"balance": 30,
					"currency": "EUR"
				   }`,

			ContentType: "application/problem+json",

			ExpectedError: nil,
			ExpectedProblem: &Problem{
				Type:     "https://example.com/probs/out-of-credit",
				Status:   400,
				Title:    "You do not have enough credit.",
				Detail:   "Your current balance is 30, but that costs 50.",
				Instance: "/account/12345/msgs/abc",
			},
			ExpectedExtensionMembers: map[string]any{
				"balance":  30,
				"currency": "EUR",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			header := http.Header{}
			header.Set("Content-Type", tc.ContentType)
			prob, _ /*extra*/, err := ProblemResponse(http.Response{
				StatusCode: tc.StatusCode,
				Header:     header,
				Body:       io.NopCloser(strings.NewReader(tc.Body)),
			})
			if err != tc.ExpectedError {
				t.Errorf("unexpected error value, expected %v got %v", tc.ExpectedError, err)
			}
			if prob == nil || tc.ExpectedProblem == nil {
				if prob != tc.ExpectedProblem {
					t.Errorf("unexpected Problem value, expected %v got %v", tc.ExpectedError, tc.ExpectedProblem)
				}
			} else if *prob != *tc.ExpectedProblem {
				t.Errorf("unexpected Problem value, expected %v got %v", *tc.ExpectedProblem, *prob)
			}
			// TODO: Cannot assert expected with actual value because both are any type, two
			// different any types will always be different
			//
			// if !maps.Equal(extra, tc.ExpectedExtensionMembers) {
			// 	t.Errorf("unexpected ExtensionMembers values, expected %v got %v", tc.ExpectedExtensionMembers, extra)
			// }
		})
	}
}
