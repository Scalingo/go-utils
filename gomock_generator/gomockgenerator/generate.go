package gomockgenerator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/logger"
)

// GenerationConfiguration lets you configure the generation of mocks for your project
type GenerationConfiguration struct {
	// MocksFilePath is the path to the JSON file containing the mock configuration.
	// The directory containing this file is considered as the base package directory.
	MocksFilePath string
	// SignaturesFilename is the filename of the signatures cache, written in the base package directory.
	SignaturesFilename string
	// ConcurrentGoroutines specifies the concurrent amount of goroutines which can execute
	ConcurrentGoroutines int
	// NoGoMod by default we'll consider go modules is enabled, mockgen will be called with -mod=mod to read interfaces in modules instead of default GOPATH
	NoGoMod bool
}

// MocksConfiguration contains the configuration of the mocks to generate.
type MocksConfiguration struct {
	BaseDirectory string `json:"base_directory"`
	// BasePackage is the project base package. E.g. github.com/Scalingo/go-utils
	BasePackage string `json:"base_package"`
	// Mocks contains the configuration of all the mocks to generate
	Mocks []MockConfiguration `json:"mocks"`
}

// MockConfiguration represents a mock and how to generate it.
type MockConfiguration struct {
	// Interface is the name of the interface we need to generate a mock for.
	Interface string `json:"interface"`
	// Mockfile is the location of the generated mock file. Relative path from the root of the
	// project. Defaults to a subfolder of SrcPackage ending with "mock".
	MockFile string `json:"mock_file,omitempty"`
	// SrcPackage is the complete name of the source package. E.g. "model/backup". Defaults to the
	// directory part of Mockfile.
	SrcPackage string `json:"src_package,omitempty"`
	// DstPackage: name of the package of Mockfile. Defaults to the name of the folder of Mockfile.
	DstPackage string `json:"dst_package,omitempty"`
	// External specifies if the generated mock is about a package external to the project.
	External bool `json:"external,omitempty"`
}

// GenerateMocks generates the mocks given in arguments
func GenerateMocks(ctx context.Context, gcfg GenerationConfiguration, mocksCfg MocksConfiguration) error {
	if mocksCfg.BasePackage == "" {
		panic(errors.New(ctx, "BasePackage is mandatory"))
	}
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(ctx, err, "fail to get current directory")
	}
	absMocksFilePath, err := filepath.Abs(gcfg.MocksFilePath)
	if err != nil {
		return errors.Wrap(ctx, err, "fail to resolve absolute path of the mocks file")
	}
	basePackageDirectory := filepath.Dir(absMocksFilePath)
	if mocksCfg.BaseDirectory == "" {
		mocksCfg.BaseDirectory = "."
	}
	packageRootDirectory := filepath.Join(basePackageDirectory, mocksCfg.BaseDirectory)
	err = os.Chdir(packageRootDirectory)
	if err != nil {
		return errors.Wrap(ctx, err, "fail to move to base package directory")
	}
	defer os.Chdir(cwd)

	ctx, log := logger.WithFieldToCtx(ctx, "nb_mocks", len(mocksCfg.Mocks))
	log.WithFields(logrus.Fields{
		"base_package":           mocksCfg.BasePackage,
		"base_package_directory": packageRootDirectory,
	}).Infof("Generating %v mocks", len(mocksCfg.Mocks))

	var mockSigs map[string]string
	mockSigsPath := filepath.Join(packageRootDirectory, gcfg.SignaturesFilename)

	sigs, err := os.ReadFile(mockSigsPath)
	if os.IsNotExist(err) {
		log.Info("No cache signatures file, generates all mocks")
	} else if err != nil {
		return errors.Wrap(ctx, err, "fail to read the signatures cache file")
	} else {
		err = json.Unmarshal(sigs, &mockSigs)
		if err != nil {
			return errors.Wrap(ctx, err, "fail to unmarshal the signatures cache file")
		}
	}
	newMockSigs := make(map[string]string, len(sigs))
	lock := sync.Mutex{}

	var wg sync.WaitGroup
	sem := make(chan bool, gcfg.ConcurrentGoroutines)
	for _, mock := range mocksCfg.Mocks {
		wg.Add(1)
		go func(mock MockConfiguration) {
			defer func() {
				wg.Done()
				<-sem
			}()
			sem <- true
			path, sig, err := generateMock(ctx, gcfg, packageRootDirectory, mocksCfg.BasePackage, mock, mockSigs)
			if err != nil {
				log.Error(err)
				return
			}
			lock.Lock()
			newMockSigs[path] = sig
			lock.Unlock()
		}(mock)
	}
	wg.Wait()

	sigs, err = json.MarshalIndent(newMockSigs, "", "  ")
	if err != nil {
		return errors.Wrap(ctx, err, "fail to marshal the signatures cache file")
	}
	err = os.WriteFile(mockSigsPath, sigs, 0644)
	if err != nil {
		return errors.Wrap(ctx, err, "fail to write the signatures cache file")
	}
	return nil
}

func generateMock(ctx context.Context, gcfg GenerationConfiguration, basePackageDirectory, basePackage string, mock MockConfiguration, sigs map[string]string) (string, string, error) {
	log := logger.Get(ctx)

	if !mock.External {
		if mock.SrcPackage == "" && mock.MockFile == "" {
			return "", "", errors.New(ctx, "SrcPackage or MockFile should be defined to know of guess the source page")
		}

		if mock.SrcPackage == "" {
			mock.SrcPackage = filepath.Dir(mock.MockFile)
		}
	}

	if mock.MockFile == "" {
		basepath := filepath.Base(mock.SrcPackage)
		// If srcPackage is empty, its Base is "."
		if basepath == "." {
			basepath = filepath.Base(basePackage)
		}

		packagePath := path.Join(mock.SrcPackage, fmt.Sprintf("%smock", basepath))
		if mock.DstPackage != "" {
			packagePath = mock.DstPackage
		}

		mock.MockFile = path.Join(
			packagePath,
			fmt.Sprintf("%s_mock.go", strings.ToLower(mock.Interface)),
		)
	}

	if mock.DstPackage == "" {
		dst := filepath.Base(filepath.Dir(mock.MockFile))
		mock.DstPackage = dst
	}

	localSrcPackage := mock.SrcPackage
	if !mock.External {
		mock.SrcPackage = path.Join(basePackage, localSrcPackage)
	}

	mockPath := filepath.Join(basePackageDirectory, mock.MockFile)
	log = log.WithFields(logrus.Fields{
		"mock_file":   mock.MockFile,
		"interface":   mock.Interface,
		"dst_package": mock.DstPackage,
	})
	ctx = logger.ToCtx(ctx, log)
	log.Debug("Generating a mock")
	log.WithFields(logrus.Fields{
		"mock_path":   mockPath,
		"src_package": mock.SrcPackage,
	}).Debug("Mock configuration")

	dir := filepath.Dir(mockPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal("fail to create directory", dir, ":", err)
	}

	selfPackage := ""
	if filepath.Dir(filepath.Join(basePackage, mock.MockFile)) == mock.SrcPackage {
		log = log.WithField("self_package", true)
		ctx = logger.ToCtx(ctx, log)
		selfPackage = "-self_package " + mock.SrcPackage
		mock.DstPackage = filepath.Base(mock.SrcPackage)
	}

	hashSourcePath := mock.SrcPackage
	hashKeyPackage := mock.SrcPackage
	if !mock.External {
		hashSourcePath = filepath.Join(basePackageDirectory, localSrcPackage)
		hashKeyPackage = localSrcPackage
	}

	hashKey := fmt.Sprintf("%s.%s", hashKeyPackage, mock.Interface)
	hash, err := interfaceHash(hashSourcePath, mock.Interface)
	if err != nil {
		return "", "", errors.Wrapf(ctx, err, "fail to get interface hash of %v:%v", mock.SrcPackage, mock.Interface)
	}
	if _, err := os.Stat(mockPath); os.IsNotExist(err) {
		hash = "NOFILE"
	}

	if sigs[hashKey] == hash && hash != "FORCE_REGENERATE" {
		log.Debug("Skipping!")
		return hashKey, hash, nil
	}

	log.WithFields(logrus.Fields{
		"hashkey":  hashKey,
		"expected": sigs[hashKey],
		"current":  hash,
	}).Info("Signature is not matching, regenerating")

	gomod := "--build_flags=--mod=mod"
	if gcfg.NoGoMod {
		gomod = ""
	}

	vendorDir := path.Join(basePackage, "vendor")
	cmd := fmt.Sprintf(
		"mockgen -write_command_comment=false %s -destination %s %s -package %s %s %s && sed -i s,%s,, %s && goimports -w %s",
		gomod, mockPath, selfPackage, mock.DstPackage, mock.SrcPackage, mock.Interface,
		vendorDir, mockPath, mockPath,
	)
	g := exec.Command("sh", "-c", cmd)
	log.WithField("cmd", cmd).Debug("Execute mockgen command")

	stdout, err := g.StdoutPipe()
	if err != nil {
		return "", "", errors.Wrap(ctx, err, "fail to get stdout")
	}
	stderr, err := g.StderrPipe()
	if err != nil {
		return "", "", errors.Wrap(ctx, err, "fail to get stderr")
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	err = g.Start()
	if err != nil {
		return "", "", errors.Wrap(ctx, err, "fail to start")
	}

	err = g.Wait()
	if err != nil {
		return "", "", errors.Wrap(ctx, err, "fail to wait")
	}

	log.Info("Done!")
	return hashKey, hash, nil
}
