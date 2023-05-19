package ally

import (
	"encoding/xml"
	"net/url"
	"time"
)

// GET utility/status
func APIStatus() IAPIStatus {
	return get[IAPIStatus]("utility/status", url.Values{}, infiniteRL)
}

// GET utility/version
func APIVersion() IAPIVersion {
	return get[IAPIVersion]("utility/version", url.Values{}, infiniteRL)
}

func (t *serverTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Note: The Server API returns the wrong time; it returns EST stamped as GMT
	const format = "Mon, 02 Jan 2006 15:04:05 GMT"

	var v string
	d.DecodeElement(&v, &start)
	est, _ := time.LoadLocation("America/New_York")
	parse, err := time.ParseInLocation(format, v, est)
	if err != nil {
		return err
	}
	t.Time = parse
	return nil
}

type serverTime struct {
	time.Time
}
