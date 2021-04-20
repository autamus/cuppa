package docker

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DataDrake/cuppa/results"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/transports"
	"github.com/containers/image/v5/types"
	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
)

func (c Provider) GetImages(url string) (rs *results.ResultSet, err error) {
	result := results.ResultSet{}
	ctx := context.Background()
	sys := &types.SystemContext{}
	// Normalize the url without the tag.
	urlData := strings.SplitN(url, ":", 3)
	urlNormalized := strings.Join(urlData[:2], ":")
	// Grab the tag pattern from the end of the string if possible.
	vexp := regexp.MustCompile(`^([0-9]{1,4}[.])+[0-9,a-d]{1,4}$`)
	filter := "*"
	if len(urlData) > 2 {
		filter = urlData[2]
		if vexp.MatchString(filter) {
			filter = "*.*"
		}
	}

	transport := transportFromImageName(urlNormalized)
	if transport.Name() != docker.Transport.Name() {
		return rs, fmt.Errorf("Unsupported transport '%v' for tag listing. Only '%v' currently supported", transport.Name(), docker.Transport.Name())
	}

	// Do transport-specific parsing and validation to get an image reference
	imgRef, err := parseDockerRepositoryReference(urlNormalized)
	if err != nil {
		return rs, err
	}

	tags, err := listDockerTags(ctx, sys, imgRef)
	if err != nil {
		return rs, err
	}

	for _, tag := range tags {
		matched, err := filepath.Match(filter, tag)
		if err != nil {
			return rs, err
		}
		if matched {
			ref, err := reference.ParseNormalizedNamed(imgRef.DockerReference().Name() + ":" + tag)
			if err != nil {
				return rs, err
			}
			tagRef, err := docker.NewReference(ref)
			if err != nil {
				return rs, err
			}
			img, err := tagRef.NewImage(ctx, sys)
			if err != nil {
				return rs, err
			}
			insp, err := img.Inspect(ctx)
			if err != nil {
				return rs, err
			}
			sha, err := docker.GetDigest(ctx, sys, tagRef)
			if err != nil {
				return rs, err
			}
			output := results.NewResult(
				sha.String(),
				tag,
				"docker://"+ref.String(),
				*insp.Created,
			)
			// Work around CUPPA's default version handling.
			if !vexp.MatchString(tag) {
				output.Version = []string{tag}
			}
			result.AddResult(output)

		}
	}

	return &result, nil
}

// TransportFromImageName converts an URL-like name to a types.ImageTransport or nil when
// the transport is unknown or when the input is invalid.
func transportFromImageName(imageName string) types.ImageTransport {
	// Keep this in sync with ParseImageName!
	parts := strings.SplitN(imageName, ":", 2)
	if len(parts) == 2 {
		return transports.Get(parts[0])
	}
	return nil
}

// Customized version of the alltransports.ParseImageName and docker.ParseReference that does not place a default tag in the reference
// Would really love to not have this, but needed to enforce tag-less and digest-less names
func parseDockerRepositoryReference(refString string) (types.ImageReference, error) {
	if !strings.HasPrefix(refString, docker.Transport.Name()+"://") {
		return nil, errors.Errorf("docker: image reference %s does not start with %s://", refString, docker.Transport.Name())
	}

	parts := strings.SplitN(refString, ":", 2)
	if len(parts) != 2 {
		return nil, errors.Errorf(`Invalid image name "%s", expected colon-separated transport:reference`, refString)
	}

	ref, err := reference.ParseNormalizedNamed(strings.TrimPrefix(parts[1], "//"))
	if err != nil {
		return nil, err
	}

	if !reference.IsNameOnly(ref) {
		return nil, errors.New(`No tag or digest allowed in reference`)
	}

	// Checks ok, now return a reference. This is a hack because the tag listing code expects a full image reference even though the tag is ignored
	return docker.NewReference(reference.TagNameOnly(ref))
}

// List the tags from a repository contained in the imgRef reference. Any tag value in the reference is ignored
func listDockerTags(ctx context.Context, sys *types.SystemContext, imgRef types.ImageReference) ([]string, error) {
	tags, err := docker.GetRepositoryTags(ctx, sys, imgRef)
	if err != nil {
		return tags, fmt.Errorf("Error listing repository tags: %v", err)
	}

	return tags, nil
}
