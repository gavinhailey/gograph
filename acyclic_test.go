package gograph

import (
	"reflect"
	"testing"
)

func TestTopologySort(t *testing.T) {
	// Create a dag with 6 vertices and 6 edges
	g := New[int](Acyclic())

	if !g.IsDirected() {
		t.Error(testErrMsgNotTrue)
	}

	if !g.IsAcyclic() {
		t.Error(testErrMsgNotTrue)
	}

	v1 := g.AddVertexByLabel(1)
	v2 := g.AddVertexByLabel(2)
	v3 := g.AddVertexByLabel(3)
	v4 := g.AddVertexByLabel(4)
	v5 := g.AddVertexByLabel(5)
	v6 := g.AddVertexByLabel(6)

	_, err := g.AddEdge(v1, v2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = g.AddEdge(v2, v3)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = g.AddEdge(v2, v4)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = g.AddEdge(v2, v5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = g.AddEdge(v3, v5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = g.AddEdge(v4, v6)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = g.AddEdge(v5, v6)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Perform a topological sort
	sortedVertices, err := TopologySort[int](g)

	// Check that there was no error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check that the sorted order is correct
	expectedOrder := []*Vertex[int]{v1, v2, v3, v4, v5, v6}
	if !reflect.DeepEqual(sortedVertices, expectedOrder) {
		t.Errorf("unexpected sort order. Got %v, expected %v", sortedVertices, expectedOrder)
	}
}

func TestStableTopologySort(t *testing.T) {
	// Create a graph where multiple valid topological sorts are possible
	g := New[int](Acyclic())

	// Create vertices with labels that we'll sort by
	v1 := g.AddVertexByLabel(1)
	v2 := g.AddVertexByLabel(2)
	v3 := g.AddVertexByLabel(3)
	v4 := g.AddVertexByLabel(4)
	v5 := g.AddVertexByLabel(5)
	v6 := g.AddVertexByLabel(6)

	// Add edges to create the graph structure
	// 1 -> 2
	// |    |
	// v    v
	// 3    4
	// |    |
	// v    v
	// 5 -> 6
	_, _ = g.AddEdge(v1, v2)
	_, _ = g.AddEdge(v1, v3)
	_, _ = g.AddEdge(v2, v4)
	_, _ = g.AddEdge(v3, v5)
	_, _ = g.AddEdge(v4, v6)
	_, _ = g.AddEdge(v5, v6)

	// Test with standard less function
	sortedVertices, err := StableTopologySort(g, func(a, b int) bool {
		return a < b
	})

	// Check that there was no error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// With this graph structure and sorting by label ascending,
	// we should get [1, 2, 3, 4, 5, 6]
	expectedOrder := []*Vertex[int]{v1, v2, v3, v4, v5, v6}
	if !reflect.DeepEqual(sortedVertices, expectedOrder) {
		t.Errorf("unexpected sort order with ascending labels. Got %v, expected %v",
			extractLabels(sortedVertices), extractLabels(expectedOrder))
	}

	// Test with reverse ordering function
	sortedVerticesReverse, err := StableTopologySort(g, func(a, b int) bool {
		return a > b
	})

	// Check that there was no error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// With this graph structure and sorting by label descending,
	// the topological constraints still apply, so we should get [1, 3, 5, 2, 4, 6]
	// This is because:
	// - v1 has no incoming edges, so it's still first
	// - Between v2 and v3 (both depend only on v1), v3 comes first due to higher label
	// - v5 depends only on v3, so it comes next
	// - v2 then comes after v5
	// - v4 depends on v2
	// - v6 depends on both v4 and v5
	expectedReverseOrder := []*Vertex[int]{v1, v3, v5, v2, v4, v6}
	if !reflect.DeepEqual(sortedVerticesReverse, expectedReverseOrder) {
		t.Errorf("unexpected sort order with descending labels. Got %v, expected %v",
			extractLabels(sortedVerticesReverse), extractLabels(expectedReverseOrder))
	}
}

// TestStableTopologySortWithTies tests the StableTopologySort function with tie resolution
func TestStableTopologySortWithTies(t *testing.T) {
	// Create a graph where multiple vertices have the same in-degree at the same time
	g := New[string](Acyclic())

	// Create vertices with the same starting in-degree
	vA := g.AddVertexByLabel("A")
	vB := g.AddVertexByLabel("B")
	vC := g.AddVertexByLabel("C")
	vD := g.AddVertexByLabel("D")
	vE := g.AddVertexByLabel("E")

	// Add edges (A->C, B->C, C->D, C->E)
	_, _ = g.AddEdge(vA, vC)
	_, _ = g.AddEdge(vB, vC)
	_, _ = g.AddEdge(vC, vD)
	_, _ = g.AddEdge(vC, vE)

	// Test with alphabetical ordering
	sortedVertices, err := StableTopologySort(g, func(a, b string) bool {
		return a < b
	})

	// Check that there was no error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// With alphabetical ordering and this graph:
	// - A and B have no incoming edges, so they come first (A before B due to sort)
	// - C depends on both A and B
	// - D and E both depend only on C (D before E due to sort)
	expectedOrder := []*Vertex[string]{vA, vB, vC, vD, vE}
	if !reflect.DeepEqual(sortedVertices, expectedOrder) {
		t.Errorf("unexpected sort order with alphabetical sorting. Got %v, expected %v",
			extractLabels(sortedVertices), extractLabels(expectedOrder))
	}

	// Test with reverse alphabetical ordering
	sortedVerticesReverse, err := StableTopologySort(g, func(a, b string) bool {
		return a > b
	})

	// Check that there was no error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// With reverse alphabetical ordering:
	// - A and B still come first (B before A due to reverse sort)
	// - C depends on both A and B
	// - D and E both depend only on C (E before D due to reverse sort)
	expectedReverseOrder := []*Vertex[string]{vB, vA, vC, vE, vD}
	if !reflect.DeepEqual(sortedVerticesReverse, expectedReverseOrder) {
		t.Errorf("unexpected sort order with reverse alphabetical. Got %v, expected %v",
			extractLabels(sortedVerticesReverse), extractLabels(expectedReverseOrder))
	}
}

// TestStableTopologySortNilComparison ensures the function handles a nil comparison function gracefully
func TestStableTopologySortNilComparison(t *testing.T) {
	// Create a simple graph
	g := New[int](Acyclic())
	v1 := g.AddVertexByLabel(1)
	v2 := g.AddVertexByLabel(2)
	_, _ = g.AddEdge(v1, v2)

	// Test with nil comparison function
	_, err := StableTopologySort(g, nil)

	// Should not panic and should return a valid result
	if err != nil {
		t.Errorf("unexpected error with nil comparison function: %v", err)
	}
}

// TestStableTopologySortCycle tests that the function correctly detects cycles
func TestStableTopologySortCycle(t *testing.T) {
	// Create a graph with a cycle
	g := New[int]()
	v1 := g.AddVertexByLabel(1)
	v2 := g.AddVertexByLabel(2)
	v3 := g.AddVertexByLabel(3)

	_, _ = g.AddEdge(v1, v2)
	_, _ = g.AddEdge(v2, v3)
	_, _ = g.AddEdge(v3, v1) // Creates a cycle

	// Test with standard less function
	_, err := StableTopologySort(g, func(a, b int) bool {
		return a < b
	})

	// Should detect a cycle
	if err != ErrDAGHasCycle {
		t.Errorf("expected cycle error, got %v", err)
	}
}

// Helper function to extract labels from a slice of vertices for easier debugging
func extractLabels[T comparable](vertices []*Vertex[T]) []T {
	labels := make([]T, len(vertices))
	for i, v := range vertices {
		labels[i] = v.Label()
	}
	return labels
}
