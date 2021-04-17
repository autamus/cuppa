package docker

import (
	"regexp"

	"github.com/DataDrake/cuppa/results"
)

var (
	// SourceRegex is the regex for Docker sources
	SourceRegex = regexp.MustCompile("docker://([^/]*/[^/.]*)*")
)

// Provider is the upstream provider interface for docker.
type Provider struct{}

// String gives the name of this provider
func (c Provider) String() string {
	return "Docker"
}

// Match check to see of this provider can handle this kind of query
func (c Provider) Match(query string) (params []string) {
	if sm := SourceRegex.FindStringSubmatch(query); len(sm) > 1 {
		params = []string{query}
	}
	return
}

// Latest finds the newest release for a docker image
func (c Provider) Latest(params []string) (r *results.Result, err error) {
	rs, err := c.GetImages(params[0])
	if err == nil {
		r = rs.Last()
	}
	return
}

// Releases finds all matching releases for a docker image
func (c Provider) Releases(params []string) (rs *results.ResultSet, err error) {
	rs, err = c.GetImages(params[0])
	return
}
