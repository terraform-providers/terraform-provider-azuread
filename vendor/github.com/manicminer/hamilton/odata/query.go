package odata

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Query struct {
	// Count includes a count of the total number of items in a collection alongside the page of data values
	Count bool

	// Expand includes the expanded resource or collection referenced by a single relationship
	Expand Expand

	// Filter retrieves just a subset of a collection, or relationships like members, memberOf, transitiveMembers, and transitiveMemberOf
	Filter string

	// Format specifies the media format of the items returned
	Format Format

	// OrderBy specify the sort order of the items returned
	OrderBy OrderBy

	// Search restricts the results of a request to match a search criterion
	Search string // complicated

	// Select returns a set of properties that are different than the default set for an individual resource or a collection of resources
	Select []string

	// Skip sets the number of items to skip at the start of a collection
	Skip int

	// Top specifies the page size of the result set
	Top int
}

func (q Query) Values() url.Values {
	p := url.Values{}
	if q.Count {
		p.Add("$count", fmt.Sprintf("%t", q.Count))
	}
	if expand := q.Expand.String(); expand != "" {
		p.Add("$expand", expand)
	}
	if q.Filter != "" {
		p.Add("$filter", q.Filter)
	}
	if format := string(q.Format); format != "" {
		p.Add("$format", format)
	}
	if orderBy := q.OrderBy.String(); orderBy != "" {
		p.Add("$orderby", orderBy)
	}
	if q.Search != "" {
		p.Add("$search", fmt.Sprintf(`"%s"`, q.Search))
	}
	if len(q.Select) > 0 {
		p.Add("$select", strings.Join(q.Select, ","))
	}
	if q.Skip > 0 {
		p.Add("$skip", strconv.Itoa(q.Skip))
	}
	if q.Top > 0 {
		p.Add("$top", strconv.Itoa(q.Top))
	}
	return p
}

type Expand struct {
	Relationship string
	Select       []string
}

func (e Expand) String() (val string) {
	val = e.Relationship
	if len(e.Select) > 0 {
		val = fmt.Sprintf("%s($select=%s)", val, strings.Join(e.Select, ","))
	}
	return
}

type Format string

const (
	FormatJson Format = "json"
	FormatAtom Format = "atom"
	FormatXml  Format = "xml"
)

type Direction string

const (
	Ascending  Direction = "asc"
	Descending Direction = "desc"
)

type OrderBy struct {
	Field     string
	Direction Direction
}

func (o OrderBy) String() (val string) {
	val = o.Field
	if val != "" && o.Direction != "" {
		val = fmt.Sprintf("%s %s", val, o.Direction)
	}
	return
}
