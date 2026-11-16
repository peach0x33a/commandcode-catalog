package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestJSONArrayAt(t *testing.T) {
	body := []byte(`xx[{"name":"a]b","nested":[1,2],"quote":"say \"hi\""},{"id":"c"}]yy`)
	got, err := jsonArrayAt(body, 2)
	if err != nil {
		t.Fatalf("jsonArrayAt: %v", err)
	}
	if want := `[{"name":"a]b","nested":[1,2],"quote":"say \"hi\""},{"id":"c"}]`; string(got) != want {
		t.Fatalf("jsonArrayAt = %s, want %s", got, want)
	}
}

func TestFindModelsArraySkipsUnrelatedKeys(t *testing.T) {
	body := []byte(`{"models":{"count":2},"meta":{"models":[{"id":"decoy"}]},` +
		`"data":{"models":[{"id":"a","name":"A","category":"opensource","planAllowanceUsd":{"go":6,"goat":20,"pro":30}}]}}`)
	got, err := findModelsArray(body)
	if err != nil {
		t.Fatalf("findModelsArray: %v", err)
	}
	var models []upstreamModel
	if err := json.Unmarshal(got, &models); err != nil {
		t.Fatalf("decoded array is not valid JSON: %v", err)
	}
	if len(models) != 1 || models[0].ID != "a" {
		t.Fatalf("got %+v, want the catalog entry", models)
	}
}

func TestCatalogFromKeepsPlanAllowances(t *testing.T) {
	body := catalogPayload(t, 60)
	models, err := catalogFrom(body)
	if err != nil {
		t.Fatalf("catalogFrom: %v", err)
	}
	if len(models) != 60 {
		t.Fatalf("got %d models, want 60", len(models))
	}
	first := models[0]
	if first.ID != "model-00" || first.Name != "Model 00" || first.Category != "opensource" {
		t.Fatalf("unexpected first row: %+v", first)
	}
	if first.Go != 6 || first.Goat != 20 || first.Pro != 30 {
		t.Fatalf("unexpected allowances: %+v", first)
	}
}

func TestModelsFromRejectsMissingAllowance(t *testing.T) {
	upstream := []upstreamModel{
		{ID: "x", Name: "X", Category: "opensource", PlanAllowanceUSD: map[string]float64{"go": 6, "goat": 20}},
	}
	if _, err := modelsFrom(upstream); err == nil {
		t.Fatal("modelsFrom accepted a model without a pro allowance")
	}
}

func TestRenderIsValidJSON(t *testing.T) {
	models := []catalogModel{
		{ID: "b", Name: "B", Category: "premium", Go: 6, Goat: 20, Pro: 20},
		{ID: "a", Name: "A", Category: "opensource", Go: 6, Goat: 20, Pro: 30},
	}
	data, err := render(models, "2026-11-16")
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	var got catalogFile
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("render wrote invalid JSON: %v", err)
	}
	if got.Schema != schemaVersion || got.Source != sourceURL || got.Captured != "2026-11-16" {
		t.Fatalf("unexpected header: %+v", got)
	}
	if len(got.Models) != 2 || got.Models[0].ID != "b" {
		t.Fatalf("unexpected models: %+v", got.Models)
	}
}

// catalogPayload builds an RSC-shaped body with n valid models.
func catalogPayload(t *testing.T, n int) []byte {
	t.Helper()
	var b strings.Builder
	b.WriteString(`1:{"models":[`)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(fmt.Sprintf(`{"id":"model-%02d","name":"Model %02d","category":"opensource","planAllowanceUsd":{"go":6,"goat":20,"pro":30}}`, i, i))
	}
	b.WriteString(`]}`)
	return []byte(b.String())
}
