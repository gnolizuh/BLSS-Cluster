package utils

var (
	HTTP_METHOD_GET = 0
	HTTP_METHOD_POST = 1
)

type Config struct {
	Listen string `json:"listen"`
	Method string `json:"method"`
}