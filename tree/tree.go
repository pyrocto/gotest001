package tree

// TODO: try defining ConsList in its own package, importing it here, defining ListTree here and importing it in test.go

/*

sym_sum ConsList[X, Y any] {
	type Cons struct {
		DataX X
		DataY Y
		Other int
		Next ConsList[X, Y]
	}
	type End struct{}
}

*/

type ConsList[X, Y any] interface {
	isConsList(X, Y)
}

func ConsList_Transform[X_in, Y_in any, X_out, Y_out any](f func(X_in, Y_in) (X_out, Y_out)) func(ConsList[X_in, Y_in]) ConsList[X_out, Y_out] {
	var consList_Transform func(ConsList[X_in, Y_in]) ConsList[X_out, Y_out]
	consList_Transform = func(c ConsList[X_in, Y_in]) ConsList[X_out, Y_out] {
		switch x := c.(type) {
		case Cons[X_in, Y_in]:
			x, y := f(c.DataX, c.DataY)
			return Cons[X_out, Y_out]{x, y, c.Other, consList_Transform(x.Next)}
		case End[X_in, Y_in]:
			return End[X_out, Y_out]{}
		}
		return nil
	}
	return consList_Transform
}

func ConsList_Reduce[X_in, Y_in any, Out any](consFunc func(X_in, Y_in, Out) Out, endFunc func() Out, nilFunc func() Out) func(ConsList[X_in, Y_in]) Out {
	var consList_Reduce func(ConsList[X_in, Y_in]) Out
	consList_Reduce = func(c ConsList[X_in, Y_in]) Out {
		switch x := c.(type) {
		case Cons[X_in, Y_in]:
			return consFunc(x.DataX, x.DataY, consList_Reduce(x.Next))
		case End[X_in, Y_in]:
			return endFunc()
		}
		return nilFunc()
	}
	return consList_Reduce
}

type Cons[X, Y any, Z any] struct {
	DataX X
	DataY Y
	DataZ Z
	Next  ConsList[X, Y]
}

func (Cons[X, Y, Z]) isConsList(_ X, _ Y) {}

type End[X, Y any, Z any] struct {
	DataZ Z
}

func (End[X, Y, Z]) isConsList(_ X, _ Y) {}

/*

sym_sum ListTree[X any] {
    type Node struct {
        Data ConsList[X, int]
        Left ListTree[X]
        Right ListTree[X]
    }
    type Leaf struct {}
}

*/

type ListTree[X any] interface {
	isListTree(X)
}

func ListTree_Transform[X_in, X_out any](f func(X_in) X_out) func(ListTree[X_in]) ListTree[X_out] {
	consList_Transform := ConsList_Transform(f)
	var listTree_Transform func(ListTree[X_in]) ListTree[X_out]
	listTree_Transform = func(l ListTree[X_in]) ListTree[X_out] {
		switch x := l.(type) {
		case Node[X_in]:
			return Node[X_out]{consList_Transform(l.Data), listTree_Transform(l.Left), listTree_Transform(l.Right)}
		case Leaf[X_in]:
			return Leaf[X_out]{}
		}
		return nil
	}
	return listTree_Transform
}

func ListTree_Reduce[X_in, Out any](consFunc func(X_in, int, Out) Out, endFunc func() Out, nodeFunc func(Out, Out, Out) Out, leafFunc func() Out, nilFunc func() Out) func(ListTree[X_in]) Out {
	consList_Reduce := ConsList_Reduce[X_in, Out](consFunc, endFunc, nilFunc)
	var listTree_Reduce func(ListTree[X_in]) Out
	listTree_Reduce = func(l ListTree[X_in]) Out {
		switch x := l.(type) {
		case Node[X_in]:
			return nodeFunc(consList_Reduce[X_in, Out](x.Data), listTree_Reduce(x.Left), listTree_Reduce(x.Right))
		case Leaf[X_in]:
			return leafFunc()
		}
		return nilFunc()
	}
	return listTree_Reduce
}

type Node[X any] struct {
	Data  X
	Left  ListTree[X]
	Right ListTree[X]
}

func (Node[X]) isListTree(X) {}

type Leaf[X any] struct{}

func (Leaf[X]) isListTree(X) {}
