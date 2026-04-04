// Copyright 2021 The Cockroach Authors.
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
// permissions and limitations under the License.

package sessiondata

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/sessiondatapb"
	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	initialElem := &SessionData{
		SessionData: sessiondatapb.SessionData{
			ApplicationName: "bob",
		},
	}
	secondElem := &SessionData{
		SessionData: sessiondatapb.SessionData{
			ApplicationName: "jane",
		},
	}
	thirdElem := &SessionData{
		SessionData: sessiondatapb.SessionData{
			ApplicationName: "t-marts",
		},
	}
	s := NewStack(initialElem.Clone())
	require.Equal(t, s.Top(), initialElem)
	require.EqualError(t, s.Pop(), "there must always be at least one element in the SessionData stack")
	require.Equal(t, s.Elems(), []*SessionData{initialElem})

	s.Push(secondElem.Clone())
	require.Equal(t, s.Top(), secondElem.Clone())
	s.Push(thirdElem.Clone())
	require.Equal(t, s.Top(), thirdElem.Clone())
	require.Equal(t, s.Elems(), []*SessionData{initialElem, secondElem, thirdElem})

	require.NoError(t, s.Pop())
	require.Equal(t, s.Top(), secondElem.Clone())
	require.NoError(t, s.Pop())
	require.Equal(t, s.Top(), initialElem.Clone())
	require.EqualError(t, s.Pop(), "there must always be at least one element in the SessionData stack")
	require.Equal(t, s.Elems(), []*SessionData{initialElem})

	s.Push(secondElem.Clone())
	s.Push(thirdElem.Clone())
	s.PushTopClone()
	require.Equal(t, s.Elems(), []*SessionData{initialElem, secondElem, thirdElem, thirdElem})
	c := s.Clone()
	s.PopAll()
	require.Equal(t, s.Elems(), []*SessionData{initialElem})
	s.PopAll()
	require.Equal(t, s.Elems(), []*SessionData{initialElem})
	require.Equal(t, c.Elems(), []*SessionData{initialElem, secondElem, thirdElem, thirdElem})

	s.Push(secondElem.Clone())
	s.Push(thirdElem.Clone())
	copyStack := s.Clone()
	require.Error(t, s.PopN(3), "there must always be at least one element in the SessionData stack")
	require.NoError(t, s.PopN(2))
	require.Equal(t, s.Elems(), []*SessionData{initialElem})

	require.Equal(t, s.Base(), initialElem)

	s.Top().Database = "some other value"
	s.Replace(copyStack)
	require.Equal(t, s.Elems(), []*SessionData{initialElem, secondElem, thirdElem})
	require.Equal(t, s.Base(), initialElem)

	copyStack.Top().Database = "some other value"
	s.Replace(copyStack)
	copiedElem := thirdElem
	copiedElem.Database = "some other value"
	require.Equal(t, s.Elems(), []*SessionData{initialElem, secondElem, copiedElem})
	require.Equal(t, s.Base(), initialElem)
}
