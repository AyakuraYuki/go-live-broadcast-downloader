package osstool

type AliImageInfo struct {
	FileSize struct {
		Value string `json:"value"`
	} `json:"FileSize"`
	Format struct {
		Value string `json:"value"`
	} `json:"Format"`
	ImageHeight struct {
		Value string `json:"value"`
	} `json:"ImageHeight"`
	ImageWidth struct {
		Value string `json:"value"`
	} `json:"ImageWidth"`
	ResolutionUnit struct {
		Value string `json:"value"`
	} `json:"ResolutionUnit"`
	XResolution struct {
		Value string `json:"value"`
	} `json:"XResolution"`
	YResolution struct {
		Value string `json:"value"`
	} `json:"YResolution"`
}
