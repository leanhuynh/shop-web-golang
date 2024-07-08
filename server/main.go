package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"myapp/internal/driver"
	"myapp/internal/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
)

const version = "1.0.0"

type config struct {
	port int
	host string
	env  string
	db   struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
	secret_key  string
	session_key string
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	DB            models.DBModel
	Redis         *redisstore.RedisStore
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting HTTP server in %s mode on port %d\n", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// load value from file environment
	err := godotenv.Load()
	if err != nil {
		errorLog.Fatal(err.Error())
	}

	PORT, _ := strconv.Atoi(os.Getenv("PORT"))
	HOST := os.Getenv("HOST")
	DSN := os.Getenv("DB_DSN")
	SMTP_HOST := os.Getenv("SMTP_HOST")
	SMTP_USERNAME := os.Getenv("SMTP_USERNAME")
	SMTP_PASSWORD := os.Getenv("SMTP_PASSWORD")
	SMTP_PORT, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	SECRET_KEY := os.Getenv("SECRET_KEY")
	SESSION_KEY := os.Getenv("SESSION_KEY")

	var cfg config
	flag.IntVar(&cfg.port, "port", PORT, "Server port to listen on")
	flag.StringVar(&cfg.host, "host", HOST, "Host")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&cfg.db.dsn, "dsn", DSN, "DSN")
	flag.StringVar(&cfg.smtp.host, "smtphost", SMTP_HOST, "smtp host")
	flag.StringVar(&cfg.smtp.username, "smtpuser", SMTP_USERNAME, "smtp user")
	flag.StringVar(&cfg.smtp.password, "smtppass", SMTP_PASSWORD, "smtp password")
	flag.IntVar(&cfg.smtp.port, "smtpport", SMTP_PORT, "smtp port")
	flag.StringVar(&cfg.secret_key, "secretkey", SECRET_KEY, "secret key")
	flag.StringVar(&cfg.session_key, "sessionkey", SESSION_KEY, "session key")
	flag.Parse()

	// connect to database
	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err.Error())
	}
	defer conn.Close()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		// Username:     "red-chu6hp5269vccp2pogt0",
		// Password:     "FYWZoKcd0Dyol1iqAzNqIohic1r9nTpD",
		Password:     "",
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		MinIdleConns: 3,
	})

	// kiểm tra liệu có kết nối tới Redis
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		errorLog.Fatalf("Không thể connect tới Redis: %v", err)
	}
	fmt.Printf("Kết nối thành công tới Redis:%s\n", pong)

	// create Session based on Redis database
	store, err := redisstore.NewRedisStore(context.Background(), client)
	if err != nil {
		errorLog.Fatalf("Không thể tạo Redis Store: %v", err)
	}

	store.KeyPrefix("session_")
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 2,
		HttpOnly: true,
		Secure:   true,
	})

	// save CONSTANT VALUES into variable app
	tc := make(map[string]*template.Template)
	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
		DB:            models.DBModel{DB: conn},
		Redis:         store,
	}

	// err = client.Set(context.Background(), "name", "huynhlean", 0).Err()
	// if err != nil {
	// 	fmt.Printf("Failed: %s", err.Error())
	// 	return
	// }
	// val, err := client.Get(context.Background(), "name").Result()
	// if err != nil {
	// 	fmt.Printf("Failed to get: %s", err.Error())
	// 	return
	// }
	// fmt.Printf("value: %s", val)

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err.Error())
		log.Fatal(err.Error())
	}
}
