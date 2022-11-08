package types

type Info struct {
	// Size of the body in bytes.
	ContentLength int64
	// A standard MIME type describing the format of the object data.
	ContentType string
	// Checksum of object.
	// For S3, the checksum is an `ETag`. It is calculated from MD5 of each part of the object.
	Checksum string
}

type ListOpts struct {
	// If this var is not set s3 will return S3ListMaxKeys objects
	MaxKeys int32
}
