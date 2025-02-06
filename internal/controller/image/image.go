package image

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
)

type Cache struct {
	cacheDir    string
	externalApi string
	client      *http.Client
	mu          sync.Mutex
	log         *slog.Logger
}

func NewCache(cacheDir, externalApi string, client *http.Client, logger *slog.Logger) (*Cache, error) {

	err := os.MkdirAll(cacheDir, os.ModePerm)

	if err != nil {
		return nil, fmt.Errorf("error creating image cache directory: %w", err)
	}

	return &Cache{
		cacheDir:    cacheDir,
		externalApi: externalApi,
		client:      client,
		log:         logger,
	}, nil
}

func (c *Cache) GetImage(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	localPath := fmt.Sprintf("%s/%s", c.cacheDir, filename)

	if _, err := os.Stat(localPath); err == nil {
		c.log.Debug(fmt.Sprintf("serving image {%v} from cache", filename))
		http.ServeFile(w, r, localPath)
		return
	}

	externalUrl := fmt.Sprintf("%s/%s", c.externalApi, filename)

	resp, err := c.client.Get(externalUrl)

	if err != nil {
		c.log.Error(fmt.Sprintf("error while retrieving image: %v", err))
		http.Error(w, "error while retrieving image", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.log.Error(fmt.Sprintf("error while retrieving image: %v", err))
		http.Error(w, "error while retrieving image", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		c.log.Error(fmt.Sprintf("error while retrieving image: %v", err))
		http.Error(w, "error while retrieving image", http.StatusInternalServerError)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		err = os.WriteFile(localPath, data, 0644)
		if err != nil {
			c.log.Error(fmt.Sprintf("error while retrieving image: %v", err))
			http.Error(w, "error while retrieving image", http.StatusInternalServerError)
			return
		}
	}

	c.log.Debug(fmt.Sprintf("successfully downloaded image {%v}", filename))

	w.Write(data)

}
