package handler

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func Resize(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(30 << 20) // Max 30MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Unable to decode image", http.StatusInternalServerError)
		return
	}

	resizedImg := resize.Resize(1000, 0, img, resize.Lanczos3)

	fileHeader := r.MultipartForm.File["image"][0]
	outFileName := "see_albania_resized_" + filepath.Base(fileHeader.Filename)
	outFile, err := os.Create(outFileName)
	if err != nil {
		http.Error(w, "Unable to create output file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, resizedImg, nil)
	if err != nil {
		http.Error(w, "Unable to encode image", http.StatusInternalServerError)
		return
	}

	log.Printf("Image resized and saved as: %s\n", outFileName)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", outFileName))

	http.ServeFile(w, r, outFileName)
}
