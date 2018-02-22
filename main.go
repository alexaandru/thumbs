package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

var (
	root  = filepath.Join("content", "images", "articole")
	sizes = []uint{257, 566}
)

func files(folder string, ch chan string) error {
	ext := ".jpg"
	f, err := os.Open(folder)
	if err != nil {
		return err
	}
	_ = f.Close()

	_ = filepath.Walk(folder, func(path string, f os.FileInfo, err error) (e error) {
		if !f.IsDir() && (ext == "" || strings.HasSuffix(path, ext)) {
			ch <- path
		}
		return
	})

	return nil
}

func isSmall(path string) bool {
	return strings.HasSuffix(path, "_small0.jpg") || strings.HasSuffix(path, "_small1.jpg")
}

func rootOf(path string) string {
	rem := 4
	if isSmall(path) {
		rem = 11
	}

	return path[:len(path)-rem]
}

func mkThumbs(path string) error {
	for i, size := range sizes {
		if err := mkThumb(path, i, size); err != nil {
			return err
		}
	}

	return nil
}

func mkThumb(path string, index int, size uint) error {
	thumb := fmt.Sprintf("%s_small%d.jpg", path, index)

	f, err := os.Open(path + ".jpg")
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		return err
	}

	m := resize.Resize(size, 0, img, resize.Lanczos2)

	out, err := os.Create(thumb)
	if err != nil {
		return err
	}

	if err := jpeg.Encode(out, m, &jpeg.Options{Quality: 80}); err != nil {
		return err
	}

	return out.Close()
}

func main() {
	ch, done := make(chan string), make(chan bool)
	thumbs := map[string]byte{}

	go func(ch chan string) {
		for path := range ch {
			r := rootOf(path)
			_, ok := thumbs[r]

			if !isSmall(path) {
				if !ok {
					thumbs[r] = 0
				}
			} else {
				if !ok {
					thumbs[r] = 1
				} else {
					thumbs[r]++
				}
			}
		}

		for k, v := range thumbs {
			if v == 2 {
				delete(thumbs, k)
			}
		}

		done <- true
	}(ch)

	if err := files(root, ch); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	close(ch)
	<-done

	for k := range thumbs {
		fmt.Printf("%s thumbnails... ", k)
		if err := mkThumbs(k); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
		}
	}
}
