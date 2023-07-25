package paging

type Pagination interface {
	FromIndex() int64
	PSize() int64
}

// Paginate 分页对象
type Paginate struct {
	// 第几页
	P int64 `json:"p" binding:"required"`
	// 每页条数
	Size int64 `json:"size" binding:"required"`
	// 总数
	Count int64 `json:"count" binding:"required"`
}

func (p *Paginate) FromIndex() int64 {
	if p.P < 1 {
		p.P = 1
	}
	pageCount := p.calcPageCount()
	if p.P > pageCount && p.Count > p.Size {
		p.P = pageCount
	}
	return (p.P - 1) * p.Size
}

func (p *Paginate) PSize() int64 {
	return p.Size
}

func (p *Paginate) calcPageCount() int64 {
	s := p.Count / p.Size
	m := p.Count % p.Size
	if m != 0 {
		s++
	}
	if p.Count == 0 {
		return 1
	} else {
		return s
	}
}

func NewPaginate(p int64, size int64, count int64) *Paginate {
	o := &Paginate{P: p, Size: size, Count: count}
	o.calcPageCount()
	o.FromIndex()
	return o
}
