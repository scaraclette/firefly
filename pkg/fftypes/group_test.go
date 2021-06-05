// Copyright © 2021 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fftypes

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupValidation(t *testing.T) {

	group := &Group{
		Namespace: "!wrong",
	}
	assert.Regexp(t, "FF10131.*namespace", group.Validate(context.Background(), false))

	group = &Group{
		Namespace:   "ok",
		Description: string(make([]byte, 4097)),
	}
	assert.Regexp(t, "FF10188.*description", group.Validate(context.Background(), false))

	group = &Group{
		Namespace:   "ok",
		Description: "ok",
	}
	assert.Regexp(t, "FF10219.*member", group.Validate(context.Background(), false))

	group = &Group{
		Namespace:   "ok",
		Description: "ok",
		Members: Members{
			{Node: NewUUID()},
		},
	}
	assert.Regexp(t, "FF10220.*member", group.Validate(context.Background(), false))

	group = &Group{
		Namespace:   "ok",
		Description: "ok",
		Members: Members{
			{Org: NewUUID()},
		},
	}
	assert.Regexp(t, "FF10221.*member", group.Validate(context.Background(), false))

	group = &Group{
		Namespace:   "ok",
		Description: "ok",
		Members: Members{
			{Org: NewUUID(), Node: NewUUID()},
		},
	}
	assert.NoError(t, group.Validate(context.Background(), false))

	assert.Regexp(t, "FF10203", group.Validate(context.Background(), true))

	group = &Group{
		Namespace:   "ok",
		Description: "ok",
		Members: Members{
			{ /* blank */ },
		},
	}
	assert.Regexp(t, "FF10220", group.Validate(context.Background(), false))

	orgID := MustParseUUID("a53fb8cc-ab22-4eed-ad21-eb104eb85902")
	nodeID := MustParseUUID("8b5c0d39-925f-4579-9c60-54f3e846ab99")
	group = &Group{
		Namespace:   "ok",
		Description: "ok",
		Members: Members{
			{Node: nodeID, Org: orgID},
			{Node: nodeID, Org: orgID},
		},
	}
	assert.Regexp(t, "FF10222", group.Validate(context.Background(), false))

	group.Members = Members{
		{Node: nodeID, Org: orgID},
	}
	b, _ := json.Marshal(&group.Members)
	assert.Equal(t, `[{"org":"a53fb8cc-ab22-4eed-ad21-eb104eb85902","node":"8b5c0d39-925f-4579-9c60-54f3e846ab99"}]`, string(b))
	group.Seal()
	assert.Equal(t, "9445512b117f19a45ef650e80e173d0ff6ab30769d9f4ace4190b44e110c24bf", group.Hash.String())

	var def Definition = group
	assert.Equal(t, fmt.Sprintf("ff-grp-%s", group.ID), def.Topic())
	def.SetBroadcastMessage(NewUUID())
	assert.NotNil(t, group.Message)
}
