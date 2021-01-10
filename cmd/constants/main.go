package constants

import (
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom_errors"
)

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

func AdTypeNameToId(addTypeName string) (int,   *custom_errors.CustomError) {
	switch addTypeName {
	case AdTypeBanner:
		return AdTypeBannerId, nil
	case AdTypeInterstitial:
		return AdTypeInterstitialId, nil
	case AdTypeRewarded:
		return AdTypeRewardedId, nil
	default:
		return -1, custom_errors.GetNotValidDataError("invalid ad type")
	}
}

func AdTypeIdToName(addTypeId int) (string,   *custom_errors.CustomError) {
	switch addTypeId {
	case AdTypeBannerId:
		return AdTypeBanner, nil
	case AdTypeInterstitialId:
		return AdTypeInterstitial, nil
	case AdTypeRewardedId:
		return AdTypeRewarded, nil
	default:
		return "", custom_errors.GetNotValidDataError("invalid ad type")
	}
}
