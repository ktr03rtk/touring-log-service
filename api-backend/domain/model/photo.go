package model

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/rwcarlsen/goexif/exif"

	"github.com/pkg/errors"
)

type Photo struct {
	ID          PhotoID   `json:"id"`
	Year        int       `json:"year"`
	Month       int       `json:"month"`
	Day         int       `json:"day"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Timestamp   time.Time `json:"timestamp"`
	S3ObjectKey string    `json:"s3_object_key"`
	UserID      UserID    `json:"user_id"`
}

type (
	PhotoID string
)

const (
	MINIMUM_LAT   = -90
	MAXIMUM_LAT   = 90
	MINIMUM_LON   = -180
	MAXIMUM_LON   = 180
	S3_KEY_PREFIX = "touring-log/photo"
)

// example key: touring-log/photo/thing=thingName/year=2022/month=01/day=12/1656422254000.gz.
var keyValidator = regexp.MustCompile(`^touring-log/photo/thing=.+/year=\d{4}/month=[01][0-9]/day=[0-2][0-9]/\d{13}.jpeg.gz$`)

func NewPhoto(id PhotoID, time time.Time, lat, lon float64, user_id UserID, unit string) (*Photo, error) {
	key := fmt.Sprintf("%s/thing=%s/year=%d/month=%02d/day=%02d/%s.jpeg.gz", S3_KEY_PREFIX, unit, time.Year(), time.Month(), time.Day(), strconv.Itoa(int(time.UnixMilli())))

	data := &Photo{
		ID:          id,
		Year:        time.Year(),
		Month:       int(time.Month()),
		Day:         time.Day(),
		Lat:         lat,
		Lon:         lon,
		Timestamp:   time,
		S3ObjectKey: key,
		UserID:      user_id,
	}

	if err := photoSpecSatisfied(data); err != nil {
		return nil, err
	}

	return data, nil
}

func photoSpecSatisfied(data *Photo) error {
	if !keyValidator.MatchString(data.S3ObjectKey) {
		return errors.Errorf("failed to validate key: %+v", data.S3ObjectKey)
	}

	if data.Lat < MINIMUM_LAT || data.Lat > MAXIMUM_LAT {
		return errors.Errorf("failed to satisfy GPS Lat Spec: %+v", data.Lat)
	}

	if data.Lon < MINIMUM_LON || data.Lon > MAXIMUM_LON {
		return errors.Errorf("failed to satisfy GPS Lon Spec: %+v", data.Lon)
	}

	return nil
}

func ExtractMetadata(r io.Reader) (lat, lon float64, time time.Time, err error) {
	x, err := exif.Decode(r)
	if err != nil {
		return 0, 0, time, errors.Wrapf(err, "failed to decode exif")
	}

	lat, lon, err = x.LatLong()
	if err != nil {
		return 0, 0, time, errors.Wrapf(err, "failed to extract lat-lon")
	}

	time, err = x.DateTime()
	if err != nil {
		return 0, 0, time, errors.Wrapf(err, "failed to extract time")
	}

	return lat, lon, time, nil
}
