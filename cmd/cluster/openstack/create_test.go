package openstack

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateOptions_Validate(t *testing.T) {
	for _, test := range []struct {
		name          string
		input         CreateOptions
		expectedError string
	}{
		{
			name:          "missing OpenStack credentials file",
			input:         CreateOptions{},
			expectedError: "OpenStack credentials file is required",
		},
	} {
		var errString string
		if err := test.input.Validate(context.Background(), nil); err != nil {
			errString = err.Error()
		}
		if diff := cmp.Diff(test.expectedError, errString); diff != "" {
			t.Errorf("got incorrect error: %v", diff)
		}
	}
}
