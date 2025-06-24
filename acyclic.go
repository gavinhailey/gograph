package gograph

import (
	"sort"
)

// TopologySort performs a topological sort of the graph using
// Kahn's algorithm. If the sorted list of vertices does not contain
// all vertices in the graph, it means there is a cycle in the graph.
//
// It returns error if it finds a cycle in the graph.
func TopologySort[T comparable](g Graph[T]) ([]*Vertex[T], error) {
	// Initialize a map to store the inDegree of each vertex
	inDegrees := make(map[*Vertex[T]]int)
	vertices := g.GetAllVertices()
	for _, v := range vertices {
		inDegrees[v] = v.inDegree
	}

	// Initialize a queue with vertices of inDegrees zero
	queue := make([]*Vertex[T], 0)
	for v, inDegree := range inDegrees {
		if inDegree == 0 {
			queue = append(queue, v)
		}
	}

	// Initialize the sorted list of vertices
	sortedVertices := make([]*Vertex[T], 0)

	// Loop through the vertices with inDegree zero
	for len(queue) > 0 {
		// Get the next vertex with inDegree zero
		curr := queue[0]
		queue = queue[1:]

		// Add the vertex to the sorted list
		sortedVertices = append(sortedVertices, curr)

		// Decrement the inDegree of each of the vertex's neighbors
		for _, neighbor := range curr.neighbors {
			inDegrees[neighbor]--
			if inDegrees[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// If the sorted list does not contain all vertices, there is a cycle
	if len(sortedVertices) != len(vertices) {
		return nil, ErrDAGHasCycle
	}

	return sortedVertices, nil
}

// StableTopologySort does the same as TopologySort, but it takes a function
// for comparing tied vertices. This is useful when you want to
// have a stable sort order for vertices with multiple topological orderings.
func StableTopologySort[T comparable](g Graph[T], cmp func(a, b T) bool) ([]*Vertex[T], error) {
	// Initialize a map to store the inDegree of each vertex
	inDegrees := make(map[*Vertex[T]]int)
	vertices := g.GetAllVertices()
	for _, v := range vertices {
		inDegrees[v] = v.inDegree
	}

	// Initialize the sorted list of vertices
	sortedVertices := make([]*Vertex[T], 0, len(vertices))

	// Collect vertices with inDegree zero
	var zeroInDegree []*Vertex[T]
	for v, inDegree := range inDegrees {
		if inDegree == 0 {
			zeroInDegree = append(zeroInDegree, v)
		}
	}
	sortVerticesWithCmp(zeroInDegree, cmp)

	// Process vertices in the sorted order
	for len(zeroInDegree) > 0 {
		// Get the next vertex with inDegree zero
		curr := zeroInDegree[0]
		zeroInDegree = zeroInDegree[1:]

		// Add the vertex to the sorted list
		sortedVertices = append(sortedVertices, curr)

		// Collect neighbors whose in-degree becomes zero after removing current vertex
		var newZeroInDegree []*Vertex[T]
		for _, neighbor := range curr.neighbors {
			inDegrees[neighbor]--
			if inDegrees[neighbor] == 0 {
				newZeroInDegree = append(newZeroInDegree, neighbor)
			}
		}

		zeroInDegree = append(zeroInDegree, newZeroInDegree...)
		sortVerticesWithCmp(zeroInDegree, cmp)
	}

	// If the sorted list does not contain all vertices, there is a cycle
	if len(sortedVertices) != len(vertices) {
		return nil, ErrDAGHasCycle
	}

	return sortedVertices, nil
}

func sortVerticesWithCmp[T comparable](vertices []*Vertex[T], cmp func(a, b T) bool) {
	sort.Slice(vertices, func(i, j int) bool {
		return cmp(vertices[i].label, vertices[j].label)
	})
}
