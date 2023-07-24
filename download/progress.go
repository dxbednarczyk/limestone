package download

import (
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// see https://github.com/charmbracelet/bubbletea/blob/master/examples/progress-download

type progressWriter struct {
	total       int
	downloaded  int
	destination *os.File
	reader      io.Reader
	onProgress  func(float64)
}

type progressErr struct {
	err error
}

func (pw *progressWriter) Start(p *tea.Program) {
	// TeeReader calls pw.Write() each time a new response is received
	_, err := io.Copy(pw.destination, io.TeeReader(pw.reader, pw))
	if err != nil {
		p.Send(progressErr{err})
	}
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)

	if pw.total > 0 && pw.onProgress != nil {
		pw.onProgress(float64(pw.downloaded) / float64(pw.total))
	}

	return len(p), nil
}
