package cloudant

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRevsDiff(t *testing.T) {
	database, err := makeDatabase()
	defer func() {
		fmt.Printf("Deleting database %s", database.Name)
		database.client.Delete(database.Name)
	}()

	if err != nil {
		t.Errorf("failed to create database")
	}

	doc := &struct {
		ID  string `json:"_id"`
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}{
		"doc-new",
		"mydata",
		57,
	}

	rev1, err1 := database.Set(doc)
	if err1 != nil {
		t.Error("failed to create document")
		return
	}
	if !strings.HasPrefix(rev1, "1-") {
		t.Error("got unexpected revision on create")
		return
	}

	// Note: lame attempt to close inconsistency window
	time.Sleep(500 * time.Millisecond)

	doc2 := &struct {
		ID  string `json:"_id"`
		Rev string `json:"_rev"`
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}{
		"doc-new",
		rev1,
		"mydata",
		57,
	}

	rev2, err2 := database.Set(doc2)
	if err2 != nil {
		t.Error("failed to update document")
	}

	// Note: lame attempt to close inconsistency window
	time.Sleep(500 * time.Millisecond)

	fakeRev := "3-b6b61a4f380712142ea80c90f172cc1e"
	rd := RevsDiffRequestBody{}
	rd["doc-new"] = []string{rev1, rev2, fakeRev}

	missing, err := database.RevsDiff(&rd)
	if err != nil {
		t.Errorf("RevsDiff failed %s", err)
	}

	for _, revs := range *missing {
		for _, rev := range revs.Missing {
			if rev != fakeRev {
				t.Errorf("RevsDiff expcted %s but found %s", fakeRev, rev)
			}
		}
	}
}

func TestBulkGet(t *testing.T) {
	database, err := makeDatabase()
	defer func() {
		fmt.Printf("Deleting database %s", database.Name)
		database.client.Delete(database.Name)
	}()

	if err != nil {
		t.Errorf("failed to create database")
	}

	doc := &struct {
		ID  string `json:"_id"`
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}{
		"doc-new",
		"mydata",
		57,
	}

	rev1, err1 := database.Set(doc)
	if err1 != nil {
		t.Error("failed to create document")
		return
	}
	if !strings.HasPrefix(rev1, "1-") {
		t.Error("got unexpected revision on create")
		return
	}

	// Note: lame attempt to close inconsistency window
	time.Sleep(500 * time.Millisecond)

	doc2 := &struct {
		ID  string `json:"_id"`
		Rev string `json:"_rev"`
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}{
		"doc-new",
		rev1,
		"mydata",
		57,
	}

	rev2, err2 := database.Set(doc2)
	if err2 != nil {
		t.Error("failed to update document")
	}

	// Note: lame attempt to close inconsistency window
	time.Sleep(500 * time.Millisecond)

	bg := &BulkGetRequest{}
	bg.Add("doc-new", rev1)
	bg.Add("doc-new", rev2)

	data, err := database.BulkGet(bg, true)

	fmt.Printf("%+v\n", data)
}
