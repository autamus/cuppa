//
// Copyright 2016-2021 Bryan T. Meyers <root@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package results

import (
	"fmt"
	"sort"
)

// ResultSet is a collection of the Results of a Provider query
type ResultSet struct {
	results []*Result
	query   string
}

// NewResultSet creates as empty ResultSet for the provided query
func NewResultSet(query string) *ResultSet {
	return &ResultSet{make([]*Result, 0), query}
}

var skipWords = []string{
	"master", "Master", "MASTER",
	"rc", "RC",
	"alpha", "Alpha", "ALPHA",
	"beta", "Beta", "BETA",
	"dev", "DEV",
	"unstable", "Unstable", "UNSTABLE",
	"eap", "EAP",
	"donotuse",
}

// AddResult appends a new Result
func (rs *ResultSet) AddResult(r *Result) {
	if r == nil {
		return
	}
	for _, part := range r.Version {
		for _, skip := range skipWords {
			if part == skip {
				return
			}
		}
	}
	rs.results = append(rs.results, r)
}

// First retrieves the first result from a query
func (rs *ResultSet) First() *Result {
	sort.Sort(rs)
	return rs.results[0]
}

// Last retrieves the first result from a query
func (rs *ResultSet) Last() *Result {
	if rs.Len() == 0 {
		return nil
	}
	sort.Sort(rs)
	return rs.results[len(rs.results)-1]
}

// PrintAll pretty-prints an entire ResultSet
func (rs *ResultSet) PrintAll() {
	fmt.Printf("%s: '%s'\n", "Results of Query", rs.query)
	fmt.Printf("%s: %d\n\n", "Total Number of Results", rs.Len())
	sort.Sort(rs)
	for _, r := range rs.results {
		r.Print()
	}
}

// Len is the number of elements in the ResultSet (sort.Interface)
func (rs *ResultSet) Len() int {
	return len(rs.results)
}

// Less reports whether the element with
// index i should sort before the element with index j. (sort.Interface)
func (rs *ResultSet) Less(i, j int) bool {
	if !rs.results[i].Published.IsZero() && !rs.results[j].Published.IsZero() {
		if rs.results[i].Published.Before(rs.results[j].Published) {
			return true
		}
		if rs.results[i].Published.After(rs.results[j].Published) {
			return false
		}
	}
	return rs.results[j].Version.Less(rs.results[i].Version)
}

// Swap swaps the elements with indexes i and j.
func (rs *ResultSet) Swap(i, j int) {
	rs.results[i], rs.results[j] = rs.results[j], rs.results[i]
}
