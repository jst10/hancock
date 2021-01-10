package api

import (
	"made.by.jst10/outfit7/hancock/cmd/structs"
	"testing"
)

func TestArePerformancesValid(t *testing.T) {
	testCases := []struct {
		Name     string
		Values   []structs.Performance
		Expected bool
	}{
		{"Empty data", []structs.Performance{}, false},
		{"Not all base types", []structs.Performance{
			structs.Performance{AdType: "a", Country: "a", App: "a", Sdk: "a", Score: 1},
		}, false},
		{"All base ad type + extra", []structs.Performance{
			structs.Performance{AdType: "banner", Country: "a", App: "a", Sdk: "a", Score: 1},
			structs.Performance{AdType: "interstitial", Country: "a", App: "a", Sdk: "a", Score: 1},
			structs.Performance{AdType: "rewarded", Country: "a", App: "a", Sdk: "a", Score: 1},
			structs.Performance{AdType: "a", Country: "a", App: "a", Sdk: "a", Score: 1},
		}, false},
		{"All base types", []structs.Performance{
			structs.Performance{AdType: "banner", Country: "a", App: "a", Sdk: "a", Score: 1},
			structs.Performance{AdType: "interstitial", Country: "a", App: "a", Sdk: "a", Score: 1},
			structs.Performance{AdType: "rewarded", Country: "a", App: "a", Sdk: "a", Score: 1},
		}, true},
		{"Not the same countries", []structs.Performance{
			structs.Performance{AdType: "banner", Country: "a", App: "a", Sdk: "a", Score: 1},
			structs.Performance{AdType: "banner", Country: "b", App: "a", Sdk: "a", Score: 1},
			structs.Performance{AdType: "interstitial", Country: "a", App: "a", Sdk: "a", Score: 1},
			structs.Performance{AdType: "rewarded", Country: "a", App: "a", Sdk: "a", Score: 1},
		}, false},
		{"Two sdks per each type", []structs.Performance{
			structs.Performance{AdType: "banner", Country: "c1", App: "a1", Sdk: "s1", Score: 1},
			structs.Performance{AdType: "banner", Country: "c1", App: "a1", Sdk: "s2", Score: 8},
			structs.Performance{AdType: "interstitial", Country: "c1", App: "a1", Sdk: "s1", Score: 10},
			structs.Performance{AdType: "interstitial", Country: "c1", App: "a1", Sdk: "s2", Score: 3},
			structs.Performance{AdType: "rewarded", Country: "c1", App: "a1", Sdk: "s1", Score: 2},
			structs.Performance{AdType: "rewarded", Country: "c1", App: "a1", Sdk: "s2", Score: 5},
		}, true},
		{"Two apps per each type", []structs.Performance{
			structs.Performance{AdType: "banner", Country: "c1", App: "a1", Sdk: "s1", Score: 1},
			structs.Performance{AdType: "banner", Country: "c1", App: "a1", Sdk: "s2", Score: 8},
			structs.Performance{AdType: "banner", Country: "c1", App: "a2", Sdk: "s1", Score: 3},
			structs.Performance{AdType: "interstitial", Country: "c1", App: "a1", Sdk: "s1", Score: 10},
			structs.Performance{AdType: "interstitial", Country: "c1", App: "a1", Sdk: "s2", Score: 3},
			structs.Performance{AdType: "interstitial", Country: "c1", App: "a2", Sdk: "s1", Score: 3},
			structs.Performance{AdType: "rewarded", Country: "c1", App: "a1", Sdk: "s1", Score: 2},
			structs.Performance{AdType: "rewarded", Country: "c1", App: "a1", Sdk: "s2", Score: 5},
			structs.Performance{AdType: "rewarded", Country: "c1", App: "a2", Sdk: "s1", Score: 7},
		}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			valid := arePerformancesValid(tc.Values)
			if valid != tc.Expected {
				t.Errorf("%t != %t", valid, tc.Expected)
			} else {
				t.Logf("%t == %t", valid, tc.Expected)
			}
		})
	}
}
