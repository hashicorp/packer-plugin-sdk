package image_test

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
)

type artifact struct {
	images map[string]string
}

func (a *artifact) BuilderId() string {
	return "example.happycloud"
}

func (a *artifact) Files() []string {
	return nil
}

func (a *artifact) Id() string {
	parts := make([]string, 0, len(a.images))
	for loc, id := range a.images {
		parts = append(parts, loc+":"+id)
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func (a *artifact) String() string {
	return a.Id()
}

func (a *artifact) State(name string) interface{} {
	return nil
}

func (a *artifact) Destroy() error {
	return nil
}

func ExampleFromMappedData() {
	a := &artifact{
		images: map[string]string{
			"west": "happycloud-1",
			"east": "happycloud-2",
		},
	}

	f := func(key, value interface{}) (*image.Image, error) {
		v, ok := value.(string)
		if !ok {
			return nil, errors.New("for happycloud maps value should always be string")
		}
		k, ok := key.(string)
		if !ok {
			return nil, errors.New("for happycloud maps key should always be string")
		}

		img := image.Image{ProviderName: "happycloud", ProviderRegion: k, ImageID: v}
		return &img, nil
	}

	hcimages, _ := image.FromMappedData(a.images, f)
	for _, hcimage := range hcimages {
		fmt.Printf("%#v\n", *hcimage)
	}
	// Unordered output:
	// image.Image{ImageID:"happycloud-1", ProviderName:"happycloud", ProviderRegion:"west", Labels:map[string]string(nil)}
	// image.Image{ImageID:"happycloud-2", ProviderName:"happycloud", ProviderRegion:"east", Labels:map[string]string(nil)}
}
