package output_port

type AuthCode interface {
	Generate4DigitCode() string
}
