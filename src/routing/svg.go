package routing

import (
	"log"
	"net/http"
	"time"
)

func (s *Service) uploadSvgHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r, s.Store)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("please login before continuing"))
		return
	}
	// max memory for parse multipart form is 50 MB
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("maximum limit of 50 MB execeeded"))
		return
	}
	for _, fileHeader := range r.MultipartForm.File["svg-files"] {
		uploadedFile, err := fileHeader.Open()
		// buffer assumed to be 1 MB
		buf := make([]byte, 1<<20)
		bytecount, err := uploadedFile.Read(buf)
		if err != nil {
			log.Println(err)
		}
		if bytecount == 1<<20 {
			log.Println("buffer for reading uploaded file is full")
		}
		createdAt := time.Now().Format(time.RFC3339)
		res, err := s.DB.Exec(`INSERT INTO svgs (name, file,
			created_at, user_id) VALUES (?, ?, ?, ?)`,
			fileHeader.Filename, string(buf[0:bytecount+1]), createdAt, userID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("server error"))
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
