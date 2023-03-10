package main

import (
	"errors"
)

func RemoveRepByMap(slc []string) []string {
	result := []string{}         //存放返回的不重复切片
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 //当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		//如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e) //当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

type unionSet struct {
	rank []int // 以数值sz[i]为根的树的深度(高度)
	set  []int
}

func NewUnionSet(size int) *unionSet {
	buf1 := make([]int, size)
	for i := 0; i < size; i++ {
		buf1[i] = i
	}
	buf2 := make([]int, size)
	for i := 0; i < size; i++ {
		buf2[i] = 1
	}

	return &unionSet{
		rank: buf2,
		set:  buf1,
	}
}

func (set *unionSet) GetSize() int {
	return len(set.set)
}

func (set *unionSet) GetID(p int) (int, error) {
	if p < 0 || p > len(set.set) {
		return 0, errors.New(
			"failed to get ID,index is illegal.")
	}

	return set.getRoot(p), nil
}

func (set *unionSet) getRoot(p int) int {
	for p != set.set[p] {
		set.set[p] = set.set[set.set[p]]
		p = set.set[p]
	}
	return p
}

func (set *unionSet) IsConnected(p, q int) (bool, error) {
	if p < 0 || p > len(set.set) || q < 0 || q > len(set.set) {
		return false, errors.New(
			"error: index is illegal.")
	}
	return set.getRoot(set.set[p]) == set.getRoot(set.set[q]), nil
}

func (set *unionSet) Union(p, q int) error {
	if p < 0 || p > len(set.set) || q < 0 || q > len(set.set) {
		return errors.New(
			"error: index is illegal.")
	}

	pRoot := set.getRoot(p)
	qRoot := set.getRoot(q)

	if pRoot != qRoot {
		if set.rank[pRoot] < set.rank[qRoot] {
			set.set[pRoot] = qRoot
		} else if set.rank[qRoot] < set.rank[pRoot] {
			set.set[qRoot] = pRoot
		} else {
			set.set[pRoot] = qRoot
			set.rank[qRoot] += 1
		}
	}
	return nil
}

// func main11() {
// 	ufs := NewUnionSet(10)
// 	ufs.Union(0, 1)
// 	ufs.Union(1, 2)
// 	ufs.Union(3, 4)
// 	ufs.Union(1, 4)
// 	ufs.Union(5, 6)
// 	ufs.Union(5, 7)
// 	ufs.Union(8, 9)

// 	// 0 1 2 3 4      5 6 7    8 9
// 	p := 0
// 	q := 1
// 	judge, _ := ufs.IsConnected(p, q)
// 	fmt.Println("%d and %d is connected?", p, q)
// 	fmt.Println(judge)

// 	p = 2
// 	q = 5
// 	judge, _ = ufs.IsConnected(p, q)
// 	fmt.Println("%d and %d is connected?", p, q)
// 	fmt.Println(judge)

// 	q = 7
// 	q = 9
// 	judge, _ = ufs.IsConnected(p, q)
// 	fmt.Println("%d and %d is connected?", p, q)
// 	fmt.Println(judge)

// 	fmt.Println("union 0 and 7")
// 	ufs.Union(0, 7)
// 	p = 2
// 	q = 5

// 	judge, _ = ufs.IsConnected(p, q)
// 	fmt.Println("%d and %d is connected?", p, q)
// 	fmt.Println(judge)

// 	fmt.Println("union 0 and 7")
// 	ufs.Union(4, 8)
// 	p = 5
// 	q = 9

// 	judge, _ = ufs.IsConnected(p, q)
// 	fmt.Println("%d and %d is connected?", p, q)
// 	fmt.Println(judge)

// 	initial_string := []string{"1", "2", "1", "3"}
// 	fmt.Println(RemoveRepByMap(initial_string))
// }
