package main

import (
	"testing"
	"github.com/graphql-go/graphql"
	"encoding/json"
)

func TestGql(t *testing.T) {
	cases := []struct {
		query    string
		expected string
	}{
		{"{hello}", `{"hello":"world"}`},
		{
			query:    `mutation { login(user: "james", password: "abc123") { token } }`,
			expected: `{"login":{"token":"fake-token"}}`,
		},
		{
			query:    `{ hello bar obj { name subobj { foo } } }`,
			expected: `{"bar":"baz","hello":"world","obj":{"name":"James","subobj":{"foo":"bar"}}}`,
		},
	}

	for i, c := range cases {
		params := graphql.Params{Schema: schema, RequestString: c.query}
		r := graphql.Do(params)
		if r.HasErrors() {
			for _, e := range r.Errors {
				t.Error(i, "Error: ", e.Message)
			}
			continue
		}
		rJSON, err := json.Marshal(r.Data)
		if err != nil {
			t.Error("request failed", err)
		}
		if string(rJSON) != c.expected {
			t.Errorf("expected '%s', got '%s'\n", c.expected, string(rJSON))
		}
	}
}
