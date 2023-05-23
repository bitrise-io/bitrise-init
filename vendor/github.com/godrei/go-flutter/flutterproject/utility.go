package flutterproject

import (
	"io"
	"os"
)

type Opener interface {
	Open(name string) (*os.File, error)
}

type FileOpener interface {
	OpenFile(pth string) (io.Reader, error)
}

type fileOpener struct {
	opener Opener
}

func NewFileOpener(opener Opener) FileOpener {
	return fileOpener{
		opener: opener,
	}
}

func (o fileOpener) OpenFile(pth string) (io.Reader, error) {
	f, err := o.opener.Open(pth)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return f, nil
}
