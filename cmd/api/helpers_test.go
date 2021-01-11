package api

import (
	"made.by.jst10/outfit7/hancock/cmd/structs"
	"testing"
)

func TestRemoveSdksFromPerformanceResults(t *testing.T) {
	testCases := []struct {
		Name     string
		Arg1     map[string]bool
		Arg2     []structs.SdkScore
		Expected []structs.SdkScore
	}{
		{"Empty data", map[string]bool{}, []structs.SdkScore{}, []structs.SdkScore{}},
		{"Nothing should removed data", map[string]bool{"r1": true},
			[]structs.SdkScore{{Sdk: "s1", Score: 1}},
			[]structs.SdkScore{{Sdk: "s1", Score: 1}}},
		{"Remove first", map[string]bool{"r1": true, "r2": true, "s1": true},
			[]structs.SdkScore{{Sdk: "s1", Score: 1}, {Sdk: "s2", Score: 2}, {Sdk: "s3", Score: 3}},
			[]structs.SdkScore{{Sdk: "s2", Score: 2}, {Sdk: "s3", Score: 3}}},
		{"Remove middle", map[string]bool{"r1": true, "r2": true, "s2": true},
			[]structs.SdkScore{{Sdk: "s1", Score: 1}, {Sdk: "s2", Score: 2}, {Sdk: "s3", Score: 3}},
			[]structs.SdkScore{{Sdk: "s1", Score: 1}, {Sdk: "s3", Score: 3}}},
		{"Remove last two", map[string]bool{"r1": true, "s2": true, "s3": true},
			[]structs.SdkScore{{Sdk: "s1", Score: 1}, {Sdk: "s2", Score: 2}, {Sdk: "s3", Score: 3}},
			[]structs.SdkScore{{Sdk: "s1", Score: 1}}},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			results := removeSdksFromPerformanceResults(tc.Arg1, tc.Arg2)
			expectedLength := len(tc.Expected)
			receivedLength := len(results)
			if expectedLength != receivedLength {
				t.Errorf("Array lengths: %d != %d", receivedLength, expectedLength)
			} else {
				for index, _ := range results {
					if results[index].Sdk != tc.Expected[index].Sdk || results[index].Score != tc.Expected[index].Score {
						t.Errorf("Items at index %d are not the sames: %s != %s", index, results[index].Sdk, tc.Expected[index].Sdk)
					}
					return
				}
			}
		})
	}
}
