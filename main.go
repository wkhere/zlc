// xflate compresses or decompresses deflate stream.
package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

type config struct {
	compress      bool
	compressLevel int

	help func()
}

func run(conf config) error {
	switch {
	case conf.compress:
		w, err := zlib.NewWriterLevel(os.Stdout, conf.compressLevel)
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

	case !conf.compress:
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

func die(exitcode int, a ...any) {
	if len(a) > 0 {
		fmt.Fprint(os.Stderr, "zlc: ")
		fmt.Fprintln(os.Stderr, a...)
	}
	os.Exit(exitcode)
}
