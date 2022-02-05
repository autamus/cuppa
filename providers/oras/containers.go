package oras

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/cuppa/version"
	"github.com/docker/distribution/reference"
)



func (c Provider) GetImages(url string) (rs *results.ResultSet, err error) {
	result := results.ResultSet{}

	// Normalize the url without the tag.
	urlNormalized := strings.Replace(url, "oras://", "", 1)
	urlData := strings.SplitN(url, ":", 3)
	
	// Regular expressions for versions
	vexp := regexp.MustCompile(`^([0-9]{1,4}[.])+[0-9,a-d]{1,4}`)
	vexpStrict := regexp.MustCompile(`^([0-9]{1,4}[.])+[0-9,a-d]{1,4}$`)
	
	// Prepare filter for tags, if provided
	filter := "*"
	if len(urlData) > 2 {
		filter = urlData[2]
	}
	tags, err := GetImageTags(urlNormalized)
	if err != nil {
		return rs, err
	}

	latest := version.NewVersion("")
	latestTag := ""
	for _, tag := range tags {
		matched, err := filepath.Match(filter, tag)
		if err != nil {
			return rs, err
		}
		if matched {
			verString := vexp.FindString(tag)
			new := version.NewVersion(verString)
			if latest.String() == "N/A" || (verString != "" && new.Compare(latest) < 0) {
				latest = new
				latestTag = tag
			}
		}
	}
	if latest.String() != "N/A" {
		filter = latestTag
	}

	for _, tag := range tags {

		if tag == "" {
			continue
		}
		matched, err := filepath.Match(filter, tag)
		if err != nil {
			return rs, err
		}
		if matched {
			ref, err := reference.ParseNormalizedNamed(urlNormalized + ":" + tag)
			if err != nil {
				return rs, err
			}
			sha, err := GetImageDigest(ref.String())
			if err != nil {
				return rs, err
			}
			output := results.NewResult(
				sha,
				tag,
				"oras://"+ref.String(),

				// Most of these don't have a time, so we use now
				time.Now(),
			)

			// Work around CUPPA's default version handling.
			if !vexpStrict.MatchString(tag) {
				output.Version = []string{tag}
			}
			result.AddResult(output)
		}
	}

	if result.Len() < 1 {
		return nil, results.NotFound
	}
	return &result, nil
}
