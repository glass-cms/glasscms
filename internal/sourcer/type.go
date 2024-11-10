package sourcer

type SourceType int

const (
	SourceTypeUnspecified SourceType = iota
	SourceTypeFilesystem
)

var (
	SourceTypeValue = map[string]SourceType{
		"unspecified": SourceTypeUnspecified,
		"filesystem":  SourceTypeFilesystem,
	}

	SourceTypeString = map[SourceType]string{
		SourceTypeUnspecified: "unspecified",
		SourceTypeFilesystem:  "filesystem",
	}
)
