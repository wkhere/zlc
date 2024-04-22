// zlc compresses or decompresses using zlib format.
package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

type action struct {
	compress      bool
	compressLevel int

	help func()
}

var defaultAction = action{
	compress:      true,
	compressLevel: 6,
}

func run(a action) error {
	switch {
	case a.compress:
		w, err := zlib.NewWriterLevel(os.Stdout, a.compressLevel)
		if err != nil {
			return fmt.Errorf("failed creating compress writer: %w", err)
		}
		_, err = io.Copy(w, os.Stdin)
		if err != nil {
			return fmt.Errorf("compress: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("compress closing: %w", err)
		}

	case !a.compress:
		r, err := zlib.NewReader(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed creating decompress reader: %w", err)
		}
		defer r.Close()
		_, err = io.Copy(os.Stdout, r)
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

func die(code int, err error) {
	fmt.Fprintln(os.Stderr, "zlc:", err)
	os.Exit(code)
}
