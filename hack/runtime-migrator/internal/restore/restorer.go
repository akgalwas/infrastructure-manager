package restore

import (
	"context"
	"fmt"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"os"
	"path"
	"sigs.k8s.io/yaml"
)

type Restorer struct {
	backupDir string
}

func NewRestorer(backupDir string) Restorer {
	return Restorer{
		backupDir: backupDir,
	}
}

func (r Restorer) Do(_ context.Context, runtimeID string, shootName string) (v1beta1.Shoot, error) {
	filePath := path.Join(r.backupDir, fmt.Sprintf("backup/%s/%s-to-restore.yaml", runtimeID, shootName))

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return v1beta1.Shoot{}, err
	}

	var shoot v1beta1.Shoot

	err = yaml.Unmarshal(fileBytes, &shoot)
	if err != nil {
		return v1beta1.Shoot{}, err
	}

	return shoot, nil
}
