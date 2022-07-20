package directusapi

type query struct {
	eqFilter  map[string]string
	sort      []string
	limit     *int
	offset    *int
	searchStr *string
}

func None() query {
	return query{
		map[string]string{},
		[]string{},
		nil,
		nil,
		nil,
	}
}

func (q query) Eq(k, v string) query {
	q.eqFilter[k] = v
	return q
}

func Eq(k, v string) query {
	return None().Eq(k, v)
}

func (q query) SortAsc(sortBy string) query {
	q.sort = append(q.sort, sortBy)
	return q
}

func SortAsc(sortBy string) query {
	return None().SortAsc(sortBy)
}

func (q query) SortDesc(sortBy string) query {
	q.sort = append(q.sort, "-"+sortBy)
	return q
}

func SortDesc(sortBy string) query {
	return None().SortDesc(sortBy)
}

func (q query) Limit(limit int) query {
	q.limit = &limit
	return q
}

func Limit(limit int) query {
	return None().Limit(limit)
}

func (q query) Offset(offset int) query {
	q.offset = &offset
	return q
}

func Offset(offset int) query {
	return None().Offset(offset)
}

func (q query) Search(str string) query {
	q.searchStr = &str
	return q
}

func Search(str string) query {
	return None().Search(str)
}
