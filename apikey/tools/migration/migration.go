package main

import (
	"context"
	"embed"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/afero"
)

//go:embed files
var Files embed.FS

func main() {
	const connString = "postgres://api_key:api_key@localhost:5435/api_key?sslmode=disable"
	m, err := Get(context.Background(), connString, Files)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		panic(err)
	}
}

// Get will return the migration prepared for the certain files embedded into the binary
func Get(_ context.Context, conn string, embeds ...embed.FS) (*migrate.Migrate, error) {
	afs, err := FromEmbeds(embeds...)
	if err != nil {
		return nil, err
	}

	d, err := iofs.New(afero.NewIOFS(afs), "files")
	if err != nil {
		return nil, err
	}

	return migrate.NewWithSourceInstance("iofs", d, conn)
}

func FromEmbeds(embeds ...embed.FS) (afero.Fs, error) {
	afs := afero.NewMemMapFs()
	if err := afs.Mkdir("files", os.ModePerm); err != nil {
		return afs, err
	}
	for _, e := range embeds {
		err := fs.WalkDir(e, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				asd, err := e.Open(path)
				if err != nil {
					return err
				}
				c, err := ioutil.ReadAll(asd)
				if err != nil {
					return err
				}
				err = afero.WriteFile(afs, path, c, os.ModePerm)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return afs, err
		}
	}
	return afs, nil
}
