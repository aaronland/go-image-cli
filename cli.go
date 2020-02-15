package cli

import (
	"bytes"
	"bufio"
	"context"
	"image"
	"github.com/aaronland/go-image-decode"
	"github.com/aaronland/go-image-encode"
	"github.com/natefinch/atomic"
	"os"
	"fmt"
)

type ProcessFunc func(context.Context, image.Image, string) (image.Image, string, error)

func Process(ctx context.Context, cb ProcessFunc, paths ...string) error {

	dec, err := decode.ImageDecoder(ctx, "image://")

	if err != nil {
		return err
	}
	
	for _, path := range paths {

		// START OF something like go-image-reader
		
		fh, err := os.Open(path)

		if err != nil {
			return err
		}

		defer fh.Close()

		im, format, err := dec.Decode(ctx, fh)

		if err != nil {
			return err
		}

		// END OF something like go-image-reader

		im, new_path, err := cb(ctx, im, path)
		
		if err != nil {
			return err
		}

		// START OF something like go-image-writer
		
		enc_uri := fmt.Sprintf("%s://", format)
		enc, err := encode.NewEncoder(ctx, enc_uri)

		if err != nil {
			return err
		}

		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)

		err = enc.Encode(ctx, im, wr)

		if err != nil {
			return err
		}

		wr.Flush()

		br := bytes.NewReader(buf.Bytes())

		err = atomic.WriteFile(new_path, br)

		if err != nil {
			return err
		}

		// END OF something like go-image-writer
	}

	return nil
}
