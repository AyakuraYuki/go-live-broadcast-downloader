package internal

const (
	MaxTick = 20000
)

type platformHandler func(*Task) error

var PlatformHandler = map[string]platformHandler{
	"asobistage": asobistage,
}
