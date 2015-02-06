// Author: Bodie Solomon
//         bodie@synapsegarden.net
//         github.com/binary132
//
//         2015-02-16

package gojsonschema_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	gjs "github.com/xeipuuv/gojsonschema"
)

// Get a *Schema given a properties map.
func schemaFromProperties(t *testing.T, properties map[string]interface{}) *gjs.Schema {
	schemaMap := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}
	loader := gjs.NewGoLoader(schemaMap)
	schema, err := gjs.NewSchema(loader)
	assert.Nil(t, err)
	return schema
}

func TestGetDocProperties(t *testing.T) {
	for i, test := range []struct {
		should         string
		usingSchema    interface{}
		expectedErr    string
		expectedResult interface{}
	}{{
		should:      "fail for nil pool",
		expectedErr: "document pool for schema not set",
	}, {
		should:      "fail for non-map document",
		usingSchema: "not-a-map",
		expectedErr: "schema document is not a map[string]interface{}",
	}, {
		should:      "fail for non schema document",
		usingSchema: map[string]interface{}{"foo": "bar"},
		expectedErr: "schema document has no properties",
	}, {
		should: "work for an OK document",
		usingSchema: map[string]interface{}{"properties": map[string]interface{}{
			"foo": "bar"}},
		expectedResult: map[string]interface{}{"foo": "bar"},
	}} {
		fmt.Printf("TestGetDocProperties test (%d) | should %s\n", i, test.should)
		schema := gjs.MakeTestingSchema(test.usingSchema)
		doc, err := schema.GetDocProperties()
		if test.expectedErr != "" {
			assert.EqualError(t, err, test.expectedErr)
			continue
		}
		assert.Equal(t, test.expectedResult, doc)
	}

	println()
}

func TestIterateAndInsert(t *testing.T) {
	for i, test := range []struct {
		should          string
		usingProperties map[string]interface{}
		into            map[string]interface{}
		expectedResult  map[string]interface{}
	}{{
		should: "insert value as expected for a simple default",
		usingProperties: map[string]interface{}{
			"num": map[string]interface{}{
				"default": 5,
				"type":    "integer",
			}},
		into:           map[string]interface{}{},
		expectedResult: map[string]interface{}{"num": 5},
	}, {
		should: "not overwrite existing values",
		usingProperties: map[string]interface{}{
			"num": map[string]interface{}{
				"default": 5,
				"type":    "integer",
			}},
		into:           map[string]interface{}{"num": 8},
		expectedResult: map[string]interface{}{"num": 8},
	}, {
		should: "create a simple map with a default inner value",
		usingProperties: map[string]interface{}{
			"num": map[string]interface{}{
				"properties": map[string]interface{}{
					"dum": map[string]interface{}{
						"default": 5,
					}}}},
		into: map[string]interface{}{},
		expectedResult: map[string]interface{}{
			"num": map[string]interface{}{
				"dum": 5,
			}},
	}, {
		should: "non-destructively insert a value into an inner map",
		usingProperties: map[string]interface{}{
			"num": map[string]interface{}{
				"properties": map[string]interface{}{
					"dum": map[string]interface{}{
						"default": 5,
					}}}},
		into: map[string]interface{}{
			"num": map[string]interface{}{
				"gum": 8,
			}},
		expectedResult: map[string]interface{}{
			"num": map[string]interface{}{
				"gum": 8,
				"dum": 5,
			}},
	}, {
		should: "non-destructively insert a value into an inner map",
		usingProperties: map[string]interface{}{
			"num": map[string]interface{}{
				"properties": map[string]interface{}{
					"dum": map[string]interface{}{
						"default": 5,
					}}}},
		into: map[string]interface{}{
			"num": map[string]interface{}{
				"gum": 8,
			},
			"foo": "bar",
		},
		expectedResult: map[string]interface{}{
			"num": map[string]interface{}{
				"gum": 8,
				"dum": 5,
			},
			"foo": "bar",
		},
	}, {
		should: "non-destructively insert a value into an inner map",
		usingProperties: map[string]interface{}{
			"num": map[string]interface{}{
				"properties": map[string]interface{}{
					"dum": map[string]interface{}{
						"default": 5,
					},
					"foo": map[string]interface{}{
						"default": map[string]interface{}{
							"bar": "baz",
						}}}},
			"foo": map[string]interface{}{
				"properties": map[string]interface{}{
					"bar": map[string]interface{}{
						"default": "baz",
					}}}},
		into: map[string]interface{}{
			"num": map[string]interface{}{
				"gum": 8,
			},
			"foo": map[string]interface{}{},
		},
		expectedResult: map[string]interface{}{
			"num": map[string]interface{}{
				"gum": 8,
				"dum": 5,
				"foo": map[string]interface{}{"bar": "baz"},
			},
			"foo": map[string]interface{}{"bar": "baz"},
		},
	}, {
		should: "not insert a value if there is no default",
		usingProperties: map[string]interface{}{
			"num": map[string]interface{}{
				"properties": map[string]interface{}{
					"dum": map[string]interface{}{
						"type":        "string",
						"description": "something dum",
					},
					"foo": map[string]interface{}{
						"default": map[string]interface{}{
							"bar": "baz",
						}}}},
			"foo": map[string]interface{}{
				"properties": map[string]interface{}{
					"bar": map[string]interface{}{
						"baz": map[string]interface{}{
							"woz": map[string]interface{}{
								"type": "integer",
							}}}}}},
		into: map[string]interface{}{},
		expectedResult: map[string]interface{}{
			"num": map[string]interface{}{
				"foo": map[string]interface{}{"bar": "baz"},
			}},
	}, {
		should: "ignore bad values",
		usingProperties: map[string]interface{}{
			"foo": map[string]interface{}{
				"properties": map[string]interface{}{
					"bar": map[string]interface{}{
						"default": 5,
					}}}},
		into:           map[string]interface{}{"foo": 5},
		expectedResult: map[string]interface{}{"foo": 5},
	}} {
		fmt.Printf("TestIterateAndInsert test (%d) | should %s\n", i, test.should)
		gjs.IterateAndInsert(test.into, test.usingProperties)
		assert.Equal(t, test.expectedResult, test.into)
	}

	println()
}

func TestDefaultObjects(t *testing.T) {
	for i, test := range []struct {
		should          string
		usingProperties map[string]interface{}
		usingParams     map[string]interface{}
		expectedResult  map[string]interface{}
		expectedError   string
	}{} {
		fmt.Printf("TestDefaultObjects test (%d) | should %s\n", i, test.should)
		schema := schemaFromProperties(t, test.usingProperties)
		err := schema.InsertDefaults(map[string]interface{}(test.usingParams))
		if test.expectedError != "" {
			assert.EqualError(t, err, test.expectedError)
			continue
		}
		assert.Nil(t, err)
		assert.Equal(t, test.expectedResult, test.usingParams)
	}

	println()
}
