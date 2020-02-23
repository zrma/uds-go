package uds

// File struct is file wrapper
type File struct {
	Name        string
	Mime        string
	Size        string
	EncodedSize string
	SizeNumeric string
	Parents     []string

	ID     string
	MD5    string
	Shared bool
}

// Init method initialize File struct's parents not to be nil
func (f *File) Init() {
	if f.Parents == nil {
		f.Parents = []string{"root"}
	}
}

// Chunk struct is split file chunk
type Chunk struct {
	Path    string
	Part    int64
	MaxSize int64
	Media   *File
	Parent  string

	RangeEnd int64
}

// Init method initialize Chunk struct's range end boundary
func (c *Chunk) Init() {
	const chunkReadLengthBytes int64 = 750000
	c.RangeEnd = (c.Part + 1) * chunkReadLengthBytes
	if c.RangeEnd > c.MaxSize {
		c.RangeEnd = c.MaxSize
	}
}
