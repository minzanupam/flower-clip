package routing

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"app.flower.clip/src/shared_types"
	"app.flower.clip/src/templates"
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
	var idString strings.Builder
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
			continue
		}
		uploadedFile.Close()

		if idString.Len() == 0 {
			if _, err = idString.WriteString(strconv.Itoa(int(fileID))); err != nil {
				log.Println(err)
			}
		} else {
			if _, err = idString.WriteString(","); err != nil {
				log.Println(err)
			}
			if _, err = idString.WriteString(strconv.Itoa(int(fileID))); err != nil {
				log.Println(err)
			}
		}
	}

	rows, err := s.DB.Query(fmt.Sprintf(`SELECT ID, name, file, created_at FROM svgs WHERE id IN (%s)`, idString.String()))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}
	var svgs []shared_types.SVG
	for rows.Next() {
		var svg shared_types.SVG
		var createdAtString string
		if err = rows.Scan(&svg.ID, &svg.Name, &svg.File, &createdAtString); err != nil {
			log.Println(err)
			continue
		}
		svg.CreatedAt, err = time.Parse(time.RFC3339, createdAtString)
		if err != nil {
			log.Println(err)
			continue
		}
		svgs = append(svgs, svg)
	}
	component := templates.RenderSvgs(svgs)
	if err = component.Render(r.Context(), w); err != nil {
		log.Println(err)
	}
}
