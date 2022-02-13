package apex

type Logger interface {
	Error()
	Info()
	Warn()
	Critical()
	Debug()
}
