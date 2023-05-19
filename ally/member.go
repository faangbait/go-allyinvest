package ally

import (
	"net/url"
)

// GET member/profile
func Profile() IProfile {
	return get[IProfile]("member/profile", url.Values{}, authedRL)
}
