package output

type Mode string

const (
	ModeTable Mode = "table"
	ModeJSON  Mode = "json"
)

func FromJSONFlag(enabled bool) Mode {
	if enabled {
		return ModeJSON
	}
	return ModeTable
}
