package registryimage

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestFromArtifact_defaults(t *testing.T) {
	artifact := new(packer.MockArtifact)

	img := FromArtifact(artifact)

	if img.ImageID != artifact.Id() {
		t.Errorf("expected resulting registryimage.Image to have a valid ImageID, but it got %q", img.ImageID)
	}

	if img.ProviderName != artifact.BuilderId() {
		t.Errorf("expected resulting registryimage.Image to have a valid ProviderName, but it got %q", img.ProviderName)
	}

	// By default no ProviderRegion is specified when generating from an Artifact
	if img.ProviderRegion != "" {
		t.Errorf("expected resulting registryimage.Image to have no ProviderRegion, but it got %q", img.ProviderRegion)
	}
}

func TestFromArtifact_WithRegion(t *testing.T) {
	region := "TheGreatBeyond"

	artifact := new(packer.MockArtifact)
	img := FromArtifact(artifact, WithRegion(region))

	if img.ProviderRegion != region {
		t.Errorf("expected resulting registryimage.Image to have ProviderRegion of %q, but it got %q", region, img.ProviderRegion)
	}
}

func TestFromArtifact_WithImageID(t *testing.T) {
	id := "some-image-id"

	artifact := new(packer.MockArtifact)
	img := FromArtifact(artifact, WithID(id))

	if img.ImageID != id {
		t.Errorf("expected resulting registryimage.Image to have ImageID of %q, but it got %q", id, img.ImageID)
	}
}

func TestFromArtifact_WithProvider(t *testing.T) {
	provider := "Provider"

	artifact := new(packer.MockArtifact)
	img := FromArtifact(artifact, WithProvider(provider))

	if img.ProviderName != provider {
		t.Errorf("expected resulting registryimage.Image to have ProviderName of %q, but it got %q", provider, img.ProviderName)
	}
}

func TestFromArtifact_SetMetadata(t *testing.T) {

	artifact := new(packer.MockArtifact)
	artifact.StateValues = map[string]interface{}{
		"cloud":            "foo",
		"non-string-value": 7,
		"slice-of-strings": []string{"foo", "bar"},
	}

	img := FromArtifact(artifact, SetMetadata(artifact.StateValues))

	if len(img.Metadata) == 0 {
		t.Errorf("expected resulting registryimage.Image to have some Metadata %q, but it got %q", artifact.StateValues["cloud"], img.Metadata)
	}

	// Only values of string should be stored as metadata
	if len(img.Metadata) > 1 {
		t.Errorf("expected resulting registryimage.Image to have some Metadata only for string values %q, but it got %q", artifact.StateValues["cloud"], img.Metadata)
	}

}
