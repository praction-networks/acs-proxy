package config_test

import (
	"testing"

	"github.com/praction-networks/acs-proxy/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {

	err := config.LoadConfig()
	assert.NoError(t, err, "Configuration should load without errors")

}
