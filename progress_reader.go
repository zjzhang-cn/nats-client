package nats_client

import "io"

type ProgressReader struct {
	io.Reader
	Total      int64 // 总大小
	ReadBytes  int64 // 已读取字节数
	OnProgress func(readBytes int64, total int64)
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.ReadBytes += int64(n)
	if pr.OnProgress != nil {
		pr.OnProgress(pr.ReadBytes, pr.Total)
	}
	return n, err
}
