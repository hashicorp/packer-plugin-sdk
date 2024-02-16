// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/rpc"
	"reflect"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	"google.golang.org/protobuf/proto"
)

// commonClient allows to rpc call funcs that can be defined on the different
// build blocks of Packer.
type commonClient struct {
	// endpoint is usually the type of build block we are connecting to.
	//
	// eg: Provisioner / PostProcessor / Builder / Artifact / Communicator
	endpoint string
	client   *rpc.Client
	mux      *muxBroker

	// useProto lets us determine whether or not we should use protobuf for serialising
	// data over RPC instead of gob.
	//
	// This is controlled by Packer using the `--use-proto` flag on plugin commands.
	useProto bool
}

type commonServer struct {
	mux *muxBroker
	// a HCL2 enabled component such as a Builder
	selfConfigurable interface {
		ConfigSpec() hcldec.ObjectSpec
	}

	// useProto lets us determine whether or not we should use protobuf for serialising
	// data over RPC instead of gob.
	//
	// This is controlled by Packer using the `--use-proto` flag on plugin commands.
	useProto bool
}

type ConfigSpecResponse struct {
	ConfigSpec []byte
}

func (p *commonClient) ConfigSpec() hcldec.ObjectSpec {
	// TODO(azr): the RPC Call can fail but the ConfigSpec signature doesn't
	// return an error; should we simply panic ? Logging this for now; will
	// decide later. The correct approach would probably be to return an error
	// in ConfigSpec but that will break a lot of things.
	resp := &ConfigSpecResponse{}
	cerr := p.client.Call(p.endpoint+".ConfigSpec", new(interface{}), resp)
	if cerr != nil {
		err := fmt.Errorf("ConfigSpec failed: %v", cerr)
		panic(err.Error())
	}

	// Legacy: this will need to be removed when we discontinue gob-encoding
	//
	// This is required for backwards compatibility for now, but using
	// gob to encode the spec objects will fail against the upstream cty
	// library, since they removed support for it.
	//
	// This will be a breaking change, as older plugins won't be able to
	// communicate with Packer any longer.
	if !p.useProto {
		log.Printf("[DEBUG] - common: receiving ConfigSpec as gob")
		res := hcldec.ObjectSpec{}
		err := gob.NewDecoder(bytes.NewReader(resp.ConfigSpec)).Decode(&res)
		if err != nil {
			panic(fmt.Errorf("failed to decode HCL spec from gob: %s", err))
		}
		return res
	}

	log.Printf("[DEBUG] - common: receiving ConfigSpec as protobuf")
	spec, err := protobufToHCL2Spec(resp.ConfigSpec)
	if err != nil {
		panic(err)
	}

	return spec
}

func (s *commonServer) ConfigSpec(_ interface{}, reply *ConfigSpecResponse) error {
	spec := s.selfConfigurable.ConfigSpec()

	if !s.useProto {
		log.Printf("[DEBUG] - common: sending ConfigSpec as gob")
		b := &bytes.Buffer{}
		err := gob.NewEncoder(b).Encode(spec)
		if err != nil {
			return fmt.Errorf("failed to encode spec from gob: %s", err)
		}
		reply.ConfigSpec = b.Bytes()

		return nil
	}

	log.Printf("[DEBUG] - common: sending ConfigSpec as protobuf")
	rawBytes, err := hcl2SpecToProtobuf(spec)
	if err != nil {
		return fmt.Errorf("failed to encode HCL spec from protobuf: %s", err)
	}
	reply.ConfigSpec = rawBytes

	return nil
}

// hcl2SpecToProtobuf converts a hcldec.ObjectSpec to a protobuf-serialised
// byte array so it can then be used to send to a Plugin/Packer.
func hcl2SpecToProtobuf(spec hcldec.ObjectSpec) ([]byte, error) {
	ret, err := ToProto(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to convert hcldec.Spec to hclspec.Spec: %s", err)
	}
	rawBytes, err := proto.Marshal(ret)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise hclspec.Spec to protobuf: %s", err)
	}

	return rawBytes, nil
}

// protobufToHCL2Spec converts a protobuf-encoded spec to a usable hcldec.Spec.
func protobufToHCL2Spec(serData []byte) (hcldec.ObjectSpec, error) {
	confSpec := &Spec{}
	err := proto.Unmarshal(serData, confSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal hclspec.Spec from raw protobuf: %q", err)
	}
	spec, err := confSpec.FromProto()
	if err != nil {
		return nil, fmt.Errorf("failed to decode HCL spec: %q", err)
	}

	obj, ok := spec.(*hcldec.ObjectSpec)
	if !ok {
		return nil, fmt.Errorf("decoded HCL spec is not an object spec: %s", reflect.TypeOf(spec).String())
	}

	return *obj, nil
}

func init() {
	gob.Register(new(hcldec.AttrSpec))
	gob.Register(new(hcldec.BlockSpec))
	gob.Register(new(hcldec.BlockAttrsSpec))
	gob.Register(new(hcldec.BlockListSpec))
	gob.Register(new(hcldec.BlockObjectSpec))
	gob.Register(new(cty.Value))
}
