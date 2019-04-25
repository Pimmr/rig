package validators

// A Repeatable validator should return an error if the value provided is not considered valid, nil otherwise.
// This validator is used on individual values of a rig.Repeatable
type Repeatable func(interface{}) error
