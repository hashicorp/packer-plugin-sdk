package registryimage

import (
	"errors"

	"github.com/hashicorp/packer-plugin-sdk/packer"
)

// ArtifactStateURI represents the key used by Packer when querying an packersdk.Artifact
// for Image metadata that a particular component would like to have stored on the HCP Packer Registry.
const ArtifactStateURI = "par.artifact.metadata"

// ArtifactOverrideFunc represents a transformation func that can be applied to a
// non-nil *Image. See With* functions for examples.
type ArtifactOverrideFunc func(*Image) error

// Image represents the metadata for some Artifact in the HCP Packer Registry.
type Image struct {
	ImageID                      string
	ProviderName, ProviderRegion string
	Metadata                     map[string]string
}

// FromArtifact returns a *Image that can be used by Packer core for publishing to the HCP Packer Registry.
// By default FromArtifact will use the a.BuilderID as the Image Provider, and the a.Id() as the ImageID that
// should be tracked within the HCP Packer Registry. No Region is selected by default as region varies per build.
// The use of one or more ArtifactOverrideFunc can be used to override any of the defaults used.
func FromArtifact(a packer.Artifact, opts ...ArtifactOverrideFunc) *Image {
	if a == nil {
		return nil
	}

	img := &Image{
		ProviderName: a.BuilderId(),
		ImageID:      a.Id(),
		Metadata:     make(map[string]string),
	}

	// Let's grab some state data
	for _, opt := range opts {
		err := opt(img)
		if err != nil {
			return nil
		}
	}

	return img
}

// WithProvider takes a name, and returns a ArtifactOverrideFunc that can be
// used to override the set name for an existing Image.
func WithProvider(name string) func(*Image) error {
	return func(img *Image) error {
		if img == nil {
			return errors.New("no go on empty image")
		}

		img.ProviderName = name
		return nil
	}
}

// WithID takes a id, and returns a ArtifactOverrideFunc that can be
// used to override the set id for an existing Image.
func WithID(id string) func(*Image) error {
	return func(img *Image) error {
		if img == nil {
			return errors.New("no go on empty image")
		}

		img.ImageID = id
		return nil
	}
}

// WithRegion takes a region, and returns a ArtifactOverrideFunc that can be
// used to override the set region for an existing Image.
func WithRegion(region string) func(*Image) error {
	return func(img *Image) error {
		if img == nil {
			return errors.New("no go on empty image")
		}

		img.ProviderRegion = region
		return nil
	}
}

// SetMetadata takes metadata, and returns a ArtifactOverrideFunc that can be
// used to set metadata for an existing Image. The incoming metadata `md`
// will be filtered only for keys whose values are of type string.
// If you wish to override this behavior you may create your own  ArtifactOverrideFunc
// for manipulating and setting Image metadata.
func SetMetadata(md map[string]interface{}) func(*Image) error {
	return func(img *Image) error {
		for k, v := range md {
			v, ok := v.(string)
			if !ok {
				continue
			}
			img.Metadata[k] = v
		}

		return nil
	}
}
