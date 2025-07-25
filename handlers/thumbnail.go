package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"golang.org/x/image/webp"

	"github.com/mtlynch/picoshare/v2/picoshare"
	"github.com/mtlynch/picoshare/v2/store"
)

const (
	defaultThumbnailSize = 64
	maxThumbnailSize     = 256
)

func (s Server) thumbnailGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseEntryID(mux.Vars(r)["id"])
		if err != nil {
			log.Printf("error parsing ID for thumbnail: %v", err)
			http.Error(w, fmt.Sprintf("bad entry ID: %v", err), http.StatusBadRequest)
			return
		}

		// Parse size parameter
		size := defaultThumbnailSize
		if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
			if parsedSize, err := strconv.Atoi(sizeParam); err == nil && parsedSize > 0 && parsedSize <= maxThumbnailSize {
				size = parsedSize
			}
		}

		entry, err := s.getDB(r).GetEntryMetadata(id)
		if _, ok := err.(store.EntryNotFoundError); ok {
			http.Error(w, "entry not found", http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("error retrieving entry metadata for thumbnail %v: %v", id, err)
			http.Error(w, "failed to retrieve entry", http.StatusInternalServerError)
			return
		}

		// Check if it's an image
		if !isImageContentType(entry.ContentType) {
			http.Error(w, "not an image", http.StatusBadRequest)
			return
		}

		// Generate thumbnail
		thumbnail, err := s.generateThumbnail(r, id.String(), entry.ContentType, size)
		if err != nil {
			log.Printf("error generating thumbnail for %v: %v", id, err)
			http.Error(w, "failed to generate thumbnail", http.StatusInternalServerError)
			return
		}

		// Set headers
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour cache
		
		// Write thumbnail
		if _, err := w.Write(thumbnail); err != nil {
			log.Printf("error writing thumbnail response: %v", err)
		}
	}
}

func (s Server) generateThumbnail(r *http.Request, id string, contentType picoshare.ContentType, size int) ([]byte, error) {
	// Read the original file  
	entryID, err := parseEntryID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry ID: %w", err)
	}
	
	entryFile, err := s.getDB(r).ReadEntryFile(entryID)
	if err != nil {
		return nil, fmt.Errorf("failed to read entry file: %w", err)
	}

	// Decode image based on content type
	var img image.Image
	switch {
	case strings.Contains(contentType.String(), "jpeg") || strings.Contains(contentType.String(), "jpg"):
		img, err = jpeg.Decode(entryFile)
	case strings.Contains(contentType.String(), "png"):
		img, err = png.Decode(entryFile)
	case strings.Contains(contentType.String(), "webp"):
		img, err = webp.Decode(entryFile)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", contentType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize image to thumbnail size
	thumbnail := resize.Thumbnail(uint(size), uint(size), img, resize.Lanczos3)

	// Encode as JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 80}); err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return buf.Bytes(), nil
}

func isImageContentType(contentType picoshare.ContentType) bool {
	ct := strings.ToLower(contentType.String())
	return strings.Contains(ct, "image/jpeg") ||
		strings.Contains(ct, "image/jpg") ||
		strings.Contains(ct, "image/png") ||
		strings.Contains(ct, "image/webp")
}