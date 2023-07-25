// The MIT License (MIT)
//
// Copyright (c) 2015
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// source repo: https://github.com/meirf/gopart
// modified by ayakurayuki
// license: MIT License

package part

// IndexRange specifies a single range. Low and High
// are the indexes in the larger collection at which this
// range begins and ends, respectively. Note that High
// is exclusive, whereas Low is inclusive.
type IndexRange struct {
	Low, High int
}

// Partition enables type-agnostic partitioning
// of anything indexable by specifying the length and
// the desired partition size of the indexable object.
// Consecutive index ranges are sent to the channel,
// each of which is the same size. The final range may
// be smaller than the others.
//
// For example, a collection with length 8 and
// partition size 3 yields ranges:
// {0, 3}, {3, 6}, {6, 8}
//
// This method should be used in a for...range loop.
// No results will be returned if the partition size is
// nonpositive. If the partition size is greater than the
// collection length, the range returned includes the
// entire collection.
func Partition(collectionLen, partitionSize int) chan IndexRange {
	c := make(chan IndexRange)
	if partitionSize <= 0 {
		close(c)
		return c
	}

	go func() {
		numFullPartitions := collectionLen / partitionSize
		var i int
		for ; i < numFullPartitions; i++ {
			c <- IndexRange{Low: i * partitionSize, High: (i + 1) * partitionSize}
		}

		if collectionLen%partitionSize != 0 { // left over
			c <- IndexRange{Low: i * partitionSize, High: collectionLen}
		}

		close(c)
	}()
	return c
}
