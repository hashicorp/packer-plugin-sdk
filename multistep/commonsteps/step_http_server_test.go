//go:generate mapstructure-to-hcl2 -type TestHttpConfig

package commonsteps

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

func TestStepHTTPServer_Run(t *testing.T) {

	tests := []struct {
		cfg         *HTTPConfig
		want        multistep.StepAction
		wantPort    interface{}
		wantContent map[string]string
	}{
		{
			&HTTPConfig{},
			multistep.ActionContinue,
			0,
			nil,
		},
		{
			&HTTPConfig{HTTPDir: "unknown_folder"},
			multistep.ActionHalt,
			nil,
			nil,
		},
		{
			&HTTPConfig{HTTPDir: "test-fixtures", HTTPPortMin: 9000},
			multistep.ActionContinue,
			9000,
			map[string]string{
				"Somedir/myfile.txt": "",
			},
		},
		{
			&HTTPConfig{HTTPContent: map[string]string{"/foo.txt": "biz", "/foo/bar.txt": "baz"}, HTTPPortMin: 9001},
			multistep.ActionContinue,
			9001,
			map[string]string{
				"foo.txt":      "biz",
				"/foo.txt":     "biz",
				"foo/bar.txt":  "baz",
				"/foo/bar.txt": "baz",
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.cfg), func(t *testing.T) {
			s := StepHTTPServerFromHTTPConfig(tt.cfg)
			state := testState(t)
			got := s.Run(context.Background(), state)
			defer s.Cleanup(state)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StepHTTPServer.Run() = %s, want %s", got, tt.want)
			}
			gotPort := state.Get("http_port")
			if !reflect.DeepEqual(gotPort, tt.wantPort) {
				t.Errorf("StepHTTPServer.Run() unexpected port = %v, want %v", gotPort, tt.wantPort)
			}
			for k, wantResponse := range tt.wantContent {
				resp, err := http.Get(fmt.Sprintf("http://:%d/%s", gotPort, k))
				if err != nil {
					t.Fatalf("http.Get: %v", err)
				}
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("readall: %v", err)
				}
				gotResponse := string(b)
				if diff := cmp.Diff(wantResponse, gotResponse); diff != "" {
					t.Fatalf("Unexpected %q content: %s", k, diff)
				}
			}
		})
	}
}
