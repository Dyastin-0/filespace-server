package types

type file struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type MoveBody struct {
	File       file   `json:"file"`
	TargetPath string `json:"targetPath"`
}
