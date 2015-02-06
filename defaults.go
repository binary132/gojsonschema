// Author: Bodie Solomon
//         bodie@synapsegarden.net
//         github.com/binary132
//
//         2015-02-16

package gojsonschema

import "errors"

// InsertDefaults takes a map[string]interface{} and inserts any missing
// default values specified in the JSON-Schema document.  If the given "into"
// map is nil, it will be created.
//
// This is an initial implementation which does not support refs, etc.
// The target value must be an object and not a bare value.
//
// If called on a bad schema, this will panic.
func (s *Schema) InsertDefaults(into map[string]interface{}) error {
	properties, err := s.getDocProperties()
	if err != nil {
		return err
	}

	// Make sure the "into" map isn't nil.
	if into == nil {
		into = make(map[string]interface{})
	}

	// Now insert the default values at the keys in the "into" map,
	// non-destructively.

	iterateAndInsert(into, properties)

	return nil
}

// getDocProperties retrieves the "properties" key of the standalone doc of
// the schema.
func (s *Schema) getDocProperties() (map[string]interface{}, error) {
	// First, check for a missing document, since we will work with the
	// raw document map.
	if s.pool == nil {
		return nil, errors.New("document pool for schema not set")
	}

	f := s.pool.GetStandaloneDocument()

	// We need a map[string]interface because we have to check for a
	// default key.
	docMap, ok := f.(map[string]interface{})
	if !ok {
		return nil, errors.New("schema document is not a map[string]interface{}")
	}

	// We need a schema with properties, since this feature does not support
	// raw value schemas.
	m, ok := docMap["properties"]
	if !ok {
		return nil, errors.New("schema document has no properties")
	}

	// Panic on bad schemas
	typedMap := m.(map[string]interface{})
	return typedMap, nil
}

// iterateAndInsert takes a target map and inserts any missing default values
// as specified in the properties map, according to JSON-Schema.
func iterateAndInsert(into, properties map[string]interface{}) {
	for property, schema := range properties {
		// Iterate over each key of the schema.  Each key should have
		// a schema describing the key's properties.  If it does not,
		// ignore this key.
		typedSchema := schema.(map[string]interface{})

		if v, ok := into[property]; ok {
			// If there is already a map value in the target for
			// this key, don't overwrite it; step in.  Ignore
			// non-map values.
			if innerMap, ok := v.(map[string]interface{}); ok {
				// The schema should have an inner
				// schema since it's an object
				schemaProperties := typedSchema["properties"]
				typedProperties := schemaProperties.(map[string]interface{})
				iterateAndInsert(innerMap, typedProperties)
			}
			continue
		}

		if d, ok := typedSchema["default"]; ok {
			// Most basic case: we have a default value. Done for
			// this key.
			into[property] = d
			continue
		}

		if p, ok := typedSchema["properties"]; ok {
			// If we have a "properties" key, this is an object,
			// and we need to go deeper.
			innerSchema := p.(map[string]interface{})
			tmpTarget := make(map[string]interface{})
			iterateAndInsert(tmpTarget, innerSchema)
			if len(tmpTarget) > 0 {
				into[property] = tmpTarget
			}
		}
	}
}
