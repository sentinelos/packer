package solver

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/sentinelos/packer/pkg/constants"
	"github.com/sentinelos/packer/pkg/types"
	"github.com/sentinelos/packer/pkg/types/v1alpha1"
)

// FilesystemPackageLoader loads packages by walking file system tree.
type FilesystemPackageLoader struct {
	*log.Logger
	Context      types.Variables
	pathContexts map[string]types.Variables
	multiErr     *multierror.Error
	pkgFile      *v1alpha1.Pkgfile
	Root         string
	absRootPath  string
	pkgFilePaths []string
	varFilePaths []string
	pkgs         []*v1alpha1.Pkg
}

func (fspl *FilesystemPackageLoader) walkFunc() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fspl.Logger.Printf("error walking %q: %s", path, err)

			return nil
		}

		if info.Name() != "." && strings.HasPrefix(info.Name(), ".") && info.IsDir() {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		switch info.Name() {
		case constants.PkgYaml:
			fspl.pkgFilePaths = append(fspl.pkgFilePaths, path)
		case constants.VarsYaml:
			fspl.varFilePaths = append(fspl.varFilePaths, path)
		}

		return nil
	}
}

// Load implements PackageLoader.
func (fspl *FilesystemPackageLoader) Load() (*LoadResult, error) {
	if fspl.Logger == nil {
		fspl.Logger = log.New(log.Writer(), "[loader] ", log.Flags())
	}

	if fspl.Root == "" {
		fspl.Root = "."
	}

	fspl.pathContexts = make(map[string]types.Variables)

	var err error

	fspl.absRootPath, err = filepath.Abs(fspl.Root)
	if err != nil {
		return nil, err
	}

	if err = fspl.loadPkgfile(); err != nil {
		return nil, err
	}

	fspl.pkgs = nil

	err = filepath.Walk(fspl.Root, fspl.walkFunc())
	if err == nil {
		sort.Slice(fspl.varFilePaths, func(i, j int) bool {
			return filepath.Dir(fspl.varFilePaths[i]) < filepath.Dir(fspl.varFilePaths[j])
		})

		for _, path := range fspl.varFilePaths {
			if err = fspl.loadVariables(path); err != nil {
				fspl.Logger.Printf("error loading variables %q: %s", path, err)
				fspl.multiErr = multierror.Append(fspl.multiErr, fmt.Errorf("error loading variables %q: %w", path, err))
			}

			fspl.Logger.Printf("loaded variables from %q", path)
		}

		for _, path := range fspl.pkgFilePaths {
			var pkg *v1alpha1.Pkg

			pkg, err = fspl.loadPkg(path)
			if err != nil {
				fspl.Logger.Printf("error loading %q: %s", path, err)
				fspl.multiErr = multierror.Append(fspl.multiErr, fmt.Errorf("error loading %q: %w", path, err))

				continue
			}

			fspl.Logger.Printf("loaded pkg %q from %q", pkg.Title, path)
			fspl.pkgs = append(fspl.pkgs, pkg)
		}
	}

	return &LoadResult{
		Pkgfile: fspl.pkgFile,
		Pkgs:    fspl.pkgs,
	}, multierror.Append(fspl.multiErr, err).ErrorOrNil()
}

func (fspl *FilesystemPackageLoader) resolveContext(basePath string) types.Variables {
	context := fspl.Context.Copy()

	dirs := strings.Split(basePath, string(filepath.Separator))

	for i := 0; i <= len(dirs); i++ {
		var subPath string

		if i == 0 {
			subPath = "."
		} else {
			subPath = strings.Join(dirs[:i], string(filepath.Separator))
		}

		if subcontext, ok := fspl.pathContexts[subPath]; ok {
			context.Merge(subcontext)
		}
	}

	return context
}

func (fspl *FilesystemPackageLoader) loadVariables(path string) error {
	absFile, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	basePath, err := filepath.Rel(fspl.absRootPath, absFile)
	if err != nil {
		return err
	}

	baseContext := fspl.resolveContext(filepath.Dir(basePath))

	var vars types.Variables

	if err = vars.Load(path, baseContext); err != nil {
		return err
	}

	fspl.pathContexts[filepath.Dir(basePath)] = vars

	return nil
}

func (fspl *FilesystemPackageLoader) loadPkg(path string) (*v1alpha1.Pkg, error) {
	absFile, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	basePath, err := filepath.Rel(fspl.absRootPath, absFile)
	if err != nil {
		return nil, err
	}

	context := fspl.resolveContext(filepath.Dir(basePath))

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close() //nolint:errcheck

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return v1alpha1.NewPkg(filepath.Dir(basePath), path, contents, context)
}

func (fspl *FilesystemPackageLoader) loadPkgfile() error {
	f, err := os.Open(filepath.Join(fspl.Root, constants.Pkgfile))
	if err != nil {
		if os.IsNotExist(err) {
			fspl.Logger.Printf("skipping %q: %s", constants.Pkgfile, err)

			return nil
		}

		return err
	}

	defer f.Close() //nolint:errcheck

	contents, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	fspl.pkgFile, err = v1alpha1.NewPkgfile(contents)
	if err != nil {
		return fmt.Errorf("error parsing %q: %w", constants.Pkgfile, err)
	}

	fspl.Context.Merge(fspl.pkgFile.Vars)
	fspl.Logger.Printf("loaded %q", constants.Pkgfile)

	return nil
}
