// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Marc Berhault (marc@cockroachlabs.com)

package keys

import "github.com/cockroachdb/cockroach/proto"

// Span describes a span of keys: [start,end).
// This will correspond to one or more range.
type Span struct {
	Start, End proto.Key
}

var (
	// Meta1Span holds all first level addressing.
	Meta1Span = Span{proto.KeyMin, Meta2Prefix}

	// ZoneConfigSpan is the range for all zone configs.
	ZoneConfigSpan = Span{ConfigZonePrefix, ConfigZonePrefix.PrefixEnd()}

	// UserDataSpan is the non-meta and non-structured portion of the key space.
	UserDataSpan = Span{SystemMax, TableDataPrefix}

	// SystemDBSpan is the range of system objects for structured data.
	SystemDBSpan = Span{TableDataPrefix, UserTableDataMin}

	// NoSplitSpans describes the ranges that should never be split.
	// Meta1Span: needed to find other ranges.
	// ZoneConfigSpan: configs are hierarchical by prefix.
	// SystemDBSpan: system objects have interdepencies.
	NoSplitSpans = []Span{Meta1Span, ZoneConfigSpan, SystemDBSpan}
)
