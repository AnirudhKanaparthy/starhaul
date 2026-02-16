package matrix

type SymmetricMatrix[T any] struct {
	data []T
}

func MakeSymmWithData[T any](data []T) SymmetricMatrix[T] {
	return SymmetricMatrix[T]{data}
}

func MakeSymmWithOrder[T any](order int) SymmetricMatrix[T] {
	size := (order*order + order) / 2
	return SymmetricMatrix[T]{make([]T, size)}
}

func (mat *SymmetricMatrix[T]) Order() int {
	// TODO: Optimize this in the future
	d := len(mat.data)
	td := 2 * d
	for n := 1; n < d; n += 1 {
		if td-n == n*n {
			return n
		}
	}
	// Unreachable
	return -1
}

func (mat *SymmetricMatrix[T]) Index(i int, j int) int {
	if i > j {
		temp := i
		i = j
		j = temp
	}
	n := mat.Order()
	idx := (n * i) - ((i - 1) * i / 2) + (j - i)
	return idx
}

func (mat *SymmetricMatrix[T]) At(i int, j int) T {
	idx := mat.Index(i, j)
	return mat.data[idx]
}

func (mat *SymmetricMatrix[T]) Set(i int, j int, v T) {
	idx := mat.Index(i, j)
	mat.data[idx] = v
}
