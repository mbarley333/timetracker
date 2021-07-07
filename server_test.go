package timetracker_test

import (
	"testing"
	"timetracker"
)

func TestBuildConnectionString(t *testing.T) {
	type testCase struct {
		fn func(string) (string, error)
		a  string
	}
	tcs := []testCase{
		{fn: timetracker.GetEnvironmentVariable, a: "TIMETRACKER_DB_HOST"},
		{fn: timetracker.GetEnvironmentVariable, a: "TIMETRACKER_DB_PORT"},
		{fn: timetracker.GetEnvironmentVariable, a: "TIMETRACKER_DB_USER"},
		{fn: timetracker.GetEnvironmentVariable, a: "TIMETRACKER_DB_NAME"},
	}
	for _, tc := range tcs {
		_, err := tc.fn(tc.a)
		if err != nil {
			t.Fatalf("error with environment variable: %s %s", tc.a, err)
		}
	}

}
