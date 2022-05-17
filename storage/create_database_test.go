package storage

import (
	"testing"
)

func TestDatabaseInitialization(t *testing.T) {

	testCases := []struct {
		name    string
		dbFile  string
		wantErr bool
	}{
		{
			name:    "positive has permissions to create database",
			dbFile:  "/var/tmp/database.sql",
			wantErr: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := initDatabase(testCase.dbFile)
			if (err != nil) != testCase.wantErr {
				t.Errorf("err: %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}
