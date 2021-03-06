package result

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tektoncd/experimental/results/pkg/api/server/db"
	pb "github.com/tektoncd/experimental/results/proto/v1alpha2/results_go_proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestParseName(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		// if want is nil, assume error
		want []string
	}{
		{
			name: "simple",
			in:   "a/results/b",
			want: []string{"a", "b"},
		},
		{
			name: "resource name reuse",
			in:   "results/results/results",
			want: []string{"results", "results"},
		},
		{
			name: "missing name",
			in:   "a/results/",
		},
		{
			name: "missing name, no slash",
			in:   "a/results",
		},
		{
			name: "missing parent",
			in:   "/results/b",
		},
		{
			name: "missing parent, no slash",
			in:   "results/b",
		},
		{
			name: "wrong resource",
			in:   "a/record/b",
		},
		{
			name: "invalid parent",
			in:   "a/b/results/c",
		},
		{
			name: "invalid name",
			in:   "a/results/b/c",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			parent, name, err := ParseName(tc.in)
			if err != nil {
				if tc.want == nil {
					// error was expected, continue
					return
				}
				t.Fatal(err)
			}
			if tc.want == nil {
				t.Fatalf("expected error, got: [%s, %s]", parent, name)
			}

			if parent != tc.want[0] && name != tc.want[1] {
				t.Errorf("want: %v, got: [%s, %s]", tc.want, parent, name)
			}
		})
	}
}

func TestToStorage(t *testing.T) {
	got, err := ToStorage("foo", "bar", &pb.Result{
		Name: "foo/results/bar",
		Id:   "a",

		// These fields are ignored for now.
		CreatedTime: timestamppb.Now(),
		Annotations: map[string]string{"a": "b"},
		Etag:        "tacocat",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := &db.Result{
		Parent: "foo",
		Name:   "bar",
		ID:     "a",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("-want,+got: %s", diff)
	}
}

func TestToAPI(t *testing.T) {
	got := ToAPI(&db.Result{
		Parent: "foo",
		Name:   "bar",
		ID:     "a",
	})
	want := &pb.Result{
		Name: "foo/results/bar",
		Id:   "a",
	}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("-want,+got: %s", diff)
	}
}
