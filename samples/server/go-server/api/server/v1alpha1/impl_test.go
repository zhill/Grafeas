// Copyright 2017 The Grafeas Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"fmt"
	"golang.org/x/net/context"
	"reflect"
	"testing"

	"github.com/grafeas/grafeas/samples/server/go-server/api/server/name"
	"github.com/grafeas/grafeas/samples/server/go-server/api/server/storage"
	"github.com/grafeas/grafeas/samples/server/go-server/api/server/testing"
	pb "github.com/grafeas/grafeas/v1alpha1/proto"
	opspb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateOperation(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	op := &opspb.Operation{}
	req := pb.CreateOperationRequest{Parent: "projects/opp", Operation: op}
	if _, err := g.CreateOperation(ctx, &req); err == nil {
		t.Error("CreateOperation(empty operation): got success, want error")
	} else if s, _ := status.FromError(err); s.Code() != codes.InvalidArgument {
		t.Errorf("CreateOperation(empty operation): got %v, want InvalidArgument", err)
	}
	op = testutil.Operation()
	pID, _, aErr := name.ParseOperation(op.Name)
	if aErr != nil {
		t.Errorf("Could not parse projectID from op.Name: %v", op.Name)
	}
	parent := name.FormatProject(pID)
	req = pb.CreateOperationRequest{Parent: parent, Operation: op}
	if _, err := g.CreateOperation(ctx, &req); error(err) != nil {
		t.Errorf("CreateOperation(%v) got %#v, want success", op, err)
	}
}

func TestCreateOccurrence(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	n := testutil.Note()
	pID, _, aErr := name.ParseNote(n.Name)
	if aErr != nil {
		t.Errorf("Could not parse projectID from n.Name: %v", n.Name)
	}
	parent := name.FormatProject(pID)
	req := &pb.CreateNoteRequest{Parent: parent, Note: n}
	if _, err := g.CreateNote(ctx, req); err != nil {
		t.Fatalf("CreateNote(%v) got %v, want success", req, err)
	}
	oReq := &pb.CreateOccurrenceRequest{Parent: parent, Occurrence: &pb.Occurrence{}}
	if _, err := g.CreateOccurrence(ctx, oReq); err == nil {
		t.Error("CreateOccurrence(empty occ): got success, want error")
	} else if s, _ := status.FromError(err); s.Code() != codes.InvalidArgument {
		t.Errorf("CreateOccurrence(empty occ): got %v, want InvalidArgument)", err)
	}
	o := testutil.Occurrence(n.Name)
	pID, _, aErr = name.ParseOccurrence(o.Name)
	if aErr != nil {
		t.Errorf("Could not parse projectID from o.Name: %v", o.Name)
	}
	parent = name.FormatProject(pID)
	oReq = &pb.CreateOccurrenceRequest{Parent: parent, Occurrence: o}
	if _, err := g.CreateOccurrence(ctx, oReq); err != nil {
		t.Errorf("CreateOccurrence(%v) got %v, want success", oReq, err)
	}
	// Try to insert an occurrence for a note that does not exist.
	o.Name = "projects/testproject/occurrences/nonote"
	o.NoteName = "projects/scan-provider/notes/notthere"
	oReq = &pb.CreateOccurrenceRequest{Parent: "projects/testproject", Occurrence: o}
	if _, err := g.CreateOccurrence(ctx, oReq); err == nil {
		t.Errorf("CreateOccurrence got success, want Error")
	} else if s, _ := status.FromError(err); s.Code() != codes.InvalidArgument {
		t.Errorf("CreateOccurrence got code %v want %v", err, codes.InvalidArgument)
	}
}

func TestCreateNote(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	n := &pb.Note{}
	req := &pb.CreateNoteRequest{Parent: "projects/foo", Note: n}
	if _, err := g.CreateNote(ctx, req); err == nil {
		t.Error("CreateNote(empty note): got success, want error")
	} else if s, _ := status.FromError(err); s.Code() != codes.InvalidArgument {
		t.Errorf("CreateNote(empty note): got %v, want %v", err, codes.InvalidArgument)
	}
	n = testutil.Note()
	pID, _, aErr := name.ParseNote(n.Name)
	if aErr != nil {
		t.Errorf("Could not parse projectID from n.Name: %v", n.Name)
	}
	parent := name.FormatProject(pID)
	req = &pb.CreateNoteRequest{Parent: parent, Note: n}
	if _, err := g.CreateNote(ctx, req); err != nil {
		t.Errorf("CreateNote(%v) got %v, want success", n, err)
	}
}

func TestDeleteNote(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	n := testutil.Note()
	pID, _, err := name.ParseNote(n.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}
	req := &pb.DeleteNoteRequest{Name: n.Name}
	if _, err := g.DeleteNote(ctx, req); err == nil {
		t.Error("DeleteNote that doesn't exist got success, want err")
	}
	parent := name.FormatProject(pID)
	cReq := &pb.CreateNoteRequest{Parent: parent, Note: n}
	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Errorf("CreateNote(%v) got %v, want success", n, err)
	}
	if _, err := g.DeleteNote(ctx, req); err != nil {
		t.Errorf("DeleteNote  got %v, want success", err)
	}
}

func TestDeleteOccurrence(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	n := testutil.Note()
	pID, _, err := name.ParseNote(n.Name)
	parent := name.FormatProject(pID)
	cReq := &pb.CreateNoteRequest{Parent: parent, Note: n}
	// CreateNote so we can create an occurrence
	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Fatalf("CreateNote(%v) got %v, want success", n, err)
	}
	o := testutil.Occurrence(n.Name)

	pID, _, err = name.ParseOccurrence(o.Name)
	if err != nil {
		t.Fatalf("Error parsing occurrence name %v", err)
	}
	parent = name.FormatProject(pID)
	oReq := &pb.CreateOccurrenceRequest{Parent: parent, Occurrence: o}
	if _, err := g.CreateOccurrence(ctx, oReq); err != nil {
		t.Fatalf("CreateOccurrence(%v) got %v, want success", n, err)
	}
	dReq := &pb.DeleteOccurrenceRequest{Name: o.Name}
	if _, err := g.DeleteOccurrence(ctx, dReq); err != nil {
		t.Errorf("DeleteOccurrence  got %v, want success", err)
	}
}

func TestDeleteOperation(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	o := testutil.Operation()
	pID, _, err := name.ParseOperation(o.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}
	req := &opspb.DeleteOperationRequest{Name: o.Name}
	if _, err := g.DeleteOperation(ctx, req); err == nil {
		t.Error("DeleteOperation that doesn't exist got success, want err")
	}
	parent := name.FormatProject(pID)
	cReq := &pb.CreateOperationRequest{Parent: parent, Operation: o}
	if _, err := g.CreateOperation(ctx, cReq); err != nil {
		t.Fatalf("CreateOperation(%v) got %v, want success", o, err)
	}
	if _, err := g.DeleteOperation(ctx, req); err != nil {
		t.Errorf("DeleteOperation  got %v, want success", err)
	}
}

func TestGetNote(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	n := testutil.Note()
	pID, _, err := name.ParseNote(n.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}
	req := &pb.GetNoteRequest{Name: n.Name}
	if _, err := g.GetNote(ctx, req); err == nil {
		t.Error("GetNote that doesn't exist got success, want err")
	}
	if err != nil {
		t.Errorf("Could not parse projectID from n.Name: %v", n.Name)
	}
	parent := name.FormatProject(pID)
	cReq := &pb.CreateNoteRequest{Parent: parent, Note: n}
	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Fatalf("CreateNote(%v) got %v, want success", n, err)
	}
	if got, err := g.GetNote(ctx, req); err != nil {
		t.Fatalf("GetNote(%v) got %v, want success", n, err)
	} else if n.Name != got.Name || !reflect.DeepEqual(n.GetVulnerabilityType(), got.GetVulnerabilityType()) {
		t.Errorf("GetNote got %v, want %v", *got, n)
	}
}

func TestGetOccurrence(t *testing.T) {
	ctx := context.Background()

	g := Grafeas{storage.NewMemStore()}
	n := testutil.Note()
	pID, _, err := name.ParseNote(n.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}

	o := testutil.Occurrence(n.Name)

	req := &pb.GetOccurrenceRequest{Name: o.Name}
	if _, err := g.GetOccurrence(ctx, req); err == nil {
		t.Error("GetOccurrence that doesn't exist got success, want err")
	}
	parent := name.FormatProject(pID)
	cReq := &pb.CreateNoteRequest{Parent: parent, Note: n}
	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Fatalf("CreateNote(%v) got %v, want success", n, err)
	}
	opID, _, err := name.ParseOccurrence(o.Name)
	if err != nil {
		t.Fatalf("Error parsing occurrence name %v", err)
	}
	oParent := name.FormatProject(opID)
	ocReq := &pb.CreateOccurrenceRequest{Parent: oParent, Occurrence: o}
	if _, err := g.CreateOccurrence(ctx, ocReq); err != nil {
		t.Fatalf("CreateOccurrence(%v) got %v, want success", n, err)
	}
	if got, err := g.GetOccurrence(ctx, req); err != nil {
		t.Fatalf("GetOccurrence(%v) got %v, want success", o, err)
	} else if o.Name != got.Name || !reflect.DeepEqual(o.GetVulnerabilityDetails(), got.GetVulnerabilityDetails()) {
		t.Errorf("GetOccurrence got %v, want %v", *got, o)
	}
}

func TestGetOperation(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	o := testutil.Operation()
	pID, _, err := name.ParseOperation(o.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}
	req := &opspb.GetOperationRequest{Name: o.Name}
	if _, err := g.GetOperation(ctx, req); err == nil {
		t.Error("GetOperation that doesn't exist got success, want err")
	}
	parent := name.FormatProject(pID)
	cReq := &pb.CreateOperationRequest{Parent: parent, Operation: o}
	if _, err := g.CreateOperation(ctx, cReq); err != nil {
		t.Fatalf("CreateOperation(%v) got %v, want success", o, err)
	}
	if got, err := g.GetOperation(ctx, req); err != nil {
		t.Fatalf("GetOperation(%v) got %v, want success", o, err)
	} else if o.Name != got.Name || !reflect.DeepEqual(got, o) {
		t.Errorf("GetOperation got %v, want %v", *got, o)
	}
}

func TestGetOccurrenceNote(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	n := testutil.Note()

	o := testutil.Occurrence(n.Name)

	req := &pb.GetOccurrenceNoteRequest{Name: o.Name}
	if _, err := g.GetOccurrenceNote(ctx, req); err == nil {
		t.Error("GetOccurrenceNote that doesn't exist got success, want err")
	}
	pID, _, err := name.ParseNote(n.Name)
	if err != nil {
		t.Fatalf("Error parsing occurrence name %v", err)
	}
	parent := name.FormatProject(pID)
	cReq := &pb.CreateNoteRequest{Parent: parent, Note: n}

	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Fatalf("CreateNote(%v) got %v, want success", n, err)
	}
	pID, _, err = name.ParseOccurrence(o.Name)
	if err != nil {
		t.Fatalf("Error parsing occurrence name %v", err)
	}
	parent = name.FormatProject(pID)
	coReq := &pb.CreateOccurrenceRequest{Parent: parent, Occurrence: o}
	if _, err := g.CreateOccurrence(ctx, coReq); err != nil {
		t.Fatalf("CreateOccurrence(%v) got %v, want success", n, err)
	}
	if got, err := g.GetOccurrenceNote(ctx, req); err != nil {
		t.Fatalf("GetOccurrenceNote(%v) got %v, want success", n, err)
	} else if n.Name != got.Name || !reflect.DeepEqual(n.GetVulnerabilityType(), got.GetVulnerabilityType()) {
		t.Errorf("GetOccurrenceNote got %v, want %v", *got, n)
	}
}

func TestUpdateNote(t *testing.T) {
	ctx := context.Background()
	// Update Note that doesn't exist
	updateDesc := "this is a new description"
	g := Grafeas{storage.NewMemStore()}
	n := testutil.Note()
	pID, _, err := name.ParseNote(n.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}
	update := testutil.Note()
	update.LongDescription = updateDesc
	req := &pb.UpdateNoteRequest{Name: n.Name, Note: n}
	if _, err := g.UpdateNote(ctx, req); err != nil {
		t.Error("UpdateNote that doesn't exist got success, want err")
	}

	parent := name.FormatProject(pID)
	cReq := &pb.CreateNoteRequest{Parent: parent, Note: n}
	// Actually create note
	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Fatalf("CreateNote(%v) got %v, want success", n, err)
	}

	// Update Note name and fail
	update.Name = "New name"
	req = &pb.UpdateNoteRequest{Name: n.Name, Note: update}
	if _, err := g.UpdateNote(ctx, req); err == nil {
		t.Error("UpdateNote that with name change got success, want err")
	}

	// Update Note and verify that update worked.
	update = testutil.Note()
	update.LongDescription = updateDesc
	req = &pb.UpdateNoteRequest{Name: n.Name, Note: update}
	if got, err := g.UpdateNote(ctx, req); err != nil {
		t.Errorf("UpdateNote got %v, want success", err)
	} else if updateDesc != update.LongDescription {
		t.Errorf("UpdateNote got %v, want %v",
			got.LongDescription, updateDesc)
	}
	if got, err := g.GetNote(ctx, &pb.GetNoteRequest{Name: n.Name}); err != nil {
		t.Fatalf("GetNote(%v) got %v, want success", n, err)
	} else if updateDesc != got.LongDescription {
		t.Errorf("GetNote got %v, want %v", got.LongDescription, updateDesc)
	}
}

func TestUpdateOccurrence(t *testing.T) {
	ctx := context.Background()
	// Update occurrence that doesn't exist
	g := Grafeas{storage.NewMemStore()}
	n := testutil.Note()
	npID, _, err := name.ParseNote(n.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}
	nParent := name.FormatProject(npID)
	cReq := &pb.CreateNoteRequest{Parent: nParent, Note: n}

	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Fatalf("CreateNote(%v) got %v, want success", n, err)
	}
	o := testutil.Occurrence(n.Name)

	req := &pb.UpdateOccurrenceRequest{Name: o.Name, Occurrence: o}
	if _, err := g.UpdateOccurrence(ctx, req); err == nil {
		t.Error("UpdateOccurrence that doesn't exist got success, want err")
	}
	// Create an occurrence to update
	pID, _, err := name.ParseOccurrence(o.Name)
	if err != nil {
		t.Fatalf("Error parsing occurrence name %v", err)
	}
	parent := name.FormatProject(pID)
	ocReq := &pb.CreateOccurrenceRequest{Parent: parent, Occurrence: o}
	if _, err := g.CreateOccurrence(ctx, ocReq); err != nil {
		t.Fatalf("CreateOccurrence(%v) got %v, want success", n, err)
	}
	// update occurrence name
	update := testutil.Occurrence(n.Name)
	update.Name = "New name"
	req = &pb.UpdateOccurrenceRequest{Name: update.Name, Occurrence: update}
	if _, err := g.UpdateOccurrence(ctx, req); err == nil {
		t.Error("UpdateOccurrence with name change got success, want err")
	}

	// update note name to a note that doesn't exist
	update = testutil.Occurrence("projects/p/notes/bar")
	req = &pb.UpdateOccurrenceRequest{Name: o.Name, Occurrence: update}
	if _, err := g.UpdateOccurrence(ctx, req); err == nil {
		t.Error("UpdateOccurrence that with note name that doesn't exist" +
			" got success, want err")
	}

	// update note name to a note that does exist
	n = testutil.Note()
	newName := fmt.Sprintf("%v-new", n.Name)
	n.Name = newName

	cReq = &pb.CreateNoteRequest{Parent: nParent, Note: n}
	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Fatalf("CreateNote(%v) got %v, want success", n, err)
	}
	update = testutil.Occurrence(n.Name)
	req = &pb.UpdateOccurrenceRequest{Name: o.Name, Occurrence: update}
	if got, err := g.UpdateOccurrence(ctx, req); err != nil {
		t.Errorf("UpdateOccurrence got %v, want success", err)
	} else if n.Name != got.NoteName {
		t.Errorf("UpdateOccurrence got %v, want %v",
			got.NoteName, n.Name)
	}
	gReq := &pb.GetOccurrenceRequest{Name: o.Name}
	if got, err := g.GetOccurrence(ctx, gReq); err != nil {
		t.Fatalf("GetOccurrence(%v) got %v, want success", n, err)
	} else if n.Name != got.NoteName {
		t.Errorf("GetOccurrence got %v, want %v",
			got.NoteName, n.Name)
	}
}

func TestListOccurrences(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	n := testutil.Note()
	npID, _, err := name.ParseNote(n.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}
	nParent := name.FormatProject(npID)
	cReq := &pb.CreateNoteRequest{Parent: nParent, Note: n}

	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Fatalf("CreateNote(%v) got %v, want success", n, err)
	}
	os := []*pb.Occurrence{}
	findProject := "findThese"
	dontFind := "dontFind"
	for i := 0; i < 20; i++ {
		o := testutil.Occurrence(n.Name)
		if i < 5 {
			o.Name = name.FormatOccurrence(findProject, string(i))
		} else {
			o.Name = name.FormatOccurrence(dontFind, string(i))
		}
		if err != nil {
			t.Fatalf("Error parsing occurrence name %v", err)
		}
		pID, _, err := name.ParseOccurrence(o.Name)
		if err != nil {
			t.Fatalf("Error parsing occurrence name %v", err)
		}
		parent := name.FormatProject(pID)
		ocReq := &pb.CreateOccurrenceRequest{Parent: parent, Occurrence: o}
		if _, err := g.CreateOccurrence(ctx, ocReq); err != nil {
			t.Fatalf("CreateOccurrence got %v want success", err)
		}
		os = append(os, o)
	}

	lReq := &pb.ListOccurrencesRequest{Parent: name.FormatProject(findProject)}
	resp, lErr := g.ListOccurrences(ctx, lReq)
	if lErr != nil {
		t.Fatalf("ListOccurrences got %v want success", err)
	}
	if len(resp.Occurrences) != 5 {
		t.Errorf("resp.Occurrences got %d, want 5", len(resp.Occurrences))
	}
}

func TestListOperations(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	findProject := "findThese"
	dontFind := "dontFind"
	for i := 0; i < 20; i++ {
		o := testutil.Operation()
		if i < 5 {
			o.Name = name.FormatOperation(findProject, string(i))
		} else {
			o.Name = name.FormatOperation(dontFind, string(i))
		}
		pID, _, err := name.ParseOperation(o.Name)
		if err != nil {
			t.Fatalf("Error parsing note name %v", err)
		}
		parent := name.FormatProject(pID)
		cReq := &pb.CreateOperationRequest{Parent: parent, Operation: o}
		if _, err := g.CreateOperation(ctx, cReq); err != nil {
			t.Fatalf("CreateOperation(%v) got %v, want success", o, err)
		}
	}

	lReq := &opspb.ListOperationsRequest{Name: name.FormatProject(findProject)}
	resp, err := g.ListOperations(ctx, lReq)
	if err != nil {
		t.Fatalf("ListOperations got %v want success", err)
	}
	if len(resp.Operations) != 5 {
		t.Errorf("resp.Operations got %d, want 5", len(resp.Operations))
	}
}

func TestListNotes(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	findProject := "findThese"
	dontFind := "dontFind"
	for i := 0; i < 20; i++ {
		n := testutil.Note()
		if i < 5 {
			n.Name = name.FormatNote(findProject, string(i))
		} else {
			n.Name = name.FormatNote(dontFind, string(i))
		}
		npID, _, err := name.ParseNote(n.Name)
		if err != nil {
			t.Fatalf("Error parsing note name %v", err)
		}
		nParent := name.FormatProject(npID)
		cReq := &pb.CreateNoteRequest{Parent: nParent, Note: n}

		if _, err := g.CreateNote(ctx, cReq); err != nil {
			t.Fatalf("CreateNote(%v) got %v, want success", n, err)
		}
	}

	req := &pb.ListNotesRequest{Parent: name.FormatProject(findProject)}
	resp, err := g.ListNotes(ctx, req)
	if err != nil {
		t.Fatalf("ListNotes got %v want success", err)
	}
	if len(resp.Notes) != 5 {
		t.Errorf("resp.Notes got %d, want 5", len(resp.Notes))
	}
}

func TestListNoteOccurrences(t *testing.T) {
	ctx := context.Background()
	g := Grafeas{storage.NewMemStore()}
	n := testutil.Note()
	npID, _, err := name.ParseNote(n.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}
	nParent := name.FormatProject(npID)
	cReq := &pb.CreateNoteRequest{Parent: nParent, Note: n}
	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Fatalf("CreateNote(%v) got %v, want success", n, err)
	}
	findProject := "project"
	dontFind := "otherProject"
	for i := 0; i < 20; i++ {
		o := testutil.Occurrence(n.Name)
		if i < 5 {
			o.Name = name.FormatOccurrence(findProject, string(i))
		} else {
			o.Name = name.FormatOccurrence(dontFind, string(i))
		}

		pID, _, err := name.ParseOccurrence(o.Name)
		if err != nil {
			t.Fatalf("Error parsing occurrence name %v", err)
		}
		parent := name.FormatProject(pID)
		ocReq := &pb.CreateOccurrenceRequest{Parent: parent, Occurrence: o}
		if _, err := g.CreateOccurrence(ctx, ocReq); err != nil {
			t.Fatalf("CreateOccurrence got %v want success", err)
		}
	}
	// Create an occurrence tied to another note, to make sure we don't find it.
	otherN := testutil.Note()
	otherN.Name = "projects/np/notes/not-to-find"
	npID, _, err = name.ParseNote(otherN.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}
	nParent = name.FormatProject(npID)
	cReq = &pb.CreateNoteRequest{Parent: nParent, Note: otherN}
	if _, err := g.CreateNote(ctx, cReq); err != nil {
		t.Fatalf("CreateNote got %v want success", err)
	}
	o := testutil.Occurrence(otherN.Name)
	pID, _, err := name.ParseOccurrence(o.Name)
	if err != nil {
		t.Fatalf("Error parsing occurrence name %v", err)
	}
	parent := name.FormatProject(pID)
	ocReq := &pb.CreateOccurrenceRequest{Parent: parent, Occurrence: o}
	if _, err := g.CreateOccurrence(ctx, ocReq); err != nil {
		t.Fatalf("CreateOccurrence got %v want success", err)
	}
	pID, _, err = name.ParseNote(n.Name)
	if err != nil {
		t.Fatalf("Error parsing note name %v", err)
	}
	lReq := &pb.ListNoteOccurrencesRequest{Name: n.Name}
	resp, lErr := g.ListNoteOccurrences(ctx, lReq)
	if lErr != nil {
		t.Fatalf("ListNoteOccurrences got %v want success", err)
	}
	if len(resp.Occurrences) != 20 {
		t.Errorf("resp.Occurrences got %d, want 20", len(resp.Occurrences))
	}
}
