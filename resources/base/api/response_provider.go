package api



type ResponseProvider interface {
	Apply(body []byte) ([]byte,error)
}