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
// Author: Peter Mattis (peter@cockroachlabs.com)

package parser

import "testing"

func TestNormalizeExpr(t *testing.T) {
	testData := []struct {
		expr     string
		expected string
	}{
		{`(a)`, `a`},
		{`((((a))))`, `a`},
		{`ROW(a)`, `(a)`},
		{`a BETWEEN b AND c`, `a >= b AND a <= c`},
		{`a NOT BETWEEN b AND c`, `a < b OR a > c`},
		{`1+1`, `2`},
		{`(1+1,2+2,3+3)`, `(2, 4, 6)`},
		{`a+(1+1)`, `a + 2`},
		{`1+1+a`, `2 + a`},
		{`a=1+1`, `a = 2`},
		{`a=1+(2*3)-4`, `a = 3`},
		{`true OR a`, `true`},
		{`false OR a`, `a`},
		{`NULL OR a`, `NULL OR a`},
		{`a OR true`, `true`},
		{`a OR false`, `a`},
		{`a OR NULL`, `a OR NULL`},
		{`true AND a`, `a`},
		{`false AND a`, `false`},
		{`NULL AND a`, `NULL AND a`},
		{`a AND true`, `a`},
		{`a AND false`, `false`},
		{`a AND NULL`, `a AND NULL`},
		{`1 IN (1, 2, 3)`, `true`},
		{`1 IN (1, 2, a)`, `1 IN (1, 2, a)`},
		{`a<1`, `a < 1`},
		{`1>a`, `a < 1`},
		{`(a+1)=2`, `a = 1`},
		{`(a-1)>=2`, `a >= 3`},
		{`(1+a)<=2`, `a <= 1`},
		{`(1-a)>2`, `a < -1`},
		{`2<(a+1)`, `a > 1`},
		{`2>(a-1)`, `a < 3`},
		{`2<(1+a)`, `a > 1`},
		{`2>(1-a)`, `a > -1`},
		{`(a+(1+1))=2`, `a = 0`},
		{`((a+1)+1)=2`, `a = 0`},
		{`a+1+1=2`, `a = 0`},
		{`1+1>=(b+c)`, `b + c <= 2`},
		{`b+c<=1+1`, `b + c <= 2`},
		{`a/2=1`, `a = 2`},
		{`1=a/2`, `a = 2`},
		{`a=lower('FOO')`, `a = 'foo'`},
		{`lower(a)='foo'`, `lower(a) = 'foo'`},
	}
	for _, d := range testData {
		q, err := ParseTraditional("SELECT " + d.expr)
		if err != nil {
			t.Fatalf("%s: %v", d.expr, err)
		}
		expr := q[0].(*Select).Exprs[0].Expr
		r, err := NormalizeExpr(expr)
		if err != nil {
			t.Fatalf("%s: %v", d.expr, err)
		}
		if s := r.String(); d.expected != s {
			t.Errorf("%s: expected %s, but found %s", d.expr, d.expected, s)
		}
	}
}
