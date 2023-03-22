package directusapi

import (
	"fmt"
	"strings"
)

type query struct {
	eqFilter    map[string]string
	nEqFilter   map[string]string
	inFilter    map[string]string
	nNullFilter []string
	sort        []string
	limit       *int
	offset      *int
	searchStr   *string
	// relational objects query where key must be a dot separated path
	deepQuery deepQuery
}

type deepQuery struct {
	eqFilter    map[string]string
	nEqFilter   map[string]string
	inFilter    map[string]string
	nNullFilter []string
	limit       *keyVal[string, int]
	offset      *keyVal[string, int]
}

type keyVal[K, V any] struct {
	key K
	val V
}

func None() query {
	return query{
		map[string]string{},
		map[string]string{},
		map[string]string{},
		[]string{},
		[]string{},
		nil,
		nil,
		nil,
		deepQuery{},
	}
}

func (q query) Nnull(k string) query {
	q.nNullFilter = append(q.nNullFilter, k)
	return q
}

func Nnull(k string) query {
	return None().Nnull(k)
}

func (q query) Neq(k, v string) query {
	q.nEqFilter[k] = v
	return q
}

func Neq(k, v string) query {
	return None().Eq(k, v)
}

func (q query) Eq(k, v string) query {
	q.eqFilter[k] = v
	return q
}

func Eq(k, v string) query {
	return None().Eq(k, v)
}

func (q query) In(k, v string) query {
	q.inFilter[k] = v
	return q
}

func In(k, v string) query {
	return None().In(k, v)
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

func (q query) DeepEq(k, v string) query {
	q.deepQuery.eqFilter[k] = v
	return q
}

func (q query) DeepLimit(k string, limit int) query {
	q.deepQuery.limit = &keyVal[string, int]{k, limit}
	return q
}

func (q query) DeepOffset(k string, offset int) query {
	q.deepQuery.offset = &keyVal[string, int]{k, offset}
	return q
}

func (q query) asKeyValue(v Version) map[string]string {
	if v == V8 {
		return q.asKeyValueV8()
	}
	return q.asKeyValueV9()
}

func (q query) asKeyValueV8() map[string]string {
	out := map[string]string{}
	for k, v := range q.eqFilter {
		out[fmt.Sprintf("filter[%s][eq]", k)] = v
	}
	for k, v := range q.nEqFilter {
		out[fmt.Sprintf("filter[%s][neq]", k)] = v
	}
	for k, v := range q.inFilter {
		out[fmt.Sprintf("filter[%s][in]", k)] = v
	}
	for _, v := range q.nNullFilter {
		out[fmt.Sprintf("filter[%s][nnull]", v)] = ""
	}
	if len(q.sort) > 0 {
		out["sort"] = strings.Join(q.sort, ",")
	}
	if q.limit != nil {
		out["limit"] = fmt.Sprint(*q.limit)
	}
	if q.offset != nil {
		out["offset"] = fmt.Sprint(*q.offset)
	}
	return out
}

func (q query) asKeyValueV9() map[string]string {
	out := map[string]string{
		"limit": "-1",
	}
	for k, v := range q.eqFilter {
		out[fmt.Sprintf("filter[%s][_eq]", k)] = v
	}
	for k, v := range q.nEqFilter {
		out[fmt.Sprintf("filter[%s][_neq]", k)] = v
	}
	for k, v := range q.inFilter {
		out[fmt.Sprintf("filter[%s][_in]", k)] = v
	}
	for _, v := range q.nNullFilter {
		out[fmt.Sprintf("filter[%s][_nnull]", v)] = "true"
	}
	if len(q.sort) > 0 {
		out["sort"] = strings.Join(q.sort, ",")
	}
	if q.limit != nil {
		out["limit"] = fmt.Sprint(*q.limit)
	}
	if q.offset != nil {
		out["offset"] = fmt.Sprint(*q.offset)
	}
	q.parseDeepQuery(out)
	return out
}

func (q query) parseDeepQuery(out map[string]string) {
	parsePath := func(path string) string {
		split := strings.Split(path, ".")
		paramPath := ""
		for _, p := range split {
			paramPath += "[" + p + "]"
		}
		return paramPath
	}
	for k, v := range q.deepQuery.eqFilter {
		out[fmt.Sprintf("deep%s[_eq]", parsePath(k))] = v
	}
	for k, v := range q.deepQuery.nEqFilter {
		out[fmt.Sprintf("deep%s[_neq]", parsePath(k))] = v
	}
	for k, v := range q.deepQuery.inFilter {
		out[fmt.Sprintf("deep%s[_in]", parsePath(k))] = v
	}
	for _, v := range q.deepQuery.nNullFilter {
		out[fmt.Sprintf("deep%s[_nnull]", parsePath(v))] = "true"
	}
	if q.deepQuery.limit != nil {
		out[fmt.Sprintf("deep%s[_limit]", parsePath(q.deepQuery.limit.key))] = fmt.Sprint(q.deepQuery.limit.val)
	}
	if q.deepQuery.offset != nil {
		out[fmt.Sprintf("deep%s[_offset]", parsePath(q.deepQuery.offset.key))] = fmt.Sprint(q.deepQuery.offset.val)
	}
}
