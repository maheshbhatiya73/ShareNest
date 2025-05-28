    package main

    import (
        "log"
        "net/http"
        "time"

        "ShareNestbackend/config"
        "ShareNestbackend/models"
        "ShareNestbackend/routes"

        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/bson/primitive"

        "github.com/gorilla/mux"
        "context"
        "os"
    )

    func enableCORS(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")


            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    }


    func cleanupExpiredFiles() {
        ticker := time.NewTicker(10 * time.Minute)

        for range ticker.C {
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

            filter := bson.M{"expiresAt": bson.M{"$lt": primitive.NewDateTimeFromTime(time.Now())}}
            cursor, err := config.Client.Database("sharenestdb").Collection("files").Find(ctx, filter)
            if err != nil {
                log.Println("Cleanup find error:", err)
                cancel()
                continue
            }

            var expiredFiles []models.File
            if err := cursor.All(ctx, &expiredFiles); err != nil {
                log.Println("Cleanup cursor error:", err)
                cancel()
                continue
            }

            for _, f := range expiredFiles {
                err := os.Remove(f.Filepath)
                if err != nil {
                    log.Println("Error deleting file:", err)
                }

                _, err = config.Client.Database("sharenestdb").Collection("files").DeleteOne(ctx, bson.M{"_id": f.ID})
                if err != nil {
                    log.Println("Error deleting from DB:", err)
                } else {
                    log.Println("Deleted expired file:", f.Filename)
                }
            }

            cancel()
        }
    }

    func main() {
        config.ConnectDB()

        router := mux.NewRouter()
        routes.FileRoutes(router)

        // Wrap router with CORS middleware
        corsRouter := enableCORS(router)

        go cleanupExpiredFiles()

        log.Println("Server started at :8080")
        if err := http.ListenAndServe(":8080", corsRouter); err != nil {
            log.Fatal(err)
        }
    }