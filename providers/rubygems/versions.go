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

package rubygems

import (
	"github.com/DataDrake/cuppa/results"
)

// Versions holds one or more Rubygems Versions
type Versions []Version

// Convert turns a Rubygems result set into a Cuppa result set
func (crs *Versions) Convert(name string) *results.ResultSet {
	rs := results.NewResultSet(name)
	for _, rel := range *crs {
		if r := rel.Convert(name); r != nil {
			rs.AddResult(r)
		}
	}
	return rs
}
