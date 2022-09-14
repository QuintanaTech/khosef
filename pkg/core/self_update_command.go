package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v42/github"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"khosef/pkg/cmd"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var (
	releaseBinary = fmt.Sprintf("khosef-%s-amd64", runtime.GOOS)
	owner         = "QuintanaTech"
	repo          = "khosef"
)

type SelfUpdateCommand struct {
	logger *Logger
}

func NewSelfUpdateCommand(verbose *bool) *cobra.Command {
	c := &SelfUpdateCommand{logger: NewNullLogger()}

	return &cobra.Command{
		Use:   "self-update",
		Short: "Update Khosef",
		RunE: cmd.NewRunFnDecorator(func(cmd *cobra.Command, args []string) error {
			if *verbose {
				c.logger = NewPrintLogger()
			}

			return nil
		}, cmd.NewSimpleRunFn(c)),
	}
}

func (s *SelfUpdateCommand) Validate() error {
	return nil
}

func (s *SelfUpdateCommand) Run() error {
	latestRelease, asset, err := getLatestRelease()
	if err != nil {
		return err
	}

	currentVersion := NewVersion()
	if !currentVersion.IsNewer(*latestRelease.TagName) {
		fmt.Println("No newer versions available")

		return nil
	}

	tmpFile, err := s.downloadAsset(asset)
	if err != nil {
		return err
	}

	i, err := NewInstaller(s.logger)
	if err != nil {
		return err
	}

	return i.InstallNix(tmpFile)
}

func (s *SelfUpdateCommand) downloadAsset(asset *github.ReleaseAsset) (*os.File, error) {
	s.logger.info("Starting download of", *asset.BrowserDownloadURL)
	resp, err := http.Get(*asset.BrowserDownloadURL)
	if err != nil {
		return nil, err
	}

	s.logger.info("Retrieved new file with length:", resp.ContentLength)
	bar := progressbar.DefaultBytes(resp.ContentLength, "downloading")

	return s.writeReaderTo(resp.Body, bar)
}

func (s *SelfUpdateCommand) writeReaderTo(sourceStream io.ReadCloser, bar *progressbar.ProgressBar) (*os.File, error) {
	// open output file
	fo, err := ioutil.TempFile(os.TempDir(), "aws2fa")
	if err != nil {
		return nil, err
	}
	s.logger.info("Created temp file at:", fo.Name())

	// close fo on exit and check for its returned error
	defer func() {
		if err := sourceStream.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(io.MultiWriter(fo, bar), sourceStream)
	if err != nil {
		return nil, err
	}

	return fo, nil
}

func getLatestRelease() (*github.RepositoryRelease, *github.ReleaseAsset, error) {
	c := github.NewClient(nil)
	releases, _, err := c.Repositories.ListReleases(context.TODO(), owner, repo, &github.ListOptions{})

	if err != nil {
		return nil, nil, err
	}

	for _, r := range releases {
		if len(r.Assets) < 2 {
			continue
		}

		assets := r.Assets
		for _, a := range assets {
			if 0 == strings.Compare(releaseBinary, *a.Name) {
				return r, a, nil
			}
		}
	}

	return nil, nil, errors.New("unable to discover matching release")
}
