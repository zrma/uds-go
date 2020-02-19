package uds

// File struct is file wrapper
type File struct {
	Name        string
	MD5         string
	Size        string
	EncodedSize string
	SizeNumeric string
	Parents     []string
}

// Chunk struct is split file chunk
type Chunk struct {
}
