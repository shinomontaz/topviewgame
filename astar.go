package main

type node struct {
	Parent *node
	Pos    *Position
	g      int
	h      int
	f      int
}

func newNode(p *node, pos *Position) *node {
	return &node{
		Parent: p,
		Pos:    pos,
	}
}

func (n *node) isEqual(o *node) bool {
	return n.Pos.X == o.Pos.X && n.Pos.Y == o.Pos.Y
}

type AStar struct {
}

func (as AStar) GetPath(l Level, start, end *Position) []Position {
	openList := []*node{}
	closedList := []*node{}
	closedmap := make(map[*node]struct{})
	openmap := make(map[*node]struct{})

	startNode := newNode(nil, start)
	endNode := newNode(nil, end)

	lh, lw := l.GetDimensions()

	openList = append(openList, startNode)
	openmap[startNode] = struct{}{}

	for {
		if len(openList) == 0 {
			break
		}
		currIdx := 0
		currNode := openList[currIdx]

		for idx, item := range openList {
			if item.f < currNode.f {
				currNode = item
				currIdx = idx
			}
		}

		openList = append(openList[:currIdx], openList[currIdx+1:]...)
		delete(openmap, currNode)

		closedList = append(closedList, currNode)

		if currNode.isEqual(endNode) {
			path := []Position{}
			curr := currNode
			for {
				if curr == nil {
					break
				}
				path = append(path, *curr.Pos)
				curr = curr.Parent
			}

			reversedPath := make([]Position, len(path))
			for i := len(path) - 1; i >= 0; i-- {
				reversedPath[len(path)-1-i] = path[i]
			}
		}

		edges := make([]*node, 0)

		if currNode.Pos.Y > 0 {
			tile := l.Tiles[l.GetIndexFromXY(currNode.Pos.X, currNode.Pos.Y-1)]
			if tile.TileType != WALL {
				upNodePos := Position{
					X: currNode.Pos.X,
					Y: currNode.Pos.Y - 1,
				}
				newNode := newNode(currNode, &upNodePos)
				edges = append(edges, newNode)
			}
		}

		if currNode.Pos.Y < lh {
			tile := l.Tiles[l.GetIndexFromXY(currNode.Pos.X, currNode.Pos.Y+1)]
			if tile.TileType != WALL {
				downNodePos := Position{
					X: currNode.Pos.X,
					Y: currNode.Pos.Y + 1,
				}
				newNode := newNode(currNode, &downNodePos)
				edges = append(edges, newNode)
			}
		}

		if currNode.Pos.X > 0 {
			tile := l.Tiles[l.GetIndexFromXY(currNode.Pos.X-1, currNode.Pos.Y)]
			if tile.TileType != WALL {
				leftNodePos := Position{
					X: currNode.Pos.X - 1,
					Y: currNode.Pos.Y,
				}
				newNode := newNode(currNode, &leftNodePos)
				edges = append(edges, newNode)
			}
		}
		if currNode.Pos.X < lw {
			tile := l.Tiles[l.GetIndexFromXY(currNode.Pos.X+1, currNode.Pos.Y)]
			if tile.TileType != WALL {
				rightNodePos := Position{
					X: currNode.Pos.X + 1,
					Y: currNode.Pos.Y,
				}
				newNode := newNode(currNode, &rightNodePos)
				edges = append(edges, newNode)
			}
		}

		for _, edge := range edges {
			if _, ok := closedmap[edge]; ok {
				continue
			}
			edge.g = currNode.g + 1
			edge.h = edge.Pos.GetManhattanDistance(endNode.Pos)
			edge.f = edge.g + edge.h

			if _, ok := openmap[edge]; ok {
				isFurther := false
				for _, n := range openList {
					if edge.g > n.g {
						isFurther = true

						break
					}
				}

				if isFurther {
					continue
				}
			}

			openList = append(openList, edge)
		}
	}

	return nil
}
