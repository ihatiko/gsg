package config

type ColumnSettings struct {
	Range      *Range
	Unique     *bool
	Annotation string
	Dictionary string
}
type Column struct {
	Name     string
	Settings *ColumnSettings
}
