// env.go (na raiz, ao lado do go.mod e .env)
package envconfig

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func Load() {
	// Descobre o caminho absoluto deste arquivo (env.go)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("[WARN] não foi possível determinar caminho de env.go; tentando Load() padrão")
		_ = godotenv.Load()
		return
	}

	dir := filepath.Dir(filename)
	envPath := filepath.Join(dir, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("[WARN] .env não carregado em %s: %v\n", envPath, err)
		return
	}

	log.Printf("[INFO] .env carregado de %s\n", envPath)
}
