// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestFromMappedData_badInput(t *testing.T) {
	_, err := FromMappedData([]string{"invalid"}, func(k, v interface{}) (*Image, error) { return nil, nil })
	if err == nil {
		t.Errorf("unexpected results for bad map input; expected err to be non-nil")
	}

	if err.Error() != "error the incoming mappedData does not appear to be a map; found type to be"+reflect.Slice.String() {
		t.Errorf("unexpected error returned for bad map input: got %s", err)
	}

}

func TestFromMappedData(t *testing.T) {

	cloudimages := map[string]string{
		"west": "happycloud-1",
		"east": "happycloud-2",
	}

	f := func(key, value interface{}) (*Image, error) {
		v, ok := value.(string)
		if !ok {
			return nil, errors.New("for happycloud maps value should always be string")
		}
		k, ok := key.(string)
		if !ok {
			return nil, errors.New("for happycloud maps key should always be string")
		}

		img := Image{ProviderName: "happycloud", ProviderRegion: k, ImageID: v}
		return &img, nil
	}

	images, err := FromMappedData(cloudimages, f)
	if err != nil {
		t.Fatalf("unexpected error when creating an images from mapped data: %v", err)
	}

	if len(images) != 2 {
		t.Errorf("expected resulting image count to be 2, but got %d", len(images))
	}
}

func TestFromArtifact_defaults(t *testing.T) {
	artifact := new(packer.MockArtifact)

	img, err := FromArtifact(artifact)
	if err != nil {
		t.Fatal("unexpected error when creating an image from a MockArtifact")
	}

	if img.ImageID != artifact.Id() {
		t.Errorf("expected resulting Image to have a valid ImageID, but it got %q", img.ImageID)
	}

	if img.ProviderName != artifact.BuilderId() {
		t.Errorf("expected resulting Image to have a valid ProviderName, but it got %q", img.ProviderName)
	}

	// By default no ProviderRegion is specified when generating from an Artifact
	if img.ProviderRegion != "" {
		t.Errorf("expected resulting Image to have no ProviderRegion, but it got %q", img.ProviderRegion)
	}
}

func TestFromArtifact_WithRegion(t *testing.T) {
	region := "TheGreatBeyond"

	artifact := new(packer.MockArtifact)
	img, err := FromArtifact(artifact, WithRegion(region))
	if err != nil {
		t.Fatal("unexpected error when creating an image from a MockArtifact with a region override.")
	}

	if img.ProviderRegion != region {
		t.Errorf("expected resulting Image to have ProviderRegion of %q, but it got %q", region, img.ProviderRegion)
	}
}

func TestFromArtifact_WithImageID(t *testing.T) {
	id := "some-image-id"

	artifact := new(packer.MockArtifact)
	img, err := FromArtifact(artifact, WithID(id))
	if err != nil {
		t.Fatal("unexpected error when creating an image from a MockArtifact with an ID override.")
	}

	if img.ImageID != id {
		t.Errorf("expected resulting Image to have ImageID of %q, but it got %q", id, img.ImageID)
	}
}

func TestFromArtifact_WithProvider(t *testing.T) {
	provider := "Provider"

	artifact := new(packer.MockArtifact)
	img, err := FromArtifact(artifact, WithProvider(provider))
	if err != nil {
		t.Fatal("unexpected error when creating an image from a MockArtifact with an Provider override.")
	}

	if img.ProviderName != provider {
		t.Errorf("expected resulting Image to have ProviderName of %q, but it got %q", provider, img.ProviderName)
	}
}

func TestFromArtifact_SetLabels(t *testing.T) {

	artifact := new(packer.MockArtifact)
	artifact.StateValues = map[string]interface{}{
		"cloud":            "foo",
		"non-string-value": 7,
		"slice-of-strings": []string{"foo", "bar"},
	}

	img, err := FromArtifact(artifact, SetLabels(artifact.StateValues))
	if err != nil {
		t.Fatal("unexpected error when creating an image from a MockArtifact with some base Metadata.")
	}

	if len(img.Labels) == 0 {
		t.Errorf("expected resulting Image to have some Metadata %q, but it got %q", artifact.StateValues["cloud"], img.Labels)
	}

	// Only values of string should be stored as metadata
	if len(img.Labels) > 1 {
		t.Errorf("expected resulting Image to have some Metadata only for string values %q, but it got %q", artifact.StateValues["cloud"], img.Labels)
	}

}
