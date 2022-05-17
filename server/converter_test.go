package server

import (
	"testing"

	"github.com/e-wrobel/90poe/storage"

	"github.com/magiconair/properties/assert"
)

func TestConverterValidity(t *testing.T) {

	testCases := []struct {
		name    string
		obj     *storage.PortEntity
		wantErr bool
	}{
		{
			name: "positive non empty input",
			obj: &storage.PortEntity{
				Identifier: "FGS6Z",
				PortDetails: &storage.PortDetails{
					Name:    "TestName",
					City:    "TestCity",
					Country: "TestCountry",
				},
			},
			wantErr: false,
		},
		{
			name:    "negative empty identifier",
			obj:     &storage.PortEntity{},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			identifier, _, err := portEntityToBytes(testCase.obj)
			if (err != nil) != testCase.wantErr {
				t.Errorf("err: %v, wantErr %v", err, testCase.wantErr)
			}
			assert.Equal(t, identifier, testCase.obj.Identifier)
		})
	}
}
