package main

import "strings"

func validateFolder(folder string) bool {
	allowed := map[string]bool{
		"inbox":       true,
		"concepts":    true,
		"practices":   true,
		"ui-ux":       true,
		"postmortems": true,
		"references":  true,
	}
	return allowed[folder]
}

func sanitizeTitle(title string) string {
	title = strings.TrimSpace(title)
	if len(title) > 255 {
		return title[:255]
	}
	return title
}

func sanitizeContent(content string) string {
	if htmlPolicy == nil {
		return content
	}
	return htmlPolicy.Sanitize(content)
}

func isAllowedImageMime(ext, mime string) bool {
	mime = strings.ToLower(strings.TrimSpace(strings.Split(mime, ";")[0]))
	allowedByExt := map[string]string{
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".webp": "image/webp",
	}
	expected, ok := allowedByExt[ext]
	if !ok {
		return false
	}
	return mime == expected
}
