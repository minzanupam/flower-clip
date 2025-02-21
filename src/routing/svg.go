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

func (s *Service) deleteSvgHandler(w http.ResponseWriter, r *http.Request) {
	svgID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("could not parse svg id"))
		return
	}
	userID, err := authenticate(r, s.Store)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("login required"))
		return
	}
	_, err = s.DB.Exec(`DELETE FROM svgs WHERE id = ? AND user_id = ?`, svgID, userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to delete svg"))
		return
	}
}

func (s *Service) editPageHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r, s.Store)
	if err != nil {
		log.Println(err)
	}
	authenticated := false
	if userID != 0 {
		authenticated = true
	}
	if !authenticated {
		component := templates.EditPage(false, []shared_types.SVG{})
		component.Render(r.Context(), w)
		return
	}
	rows, err := s.DB.Query(`SELECT id, name, file, created_at FROM svgs WHERE user_id = ?`, userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
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
	component := templates.EditPage(authenticated, svgs)
	component.Render(r.Context(), w)
}

func (s *Service) downloadSvgHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse svg id"))
		return
	}
	userID, err := authenticate(r, s.Store)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("login required"))
	}
	row := s.DB.QueryRow(`SELECT file FROM svgs WHERE id = ? AND user_id = ?`, id, userID)
	var svg_file string
	if err = row.Scan(&svg_file); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("file not found"))
	}
	w.Header().Set("Content-Type", "image/svg")
	w.Write([]byte(svg_file))
}
