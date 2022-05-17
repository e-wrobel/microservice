package configuration

import "testing"

func TestParametersValidity(t *testing.T) {

	testCases := []struct {
		name           string
		requiredParams []string
		wantErr        bool
	}{
		{
			name:           "negative required parameter is missing",
			requiredParams: []string{"test1"},
			wantErr:        true,
		},
		{
			name:           "positive no required parameters",
			requiredParams: nil,
			wantErr:        false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateConfig(testCase.requiredParams)
			if (err != nil) != testCase.wantErr {
				t.Errorf("err: %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}
