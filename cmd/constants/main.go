package constants

import "errors"

const UserRoleGuest = 10
const UserRoleAdmin = 20
const AdTypeBanner = "banner"
const AdTypeInterstitial = "interstitial"
const AdTypeRewarded = "rewarded"
const AdTypeBannerId = 0
const AdTypeInterstitialId = 1
const AdTypeRewardedId = 2

var AddTypes = []string{AdTypeBanner, AdTypeInterstitial, AdTypeRewarded}
var AddTypesIds = []int{AdTypeBannerId, AdTypeInterstitialId, AdTypeRewardedId}

func adTypeNameToId(addTypeName string) (int, error) {
	switch addTypeName {
	case AdTypeBanner:
		return AdTypeBannerId, nil
	case AdTypeInterstitial:
		return AdTypeInterstitialId, nil
	case AdTypeRewarded:
		return AdTypeRewardedId, nil
	default:
		return -1, errors.New("invalid ad type")
	}
}

func adTypeIdToName(addTypeId int) (string, error) {
	switch addTypeId {
	case AdTypeBannerId:
		return AdTypeBanner, nil
	case AdTypeInterstitialId:
		return AdTypeInterstitial, nil
	case AdTypeRewardedId:
		return AdTypeRewarded, nil
	default:
		return "", errors.New("invalid ad type")
	}
}
