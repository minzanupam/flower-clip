package routing

import (
	"log"
	"net/http"
	"time"
)

func (s *Service) uploadSvgHandler(w http.ResponseWriter, r *http.Request) {
	// max memory for parse multipart form is 50 MB
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		log.Println(err)
		w.Write([]byte("maximum limit of 50 MB execeeded"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, fileHeader := range r.MultipartForm.File["svg-files"] {
		uploadedFile, err := fileHeader.Open()
		buf := make([]byte, 10<<20)
		bytecount, err := uploadedFile.Read(buf)
		if err != nil {
			log.Println(err)
		}
		if bytecount == 10<<20 {
			log.Println("buffer for reading uploaded file is full")
		}
		createdAt := time.Now().Format(time.RFC3339)
		res, err := s.DB.Exec(`INSERT INTO svgs (name, file, created_at) VALUES (?, ?, ?)`, fileHeader.Filename, buf, createdAt)
		if err != nil {
			log.Println(err)
			w.Write([]byte("server error"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fileID, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
		}
		log.Println(fileID)
		uploadedFile.Close()
	}
}
