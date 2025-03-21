// zlc compresses or decompresses using zlib format.
package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

const fileExt = ".zl"
const discard = "∅"

type action struct {
	fileIn        string
	fileOut       string
	compress      bool
	compressLevel int
	force         bool

	help func()
}

var defaultAction = action{
	compress:      true,
	compressLevel: 6,
}

func openIn(path string) (*os.File, error) {
	if path == "-" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

func openOut(path string, force bool) (io.Writer, error) {
	if path == discard {
		return io.Discard, nil
	}
	if path == "-" {
		return os.Stdout, nil
	}
	var wflag int
	if force {
		wflag = os.O_TRUNC
	} else {
		wflag = os.O_EXCL
	}
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|wflag, 0644)
}

func run(a action) (rerr error) {
	in, err := openIn(a.fileIn)
	if err != nil {
		return err
	}
	defer in.Close()

	switch {
	case a.compress:
		out, err := openOut(a.fileOut, a.force)
		if err != nil {
			return err
		}
		w, err := zlib.NewWriterLevel(out, a.compressLevel)
		if err != nil {
			return fmt.Errorf("failed creating compress writer: %w", err)
		}
		defer safeCloseWriter(out, &rerr)

		_, err = io.Copy(w, in)
		if err != nil {
			return fmt.Errorf("compress: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("compress closing: %w", err)
		}

	case !a.compress:
		r, err := zlib.NewReader(in)
		if err != nil {
			return fmt.Errorf("failed creating decompress reader: %w", err)
		}
		defer safeClose(r, &rerr)

		out, err := openOut(a.fileOut, a.force)
		if err != nil {
			return err
		}
		defer safeCloseWriter(out, &rerr)

		_, err = io.Copy(out, r)
		if err != nil {
			return fmt.Errorf("decompress: %w", err)
		}
	}

	return nil
}

func main() {
	conf, err := parseArgs(os.Args[1:])
	if err != nil {
		die(2, err)
	}
	if conf.help != nil {
		conf.help()
		os.Exit(0)
	}

	err = run(conf)
	if err != nil {
		die(1, err)
	}
}

func safeClose(c io.Closer, errp *error) {
	cerr := c.Close()
	if cerr != nil && *errp == nil {
		*errp = cerr
	}
}

func safeCloseWriter(w io.Writer, errp *error) {
	if c, ok := w.(io.Closer); ok {
		safeClose(c, errp)
	}
}

func die(code int, err error) {
	fmt.Fprintln(os.Stderr, "zlc:", err)
	os.Exit(code)
}
