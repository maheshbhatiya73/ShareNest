package controllers

import (
	"ShareNestbackend/config"
	"ShareNestbackend/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const uploadPath = "./uploads"

func hashPassword(pw string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
    return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    err := r.ParseMultipartForm(10 << 30) // 10 GB max
    if err != nil {
        http.Error(w, "Invalid form data: "+err.Error(), http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error getting file: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    password := r.FormValue("password")
    if password == "" {
        http.Error(w, "Password is required", http.StatusBadRequest)
        return
    }

    hashedPw, err := hashPassword(password)
    if err != nil {
        http.Error(w, "Error hashing password: "+err.Error(), http.StatusInternalServerError)
        return
    }

    fileID := uuid.New().String()

    // MODIFY THIS LINE: prepend UUID + "-" + original filename
    filename := fileID + "-" + handler.Filename

    filepath := filepath.Join(uploadPath, filename)

    if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
        os.Mkdir(uploadPath, os.ModePerm)
    }

    dst, err := os.Create(filepath)
    if err != nil {
        http.Error(w, "Unable to create file: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
        return
    }

    now := primitive.NewDateTimeFromTime(time.Now())
    expires := primitive.NewDateTimeFromTime(time.Now().Add(time.Hour))

    downloadLink := fmt.Sprintf("/download/%s", fileID)

    fileDoc := models.File{
        Filename:     handler.Filename, // original filename stored here
        Filepath:     filepath,         // full saved path with UUID prefix
        DownloadLink: downloadLink,
        PasswordHash: hashedPw,
        CreatedAt:    now,
        ExpiresAt:    expires,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err = config.Client.Database("sharenestdb").Collection("files").InsertOne(ctx, fileDoc)
    if err != nil {
        http.Error(w, "Error saving metadata: "+err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "downloadLink": downloadLink,
        "expiresIn":    "1 hour",
    })
}

type DownloadRequest struct {
    Password string `json:"password"`
}


func DownloadFilePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
     fileID := vars["id"]
    log.Printf("[Download] Request received for fileID: %s\n", fileID)

    var req DownloadRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Password == "" {
        log.Printf("[Download] Invalid or missing password in request body: %v\n", err)
        http.Error(w, "Password required in JSON body", http.StatusBadRequest)
        return
    }
    log.Printf("[Download] Password received (hashed verification next)")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var fileDoc models.File

    queryLink := "/download/" + fileID
    log.Printf("[Download] Querying DB for downloadLink: %s\n", queryLink)

    err = config.Client.Database("sharenestdb").Collection("files").FindOne(ctx, bson.M{"downloadLink": queryLink}).Decode(&fileDoc)
    if err != nil {
        log.Printf("[Download] File not found in DB for downloadLink %s: %v\n", queryLink, err)
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }
    log.Printf("[Download] File found in DB: %+v\n", fileDoc)

    if fileDoc.ExpiresAt.Time().Before(time.Now()) {
        log.Printf("[Download] File expired at %v, now %v\n", fileDoc.ExpiresAt.Time(), time.Now())
        http.Error(w, "File expired", http.StatusGone)
        return
    }

    if !checkPasswordHash(req.Password, fileDoc.PasswordHash) {
        log.Printf("[Download] Invalid password attempt for fileID: %s\n", fileID)
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }
    log.Printf("[Download] Password verified successfully")

    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileDoc.Filename))
    w.Header().Set("Content-Type", "application/octet-stream")

    log.Printf("[Download] Serving file from path: %s\n", fileDoc.Filepath)
    http.ServeFile(w, r, fileDoc.Filepath)
}
