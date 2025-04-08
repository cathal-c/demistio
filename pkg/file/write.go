package file

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

func Write(ctx context.Context, path string, content []byte) error {
	log := zerolog.Ctx(ctx)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("failed to close file")
		}
	}()

	res, err := f.Write(content)
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	log.Debug().Msgf("wrote %d bytes to %s", res, path)

	return nil
}
