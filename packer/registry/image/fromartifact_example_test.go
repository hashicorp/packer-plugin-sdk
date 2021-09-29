package image_test

import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
)

type simpleArtifact struct {
	image_id string
}

func (a *simpleArtifact) BuilderId() string {
	return "example.happycloud"
}

func (a *simpleArtifact) Files() []string {
	return nil
}

func (a *simpleArtifact) Id() string {
	return a.image_id
}

func (a *simpleArtifact) String() string {
	return fmt.Sprintf("Imported image URL: %s", a.Id())
}

func (a *simpleArtifact) State(name string) interface{} {
	return nil
}

func (a *simpleArtifact) Destroy() error {
	return nil
}

func ExampleFromArtifact() {

	a := &simpleArtifact{
		image_id: "service-id-123",
	}

	hcimage, _ := image.FromArtifact(a)
	fmt.Printf("%#v", *hcimage)
	// Output:
	// image.Image{ImageID:"service-id-123", ProviderName:"example.happycloud", ProviderRegion:"", Labels:map[string]string{}, SourceImageID:""}
}

func ExampleWithProvider() {
	a := &simpleArtifact{
		image_id: "service-id-123",
	}

	// This example also includes an override for the ProviderRegion to illustrate how ArtifactOverrideFunc(s) can be chained.
	hcimage, _ := image.FromArtifact(a, image.WithProvider("happycloud"), image.WithRegion("west"))
	fmt.Printf("%#v", *hcimage)
	// Output:
	// image.Image{ImageID:"service-id-123", ProviderName:"happycloud", ProviderRegion:"west", Labels:map[string]string{}, SourceImageID:""}
}

func ExampleWithSourceImageImageID() {
	a := &simpleArtifact{
		image_id: "service-id-123",
	}

	// This example also includes an override for the ProviderRegion to illustrate how ArtifactOverrideFunc(s) can be chained.
	hcimage, _ := image.FromArtifact(a, image.WithProvider("happycloud"), image.WithRegion("west"), image.WithSourceID("ami-12345"))
	fmt.Printf("%#v", *hcimage)
	// Output:
	// image.Image{ImageID:"service-id-123", ProviderName:"happycloud", ProviderRegion:"west", Labels:map[string]string{}, SourceImageID:"ami-12345"}
}

func ExampleSetLabels() {
	a := &simpleArtifact{
		image_id: "service-id-123",
	}

	hcimage, _ := image.FromArtifact(a, image.SetLabels(map[string]interface{}{"kernel": "4.0", "python": "3.5"}))
	fmt.Printf("%v", hcimage.Labels)
	// Unordered output:
	// map[kernel:4.0 python:3.5]
}
