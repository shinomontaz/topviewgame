package main

type Heap[T any] struct {
	list []int
	obj  []T
}

func (h *Heap[T]) Push(key int, o T) {
	h.list = append(h.list, key)
	h.obj = append(h.obj, o)

	h.up()
}

func (h *Heap[T]) Pop() (int, T) {
	if len(h.list) == 0 {
		var zero T
		return -1, zero
	}

	h.list[0], h.list[len(h.list)-1] = h.list[len(h.list)-1], h.list[0]
	h.obj[0], h.obj[len(h.obj)-1] = h.obj[len(h.obj)-1], h.obj[0]

	x := h.list[len(h.list)-1]
	o := h.obj[len(h.obj)-1]

	h.list = h.list[:len(h.list)-1]
	h.obj = h.obj[:len(h.obj)-1]

	h.down()

	return x, o
}

func (h *Heap[T]) up() {
	i := len(h.list) - 1
	for i > 0 {
		a := (i - 1) / 2
		if h.list[a] <= h.list[i] {
			break
		}
		h.list[a], h.list[i] = h.list[i], h.list[a]
		h.obj[a], h.obj[i] = h.obj[i], h.obj[a]

		i = a
	}
}

func (h *Heap[T]) down() {
	i := 0
	for {
		a, b := 2*i+1, 2*i+2
		j := i
		if b < len(h.list) && h.list[b] < h.list[j] {
			j = b
		}
		if a < len(h.list) && h.list[a] < h.list[j] {
			j = a
		}
		if j == i {
			break
		}

		h.list[i], h.list[j] = h.list[j], h.list[i]
		h.obj[i], h.obj[j] = h.obj[j], h.obj[i]

		i = j
	}
}
