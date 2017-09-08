package arrgo

import (
	"fmt"
	"strings"
	"math"
)

type Arrf struct {
	shape   []int
	strides []int
	data    []float64
}

//通过[]float64，形状来创建多维数组。
//输入参数1：data []float64，以·C· 顺序存储，作为多维数组的输入数据，内部复制一份新的internalData，不改变data。
//输入参数2：shape ...int，指定多维数组的形状，多维，类似numpy中的shape。
//	如果某一个（仅支持一个维度）维度为负数，则根据len(data)推断该维度的大小。
//情况1：如果不指定shape，而且data为nil，则创建一个空的*Arrf。
//情况2：如果不指定shape，而且data不为nil，则创建一个len(data)大小的一维*Arrf。
//情况3：如果指定shape，而且data不为nil，则根据data大小创建多维数组，如果len(data)不等于shape，或者len(data)不能整除shape，抛出异常。
//情况4：如果指定shape，而且data为nil，则创建shape大小的全为0.0的多维数组。
func Array(data []float64, shape ...int) *Arrf {
	if len(shape) == 0 && data == nil {
		return &Arrf{
			shape:   []int{0},
			strides: []int{0, 1},
			data:    []float64{},
		}
	}

	if len(shape) == 0 && data != nil {
		internalData := make([]float64, len(data)) //复制data，不影响输入的值。
		copy(internalData, data)
		return &Arrf{
			shape:   []int{len(data)},
			strides: []int{len(data), 1},
			data:    internalData,
		}
	}

	if data == nil {
		for _, v := range shape {
			if v <= 0 {
				fmt.Println("shape should be positive when data is nill")
				panic(SHAPE_ERROR)
			}
		}
		length := ProductIntSlice(shape)
		internalShape := make([]int, len(shape))
		copy(internalShape, shape)
		strides := make([]int, len(shape)+1)
		strides[len(shape)] = 1
		for i := len(shape) - 1; i >= 0; i-- {
			strides[i] = strides[i+1] * internalShape[i]
		}

		return &Arrf{
			shape:   internalShape,
			strides: strides,
			data:    make([]float64, length),
		}
	}

	var dataLength = len(data)
	negativeIndex := -1
	internalShape := make([]int, len(shape))
	copy(internalShape, shape)
	for k, v := range shape {
		if v < 0 {
			if negativeIndex < 0 {
				negativeIndex = k
				internalShape[k] = 1
			} else {
				fmt.Println("shape can only have one negative demention.")
				panic(SHAPE_ERROR)
			}
		}
	}
	shapeLength := ProductIntSlice(internalShape)

	if dataLength < shapeLength {
		fmt.Println("data length is shorter than shape length.")
		panic(SHAPE_ERROR)
	}
	if (dataLength % shapeLength) != 0 {
		fmt.Println("data length cannot divided by shape length")
		panic(SHAPE_ERROR)
	}

	if negativeIndex >= 0 {
		internalShape[negativeIndex] = dataLength / shapeLength
	}

	strides := make([]int, len(internalShape)+1)
	strides[len(internalShape)] = 1
	for i := len(internalShape) - 1; i >= 0; i-- {
		strides[i] = strides[i+1] * internalShape[i]
	}

	internalData := make([]float64, len(data))
	copy(internalData, data)

	return &Arrf{
		shape:   internalShape,
		strides: strides,
		data:    internalData,
	}
}

// 通过指定起始、终止和步进量来创建一维Array。
// 输入参数： vals，可以有三种情况，详见下面描述。
// 情况1：Arange(stop): 以0开始的序列，创建Array [0, 0+(-)1, ..., stop)，不包括stop，stop符号决定升降序。
// 情况2：Arange(start, stop):创建Array [start, start +(-)1, ..., stop)，如果start小于start则递增，否则递减。
// 情况3：Arange(start, stop, step):创建Array [start, start + step, ..., stop)，step符号决定升降序。
// 输入参数多于三个的都会被忽略。
// 输出序列为“整型数”序列。
func Arange(vals ...int) *Arrf {
	var start, stop, step int = 0, 0, 1

	switch len(vals) {
	case 0:
		fmt.Println("range function should have range")
		panic(PARAMETER_ERROR)
	case 1:
		if vals[0] <= 0 {
			step = -1
			stop = vals[0] + 1
		} else {
			stop = vals[0] - 1
		}
	case 2:
		if vals[1] < vals[0] {
			step = -1
			stop = vals[1] + 1
		} else {
			stop = vals[1] - 1
		}
		start = vals[0]
	default:
		if vals[1] < vals[0] {
			if vals[2] >= 0 {
				fmt.Println("increment should be negative.")
				panic(PARAMETER_ERROR)
			}
			stop = vals[1] + 1
		} else {
			if vals[2] <= 0 {
				fmt.Println("increment should be positive.")
				panic(PARAMETER_ERROR)
			}
			stop = vals[1] - 1
		}
		start, step = vals[0], vals[2]
	}

	a := Array(nil, int(math.Abs(float64((stop-start)/step)))+1)
	for i, v := 0, start; i < len(a.data); i, v = i+1, v+step {
		a.data[i] = float64(v)
	}
	return a
}

//判断Arrf是否为空数组。
//如果内部的data长度为0或者为nil，返回true，否则位false。
func (a *Arrf) IsEmpty() bool {
	return len(a.data) == 0 || a.data == nil
}

//创建shape形状的多维数组，全部填充为fullvalue。
//必须指定shape，否则抛出异常。
func Full(fullValue float64, shape ...int) *Arrf {
	if len(shape) == 0 {
		fmt.Println("shape is empty!")
		panic(SHAPE_ERROR)
	}
	arr := Array(nil, shape...)
	for i := range arr.data {
		arr.data[i] = fullValue
	}

	return arr
}

//根据shape创建全为1.0的多维数组。
func Ones(shape ...int) *Arrf {
	return Full(1, shape...)
}

//根据输入的多维数组的形状创建全1的多维数组。
func OnesLike(a *Arrf) *Arrf {
	return Ones(a.shape...)
}

//根据shape创建全为0的多维数组。
func Zeros(shape ...int) *Arrf {
	return Full(0, shape...)
}

//根据输入的多维数组的形状创建全0的多维数组。
func ZerosLike(a *Arrf) *Arrf {
	return Zeros(a.shape...)
}

// String Satisfies the Stringer interface for fmt package
func (a *Arrf) String() (s string) {
	switch {
	case a == nil:
		return "<nil>"
	case a.data == nil || a.shape == nil || a.strides == nil:
		return "<nil>"
	case a.strides[0] == 0:
		return "[]"
	case len(a.shape) == 1:
		return fmt.Sprint(a.data)
	}

	stride := a.shape[len(a.shape)-1]

	for i, k := 0, 0; i+stride <= len(a.data); i, k = i+stride, k+1 {

		t := ""
		for j, v := range a.strides {
			if i%v == 0 && j < len(a.strides)-2 {
				t += "["
			}
		}

		s += strings.Repeat(" ", len(a.shape)-len(t)-1) + t
		s += fmt.Sprint(a.data[i: i+stride])

		t = ""
		for j, v := range a.strides {
			if (i+stride)%v == 0 && j < len(a.strides)-2 {
				t += "]"
			}
		}

		s += t + strings.Repeat(" ", len(a.shape)-len(t)-1)
		if i+stride != len(a.data) {
			s += "\n"
			if len(t) > 0 {
				s += "\n"
			}
		}
	}
	return
}

//获取index指定位置的元素。
//index必须在shape规定的范围内，否则会抛出异常。
//index的长度必须小于等于维度的个数，否则会抛出异常。
//如果index的个数小于维度个数，则会取后面的第一个值。
func (a *Arrf) At(index ...int) float64 {
	idx, err := a.valIndex(index...)
	if err != nil {
		panic(err)
	}
	return a.data[idx]
}

//详见At函数。
func (a *Arrf) Get(index ...int) float64 {
	return a.At(index...)
}

//At函数的内部实现，返回index指定的元素在切片中的位置，如果有错误，则返回error。
func (a *Arrf) valIndex(index ...int) (int, error) {
	idx := 0
	if len(index) > len(a.shape) {
		fmt.Println("index len should not longer than shape.")
		return -1, INDEX_ERROR
	}
	for i, v := range index {
		if v >= a.shape[i] || v < 0 {
			fmt.Println("index value out of range.")
			return -1, INDEX_ERROR
		}
		idx += v * a.strides[i+1]
	}
	return idx, nil
}

//获取多维数组元素的个数。
func (a *Arrf) Length() int {
	return len(a.data)
}

//改变原始多维数组的形状，并返回改变后的多维数组的指引引用。
//不会创建新的数据副本。
//如果新的shape的大小和原来多维数组的大小不同，则抛出异常。
func (a *Arrf) Reshape(shape ...int) *Arrf {
	if a.Length() != ProductIntSlice(shape) {
		fmt.Println("new shape length does not equal to original array length.")
		panic(SHAPE_ERROR)
	}

	internalShape := make([]int, len(shape))
	copy(internalShape, shape)
	a.shape = internalShape

	a.strides = make([]int, len(a.shape)+1)
	a.strides[len(a.shape)] = 1
	for i := len(a.shape) - 1; i >= 0; i-- {
		a.strides[i] = a.strides[i+1] * a.shape[i]
	}

	return a
}

//创建一个n X n 的2维单位矩阵(数组)。
func Eye(n int) *Arrf {
	arr := Zeros(n, n)
	for i := 0; i < n; i++ {
		arr.Set(1, i, i)
	}
	return arr
}

//Eye的另一种称呼，详见Eye函数。
func Identity(n int) *Arrf {
	return Eye(n)
}

func (a *Arrf) Set(val float64, index ...int) *Arrf {
	idx, _ := a.valIndex(index...)
	a.data[idx] = val
	return a
}

func (a *Arrf) Values() []float64 {
	return a.data
}

//Return evenly spaced numbers over ta specified interval.
//
//Returns `num` evenly spaced samples, calculated over the
//interval [`start`, `stop`].
//
//The endpoint of the interval can optionally be excluded.
//
//Parameters
//----------
//start : scalar
//The starting value of the sequence.
//stop : scalar
//The end value of the sequence, unless `endpoint` is set to False.
//In that case, the sequence consists of all but the last of ``num + 1``
//evenly spaced samples, so that `stop` is excluded.  Note that the step
//size changes when `endpoint` is False.
//num : int, optional
//Number of samples to generate. Default is 50. Must be non-negative.
//endpoint : bool, optional
//If True, `stop` is the last sample. Otherwise, it is not included.
//Default is True.
//retstep : bool, optional
//If True, return (`samples`, `step`), where `step` is the spacing
//between samples.
//dtype : dtype, optional
//The type of the output array.  If `dtype` is not given, infer the data
//type from the other input arguments.
//
//.. versionadded:: 1.9.0
//
//Returns
//-------
//samples : ndarray
//There are `num` equally spaced samples in the closed interval
//``[start, stop]`` or the half-open interval ``[start, stop)``
//(depending on whether `endpoint` is True or False).
func Linspace(start, stop, num int) *Arrf {
	var data = make([]float64, num)
	var startF, stopF = float64(start), float64(stop)
	if startF <= stopF {
		var step = (stopF - startF) / (float64(num - 1.0))
		for i := range data {
			data[i] = startF + float64(i)*step
		}
		return Array(data, num)
	} else {
		var step = (startF - stopF) / (float64(num - 1.0))
		for i := range data {
			data[i] = startF - float64(i)*step
		}
		return Array(data, num)
	}
}

func (a *Arrf) Copy() *Arrf {
	b := ZerosLike(a)
	copy(b.data, a.data)
	return b
}

func (a *Arrf) Ndims() int {
	return len(a.shape)
}

//Returns ta view of the array with axes transposed.
//
//For ta 1-D array, this has no effect. (To change between column and
//row vectors, first cast the 1-D array into ta matrix object.)
//For ta 2-D array, this is the usual matrix transpose.
//For an n-D array, if axes are given, their order indicates how the
//axes are permuted (see Examples). If axes are not provided and
//``ta.shape = (i[0], i[1], ... i[n-2], i[n-1])``, then
//``ta.transpose().shape = (i[n-1], i[n-2], ... i[1], i[0])``.
//
//Parameters
//----------
//axes : None, tuple of ints, or `n` ints
//
//* None or no argument: reverses the order of the axes.
//
//* tuple of ints: `i` in the `j`-th place in the tuple means `ta`'s
//`i`-th axis becomes `ta.transpose()`'s `j`-th axis.
//
//* `n` ints: same as an n-tuple of the same ints (this form is
//intended simply as ta "convenience" alternative to the tuple form)
//
//Returns
//-------
//out : ndarray
//View of `ta`, with axes suitably permuted.
func (a *Arrf) Transpose(axes ...int) *Arrf {
	var n = a.Ndims()
	var permutation []int
	var nShape []int

	switch len(axes) {
	case 0:
		permutation = make([]int, n)
		for i := 0; i < n; i++ {
			permutation[i] = n - 1 - i
			nShape[i] = a.shape[permutation[i]]
		}

	case n:
		permutation = axes
		nShape = make([]int, n)
		for i := range nShape {
			nShape[i] = a.shape[permutation[i]]
		}

	default:
		panic(DIMENTION_ERROR)
	}

	var totalIndexSize = 1
	for i := range a.shape {
		totalIndexSize *= a.shape[i]
	}

	var indexsSrc = make([][]int, totalIndexSize)
	var indexsDst = make([][]int, totalIndexSize)

	var b = Zeros(nShape...)
	var index = make([]int, n)
	for i := 0; i < totalIndexSize; i++ {
		tindexSrc := make([]int, n)
		copy(tindexSrc, index)
		indexsSrc[i] = tindexSrc
		var tindexDst = make([]int, n)
		for j := range tindexDst {
			tindexDst[j] = index[permutation[j]]
		}
		indexsDst[i] = tindexDst

		var j = n - 1
		index[j]++
		for {
			if j > 0 && index[j] >= a.shape[j] {
				index[j-1]++
				index[j] = 0
				j--
			} else {
				break
			}
		}
	}
	for i := range indexsSrc {
		b.Set(a.Get(indexsSrc[i]...), indexsDst[i]...)
	}
	return b
}

func (a *Arrf) Count(axis ...int) int {
	if len(axis) == 0 {
		return a.strides[0]
	}

	var cnt = 1
	for _, w := range axis {
		cnt *= a.shape[w]
	}
	return cnt
}

func (a *Arrf) Flatten() *Arrf {
	ra := make([]float64, len(a.data))
	copy(ra, a.data)
	return Array(ra, len(a.data))
}
