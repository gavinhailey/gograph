![build](https://github.com/gavinhailey/gograph/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/gavinhailey/gograph)](https://goreportcard.com/report/github.com/gavinhailey/gograph)
[![Go Reference](https://pkg.go.dev/badge/github.com/gavinhailey/gograph.svg)](https://pkg.go.dev/github.com/gavinhailey/gograph)

# gograph

Fork of [hmdsefi/gograph](https://github.com/hmdsefi/gograph)

<img alt="golang generic graph package" src="https://user-images.githubusercontent.com/11541936/221823924-358994d2-44ff-4236-bbc8-b404de62293e.png" style="width:40%" align="right" title="gograph"/>
<br/>
<br/>
<p>GoGraph is a lightweight, efficient, and easy-to-use graph data structure
implementation written in Go. It provides a versatile framework for representing 
graphs and performing various operations on them, making it ideal for both
educational purposes and practical applications.</p>
<br/><br/><br/>

## Table of Contents

<!-- TOC -->
* [gograph](#gograph)
  * [Table of Contents](#table-of-contents)
  * [Install](#install)
  * [How to Use](#how-to-use)
    * [Graph](#graph)
      * [Directed](#directed)
      * [Acyclic](#acyclic)
      * [Undirected](#undirected)
      * [Weighted](#weighted)
    * [Traverse](#traverse)
  * [License](#license)
<!-- TOC -->

## Install

Use `go get` command to get the latest version of the `gograph`:

```shell
go get github.com/gavinhailey/gograph
```

Then you can use import the `gograph` to your code:

```go
package main

import "github.com/gavinhailey/gograph"
```

## How to Use

### Graph

`gograph` contains the `Graph[T comparable]` interface that provides all needed APIs to
manage a graph. All the supported graph types in `gograph` library implemented this interface.

```go
type Graph[T comparable] interface {
GraphType

AddEdge(from, to *Vertex[T], options ...EdgeOptionFunc) (*Edge[T], error)
GetAllEdges(from, to *Vertex[T]) []*Edge[T]
GetEdge(from, to *Vertex[T]) *Edge[T]
EdgesOf(v *Vertex[T]) []*Edge[T]
RemoveEdges(edges ...*Edge[T])
AddVertexByLabel(label T, options ...VertexOptionFunc) *Vertex[T]
AddVertex(v *Vertex[T])
GetVertexByID(label T) *Vertex[T]
GetAllVerticesByID(label ...T) []*Vertex[T]
GetAllVertices() []*Vertex[T]
RemoveVertices(vertices ...*Vertex[T])
ContainsEdge(from, to *Vertex[T]) bool
ContainsVertex(v *Vertex[T]) bool
}
```

The generic type of the `T` in `Graph` interface represents the vertex label. The type of `T`
should be comparable. You cannot use slices and function types for `T`.

#### Directed

![directed-graph](https://user-images.githubusercontent.com/11541936/221904292-face2083-16da-491f-a339-2164b7040264.png)

```go
graph := New[int](gograph.Directed())

graph.AddEdge(gograph.NewVertex(1), gograph.NewVertex(2))
graph.AddEdge(gograph.NewVertex(1), gograph.NewVertex(3))
graph.AddEdge(gograph.NewVertex(2), gograph.NewVertex(2))
graph.AddEdge(gograph.NewVertex(3), gograph.NewVertex(4))
graph.AddEdge(gograph.NewVertex(4), gograph.NewVertex(5))
graph.AddEdge(gograph.NewVertex(5), gograph.NewVertex(6))
```

#### Acyclic

![acyclic-graph](https://user-images.githubusercontent.com/11541936/221911652-ce2dfb5f-5547-4f26-8412-94ad9124d4fa.png)

```go
graph := New[int](gograph.Directed())

graph.AddEdge(gograph.NewVertex(1), gograph.NewVertex(2))
graph.AddEdge(gograph.NewVertex(2), gograph.NewVertex(3))
_, err := graph.AddEdge(gograph.NewVertex(3), gograph.NewVertex(1))
if err != nil {
// do something
}
```

#### Undirected

![undirected-graph](https://user-images.githubusercontent.com/11541936/221908261-a009049d-2b71-46c3-9026-faa4dcc2a693.png)

```go
// by default graph is undirected
graph := New[string]()

graph.AddEdge(gograph.NewVertex("A"), gograph.NewVertex("B"))
graph.AddEdge(gograph.NewVertex("A"), gograph.NewVertex("D"))
graph.AddEdge(gograph.NewVertex("B"), gograph.NewVertex("C"))
graph.AddEdge(gograph.NewVertex("B"), gograph.NewVertex("D"))
```

#### Weighted

![weighted-edge](https://user-images.githubusercontent.com/11541936/221908269-b6db15fb-6104-49d9-b9b9-acc062d94e4a.png)

```go
graph := New[string]()

vA := gograph.AddVertexByLabel("A")
vB := gograph.AddVertexByLabel("B")
vC := gograph.AddVertexByLabel("C")
vD := gograph.AddVertexByLabel("D")

graph.AddEdge(vA, vB, gograph.WithEdgeWeight(4))
graph.AddEdge(vA, vD, gograph.WithEdgeWeight(3))
graph.AddEdge(vB, vC, gograph.WithEdgeWeight(3))
graph.AddEdge(vB, vD, gograph.WithEdgeWeight(1))
graph.AddEdge(vC, vD, gograph.WithEdgeWeight(2))
```

![weighted-vertex](https://user-images.githubusercontent.com/11541936/221908278-83f3138d-8b28-4c38-825a-627a46d65294.png)

```go
graph := New[string]()
vA := gograph.AddVertexByLabel("A", gograph.WithVertexWeight(3))
vB := gograph.AddVertexByLabel("B", gograph.WithVertexWeight(2))
vC := gograph.AddVertexByLabel("C", gograph.WithVertexWeight(4))

graph.AddEdge(vA, vB)
graph.AddEdge(vB, vC)
```

### Traverse

Traverse package provides the iterator interface that guarantees all the algorithm export the same APIs:

```go
type Iterator[T comparable] interface {
	HasNext() bool
	Next() *gograph.Vertex[T]
	Iterate(func(v *gograph.Vertex[T]) error) error
	Reset()
}
```

This package contains the following iterators:

- [Breadth-First iterator](https://github.com/gavinhailey/gograph/tree/master/traverse#BFS)
- [Depth-First iterator](https://github.com/gavinhailey/gograph/tree/master/traverse#DFS)
- [Topological iterator](https://github.com/gavinhailey/gograph/tree/master/traverse#Topological-Sort)
- [Closest-First iterator](https://github.com/gavinhailey/gograph/tree/master/traverse#Closest-First)
- [Random Walk iterator](https://github.com/gavinhailey/gograph/tree/master/traverse#random-walk)

## License

Apache License, please see [LICENSE](https://github.com/gavinhailey/gograph/blob/master/LICENSE) for details.