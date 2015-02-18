// Author: Bodie Solomon
//         bodie@synapsegarden.net
//         github.com/binary132
//
//         2015-02-16

package gojsonschema

// This returns a *Schema suitable for stubbing in defaults test, which only
// cares about the pool document.
func MakeTestingSchema(doc interface{}) *Schema {
	var testingPool *schemaPool
	if doc != nil {
		testingPool = &schemaPool{standaloneDocument: doc}
	}
	return &Schema{pool: testingPool}
}

// This just uses the internal definition.
func (s *Schema) GetDocProperties() map[string]interface{} {
	return s.getDocProperties()
}

// Same
func IterateAndInsert(into map[string]interface{}, properties map[string]interface{}) {
	iterateAndInsert(into, properties)
}
