package server

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

func init() {
	getWd, err := os.Getwd()
	if err != nil {

	}
	parentDir := filepath.Join(getWd, "..")
	err = os.Chdir(parentDir)

	cwd, err := os.Getwd()

	viper.SetDefault("MIGRATIONS_DIR", fmt.Sprintf("%s/migrations", cwd))
}

func (s *server) migrations(pool *pgxpool.Pool) (err error) {
	db, err := goose.OpenDBWithDriver("postgres", pool.Config().ConnConfig.ConnString())
	if err != nil {
		return err
	}
	defer func() {
		if errClose := db.Close(); errClose != nil {
			err = errClose
			return
		}
	}()

	dir := viper.GetString("MIGRATIONS_DIR")
	goose.SetTableName("small_version")
	files, err := ioutil.ReadDir(dir)
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		fmt.Println("hello")
	}

	if err = goose.Run("up", db, dir); err != nil {
		return err
	}
	return
}

func (s *server) logRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		since := time.Since(start).String()
		s.log.Info(fmt.Sprintf("Method: %s, URI: %s, Status: %v, Size: %v, Time: %s", req.Method, req.URL, status, size, since))
		return err
	}
}
