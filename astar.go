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
	openList := Heap[*node]{}
	openMap := map[Position]*node{}
	closedMap := map[Position]struct{}{}

	startNode := newNode(nil, start)
	startNode.g = 0
	startNode.h = start.GetManhattanDistance(end)
	startNode.f = startNode.g + startNode.h

	openList.Push(startNode.f, startNode)
	openMap[*start] = startNode

	for {
		currIdx, currNode := openList.Pop()
		if currIdx == -1 {
			break
		}
		delete(openMap, *currNode.Pos)
		closedMap[*currNode.Pos] = struct{}{}

		if currNode.Pos.X == end.X && currNode.Pos.Y == end.Y {
			path := []Position{}
			for n := currNode; n != nil; n = n.Parent {
				path = append(path, *n.Pos)
			}
			// Reverse path
			for i := 0; i < len(path)/2; i++ {
				path[i], path[len(path)-1-i] = path[len(path)-1-i], path[i]
			}

			return path
		}

		// Expand neighbors
		dirs := []Position{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		for _, dir := range dirs {
			neighbor := Position{X: currNode.Pos.X + dir.X, Y: currNode.Pos.Y + dir.Y}

			if !l.InBounds(neighbor.X, neighbor.Y) || l.TileAt(neighbor.X, neighbor.Y).TileType == WALL {
				continue
			}

			if _, visited := closedMap[neighbor]; visited {
				continue
			}

			tentativeG := currNode.g + 1
			existingNode, inOpen := openMap[neighbor]

			if !inOpen || tentativeG < existingNode.g {
				newN := &node{
					Parent: currNode,
					Pos:    &Position{X: neighbor.X, Y: neighbor.Y},
					g:      tentativeG,
					h:      neighbor.GetManhattanDistance(end),
				}
				newN.f = newN.g + newN.h

				openList.Push(newN.f, newN)
				openMap[neighbor] = newN
			}
		}
	}

	return nil
}
