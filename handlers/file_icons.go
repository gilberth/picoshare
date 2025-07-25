package handlers

import (
	"path/filepath"
	"strings"

	"github.com/mtlynch/picoshare/v2/picoshare"
)

// FileIcon represents an icon for a file type
type FileIcon struct {
	Icon      string
	Color     string
	Extension string
}

// GetFileIcon returns the appropriate icon information for a file
func GetFileIcon(filename picoshare.Filename, contentType picoshare.ContentType) FileIcon {
	ext := strings.ToLower(filepath.Ext(filename.String()))
	if ext != "" && ext[0] == '.' {
		ext = ext[1:] // Remove the dot
	}

	ct := strings.ToLower(contentType.String())

	// Check content type first, then extension
	switch {
	// Images
	case strings.Contains(ct, "image/"):
		return FileIcon{Icon: "fas fa-image", Color: "var(--file-image)", Extension: ext}

	// PDFs
	case strings.Contains(ct, "pdf") || ext == "pdf":
		return FileIcon{Icon: "fas fa-file-pdf", Color: "var(--file-pdf)", Extension: "pdf"}

	// Microsoft Office Documents
	case strings.Contains(ct, "msword") || 
		 strings.Contains(ct, "officedocument.wordprocessingml") || 
		 ext == "doc" || ext == "docx":
		return FileIcon{Icon: "fas fa-file-word", Color: "var(--file-doc)", Extension: "doc"}

	case strings.Contains(ct, "ms-excel") || 
		 strings.Contains(ct, "officedocument.spreadsheetml") || 
		 ext == "xls" || ext == "xlsx":
		return FileIcon{Icon: "fas fa-file-excel", Color: "var(--file-excel)", Extension: "xls"}

	case strings.Contains(ct, "ms-powerpoint") || 
		 strings.Contains(ct, "officedocument.presentationml") || 
		 ext == "ppt" || ext == "pptx":
		return FileIcon{Icon: "fas fa-file-powerpoint", Color: "var(--file-powerpoint)", Extension: "ppt"}

	// Archives
	case strings.Contains(ct, "zip") || strings.Contains(ct, "x-zip") || ext == "zip":
		return FileIcon{Icon: "fas fa-file-archive", Color: "var(--file-zip)", Extension: "zip"}
	case ext == "rar":
		return FileIcon{Icon: "fas fa-file-archive", Color: "var(--file-zip)", Extension: "rar"}
	case ext == "7z":
		return FileIcon{Icon: "fas fa-file-archive", Color: "var(--file-zip)", Extension: "7z"}
	case strings.Contains(ct, "gzip") || ext == "gz" || ext == "tar":
		return FileIcon{Icon: "fas fa-file-archive", Color: "var(--file-zip)", Extension: ext}

	// Videos
	case strings.Contains(ct, "video/"):
		return FileIcon{Icon: "fas fa-file-video", Color: "var(--file-video)", Extension: ext}

	// Audio
	case strings.Contains(ct, "audio/"):
		return FileIcon{Icon: "fas fa-file-audio", Color: "var(--file-audio)", Extension: ext}

	// Code files
	case ext == "js" || ext == "ts" || ext == "jsx" || ext == "tsx":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: ext}
	case ext == "html" || ext == "htm":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "html"}
	case ext == "css" || ext == "scss" || ext == "sass":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "css"}
	case ext == "py":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "py"}
	case ext == "java":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "java"}
	case ext == "php":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "php"}
	case ext == "rb":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "rb"}
	case ext == "go":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "go"}
	case ext == "c" || ext == "cpp" || ext == "h" || ext == "hpp":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: ext}
	case ext == "cs":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "cs"}

	// Text files
	case strings.Contains(ct, "text/") || ext == "txt":
		return FileIcon{Icon: "fas fa-file-alt", Color: "var(--file-text)", Extension: "txt"}
	case ext == "md" || ext == "markdown":
		return FileIcon{Icon: "fas fa-file-alt", Color: "var(--file-text)", Extension: "md"}
	case ext == "json":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "json"}
	case ext == "xml":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "xml"}
	case ext == "yaml" || ext == "yml":
		return FileIcon{Icon: "fas fa-file-code", Color: "var(--file-code)", Extension: "yml"}

	// Executables
	case ext == "exe":
		return FileIcon{Icon: "fas fa-cog", Color: "var(--file-executable)", Extension: "exe"}
	case ext == "msi":
		return FileIcon{Icon: "fas fa-cog", Color: "var(--file-executable)", Extension: "msi"}
	case ext == "app":
		return FileIcon{Icon: "fas fa-cog", Color: "var(--file-executable)", Extension: "app"}
	case ext == "dmg":
		return FileIcon{Icon: "fas fa-compact-disc", Color: "var(--file-disk)", Extension: "dmg"}
	case ext == "iso":
		return FileIcon{Icon: "fas fa-compact-disc", Color: "var(--file-disk)", Extension: "iso"}

	// Default
	default:
		return FileIcon{Icon: "fas fa-file", Color: "var(--file-generic)", Extension: ext}
	}
}

// IsImageFile checks if a file is an image based on content type
func IsImageFile(contentType picoshare.ContentType) bool {
	ct := strings.ToLower(contentType.String())
	return strings.Contains(ct, "image/jpeg") ||
		strings.Contains(ct, "image/jpg") ||
		strings.Contains(ct, "image/png") ||
		strings.Contains(ct, "image/webp") ||
		strings.Contains(ct, "image/gif") ||
		strings.Contains(ct, "image/svg")
}