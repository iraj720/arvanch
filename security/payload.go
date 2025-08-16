package security

type PayloadTransformer interface {
	Transform(text string) (string, error)
}
