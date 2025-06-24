package traverse

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gavinhailey/gograph"
)

func TestBreadthFirstIterator(t *testing.T) {
	// Create a new graph
	g := gograph.New[string]()

	// the example graph
	//	A -> B -> C
	//	|    |    |
	//	v    v    v
	//	D -> E -> F

	vertices := map[string]*gograph.Vertex[string]{
		"A": g.AddVertexByLabel("A"),
		"B": g.AddVertexByLabel("B"),
		"C": g.AddVertexByLabel("C"),
		"D": g.AddVertexByLabel("D"),
		"E": g.AddVertexByLabel("E"),
		"F": g.AddVertexByLabel("F"),
	}

	// add some edges
	_, _ = g.AddEdge(vertices["A"], vertices["B"])
	_, _ = g.AddEdge(vertices["A"], vertices["D"])
	_, _ = g.AddEdge(vertices["B"], vertices["C"])
	_, _ = g.AddEdge(vertices["B"], vertices["E"])
	_, _ = g.AddEdge(vertices["C"], vertices["F"])
	_, _ = g.AddEdge(vertices["D"], vertices["E"])
	_, _ = g.AddEdge(vertices["E"], vertices["F"])

	// create an iterator with a vertex that doesn't exist
	_, err := NewBreadthFirstIterator(g, "X")
	if err == nil {
		t.Error("Expect NewBreadthFirstIterator returns error, but got nil")
	}

	// test depth first iteration
	iter, err := NewBreadthFirstIterator(g, "A")
	if err != nil {
		t.Errorf("Expect NewBreadthFirstIterator doesn't return error, but got %s", err)
	}

	expected := []string{"A", "B", "D", "C", "E", "F"}

	for i, label := range expected {
		if !iter.HasNext() {
			t.Errorf("Expected iter.HasNext() to be true, but it was false for label %s", label)
		}

		v := iter.Next()
		if v.Label() != expected[i] {
			t.Errorf("Expected iter.Next().Label() to be %s, but got %s", expected[i], v.Label())
		}
	}

	if iter.HasNext() {
		t.Error("Expected iter.HasNext() to be false, but it was true")
	}

	v := iter.Next()
	if v != nil {
		t.Errorf("Expected nil, but got %+v", v)
	}

	// test the Reset method
	iter.Reset()
	if !iter.HasNext() {
		t.Error("Expected iter.HasNext() to be true, but it was false after reset")
	}

	v = iter.Next()
	if v.Label() != "A" {
		t.Errorf("Expected iter.Next().Label() to be %s, but got %s", "A", v.Label())
	}

	// test Iterate method
	iter.Reset()
	var ordered []string
	err = iter.Iterate(
		func(vertex *gograph.Vertex[string]) error {
			ordered = append(ordered, vertex.Label())
			return nil
		},
	)
	if err != nil {
		t.Errorf("Expect iter.Iterate(func) returns no error, but got one %s", err)
	}

	if !reflect.DeepEqual(expected, ordered) {
		t.Errorf(
			"Expect same vertex order, but got different one expected: %v, actual: %v",
			expected, ordered,
		)
	}

	iter.Reset()
	expectedErr := errors.New("something went wrong")
	err = iter.Iterate(
		func(vertex *gograph.Vertex[string]) error {
			return expectedErr
		},
	)
	if err == nil {
		t.Error("Expect iter.Iterate(func) returns error, but got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Expect %+v error, but got %+v", expectedErr, err)
	}

	// Test the depth tracking functionality
	t.Run("DepthTracking", func(t *testing.T) {
		// Create a new graph for depth testing
		g := gograph.New[string]()

		// Create a graph with clear levels
		//    A
		//   / \
		//  B   C
		// / \ / \
		// D  E   F

		vertices := map[string]*gograph.Vertex[string]{
			"A": g.AddVertexByLabel("A"),
			"B": g.AddVertexByLabel("B"),
			"C": g.AddVertexByLabel("C"),
			"D": g.AddVertexByLabel("D"),
			"E": g.AddVertexByLabel("E"),
			"F": g.AddVertexByLabel("F"),
		}

		// Add edges that create clear depth levels
		_, _ = g.AddEdge(vertices["A"], vertices["B"]) // A -> B (depth 1)
		_, _ = g.AddEdge(vertices["A"], vertices["C"]) // A -> C (depth 1)
		_, _ = g.AddEdge(vertices["B"], vertices["D"]) // B -> D (depth 2)
		_, _ = g.AddEdge(vertices["B"], vertices["E"]) // B -> E (depth 2)
		_, _ = g.AddEdge(vertices["C"], vertices["F"]) // C -> F (depth 2)

		// Create breadth-first iterator starting from A
		iter, err := NewBreadthFirstIterator(g, "A")
		if err != nil {
			t.Fatalf("Failed to create iterator: %v", err)
		}

		// Using type assertion to access the depth functionality
		bfsIter, ok := iter.(*breadthFirstIterator[string])
		if !ok {
			t.Fatal("Failed to assert iterator as breadthFirstIterator")
		}

		// Expected vertices with their depths
		expectedDepths := map[string]int{
			"A": 0, // Starting vertex has depth 0
			"B": 1, // Direct neighbors of A
			"C": 1,
			"D": 2, // Neighbors of B and C
			"E": 2,
			"F": 2,
		}

		// Track visited vertices and their depths
		visitedDepths := make(map[string]int)

		// Iterate through vertices and check their depths
		for bfsIter.HasNext() {
			vertex := bfsIter.Next()
			depth := bfsIter.GetCurrentDepth()
			visitedDepths[vertex.Label()] = depth

			// Check if depth matches expected
			expectedDepth, exists := expectedDepths[vertex.Label()]
			if !exists {
				t.Errorf("Unexpected vertex: %s", vertex.Label())
			} else if depth != expectedDepth {
				t.Errorf("Expected depth %d for vertex %s, but got %d",
					expectedDepth, vertex.Label(), depth)
			}
		}

		// Make sure we visited all vertices with correct depths
		if len(visitedDepths) != len(expectedDepths) {
			t.Errorf("Expected to visit %d vertices, but visited %d",
				len(expectedDepths), len(visitedDepths))
		}

		// Test GetDepthOfVertex method
		bfsIter.Reset()

		// Let the iterator run completely to build depth information
		for bfsIter.HasNext() {
			bfsIter.Next()
		}

		// Check depths of each vertex using GetDepthOfVertex
		for label, expectedDepth := range expectedDepths {
			actualDepth := bfsIter.GetDepthOfVertex(label)
			if actualDepth != expectedDepth {
				t.Errorf("GetDepthOfVertex: Expected depth %d for vertex %s, but got %d",
					expectedDepth, label, actualDepth)
			}
		}

		// Check non-existent vertex
		nonExistentDepth := bfsIter.GetDepthOfVertex("Z")
		if nonExistentDepth != -1 {
			t.Errorf("Expected depth -1 for non-existent vertex, but got %d", nonExistentDepth)
		}
	})

	// Test the IterateWithDepth method
	t.Run("IterateWithDepth", func(t *testing.T) {
		// Create a simple graph
		g := gograph.New[string]()

		// Create a linear path A -> B -> C -> D
		vertices := map[string]*gograph.Vertex[string]{
			"A": g.AddVertexByLabel("A"),
			"B": g.AddVertexByLabel("B"),
			"C": g.AddVertexByLabel("C"),
			"D": g.AddVertexByLabel("D"),
		}

		// Add edges
		_, _ = g.AddEdge(vertices["A"], vertices["B"]) // A -> B
		_, _ = g.AddEdge(vertices["B"], vertices["C"]) // B -> C
		_, _ = g.AddEdge(vertices["C"], vertices["D"]) // C -> D

		// Create iterator and type assert
		iter, err := NewBreadthFirstIterator(g, "A")
		if err != nil {
			t.Fatalf("Failed to create iterator: %v", err)
		}

		bfsIter, ok := iter.(*breadthFirstIterator[string])
		if !ok {
			t.Fatal("Failed to assert iterator as breadthFirstIterator")
		}

		// Expected vertices with their depths (linear path)
		expectedVisits := []struct {
			label string
			depth int
		}{
			{"A", 0},
			{"B", 1},
			{"C", 2},
			{"D", 3},
		}

		// Track the visits
		var actualVisits []struct {
			label string
			depth int
		}

		// Use IterateWithDepth to visit vertices
		err = bfsIter.IterateWithDepth(func(v *gograph.Vertex[string], depth int) error {
			actualVisits = append(actualVisits, struct {
				label string
				depth int
			}{v.Label(), depth})
			return nil
		})

		if err != nil {
			t.Errorf("IterateWithDepth returned error: %v", err)
		}

		// Check if visits match expected
		if len(actualVisits) != len(expectedVisits) {
			t.Errorf("Expected %d visits, got %d", len(expectedVisits), len(actualVisits))
		} else {
			for i, visit := range expectedVisits {
				if actualVisits[i].label != visit.label || actualVisits[i].depth != visit.depth {
					t.Errorf("Visit #%d: Expected (%s, %d), got (%s, %d)",
						i, visit.label, visit.depth, actualVisits[i].label, actualVisits[i].depth)
				}
			}
		}

		// Test error propagation
		bfsIter.Reset()
		expectedErr := errors.New("intentional test error")

		err = bfsIter.IterateWithDepth(func(v *gograph.Vertex[string], depth int) error {
			return expectedErr
		})

		if err == nil {
			t.Error("Expected IterateWithDepth to propagate error, but got nil")
		}
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}
