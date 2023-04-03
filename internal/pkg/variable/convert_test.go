// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package variable

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestConvertCtyToInterface(t *testing.T) {

	// test basic type
	testCases := []struct {
		name string
		val  cty.Value
		t    reflect.Kind
	}{
		{"bool", cty.BoolVal(true), reflect.Bool},
		{"string", cty.StringVal("test"), reflect.String},
		{"number", cty.NumberIntVal(1), reflect.Int},
		{"map", cty.MapVal(map[string]cty.Value{"test": cty.BoolVal(true)}), reflect.Map},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := convertCtyToInterface(tc.val)
			require.NoError(t, err)

			resType := reflect.TypeOf(res).Kind()
			require.Equal(t, tc.t, resType)
		})
	}

	// test list of list
	t.Run("lists of lists", func(t *testing.T) {
		testListOfList := cty.ListVal([]cty.Value{
			cty.ListVal([]cty.Value{
				cty.BoolVal(true),
			}),
		})

		resListOfList, err := convertCtyToInterface(testListOfList)
		require.NoError(t, err)

		tempList, ok := resListOfList.([]interface{})
		require.True(t, ok)

		_, ok = tempList[0].([]interface{})
		require.True(t, ok)
	})

	// test list of maps
	t.Run("list of maps", func(t *testing.T) {
		testListOfMaps := cty.ListVal([]cty.Value{
			cty.MapVal(map[string]cty.Value{
				"test": cty.BoolVal(true),
			}),
		})

		resListOfMaps, err := convertCtyToInterface(testListOfMaps)
		require.NoError(t, err)

		_, ok := resListOfMaps.([]map[string]interface{})
		require.True(t, ok)
	})

	// test map of maps
	t.Run("map of maps", func(t *testing.T) {
		testMapOfMaps := cty.MapVal(map[string]cty.Value{
			"test": cty.MapVal(map[string]cty.Value{"test": cty.BoolVal(true)}),
		})

		restMapOfMaps, err := convertCtyToInterface(testMapOfMaps)
		require.NoError(t, err)

		tempMapOfMaps, ok := restMapOfMaps.(map[string]interface{})
		require.True(t, ok)

		_, ok = tempMapOfMaps["test"].(map[string]interface{})
		require.True(t, ok)
	})
}
