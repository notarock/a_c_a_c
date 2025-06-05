package filters

// Filter defines an interface for message filtering.
type Filter interface {
	Filter(message string) bool
}
