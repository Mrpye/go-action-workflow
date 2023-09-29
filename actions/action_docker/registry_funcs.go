package action_docker

import (
	"archive/tar"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/google/go-containerregistry/pkg/v1/partial"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func extractManifest(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	/*gzipReader, err := gzip.NewReader(bufio.NewReader(file))
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()*/

	tarReader := tar.NewReader(bufio.NewReader(file))

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		fileInfo := header.FileInfo()

		if fileInfo.Name() == "manifest.json" {
			buf := bytes.NewBuffer(nil)
			writer := bufio.NewWriter(buf)

			buffer := make([]byte, 4096)
			for {
				n, err := tarReader.Read(buffer)
				if err != nil && err != io.EOF {
					panic(err)
				}
				if n == 0 {
					break
				}

				_, err = writer.Write(buffer[:n])
				if err != nil {
					return "", err
				}
			}

			err = writer.Flush()
			if err != nil {
				return "", err
			}

			err = file.Close()
			if err != nil {
				return "", err
			}

			//decode the manifest

			return buf.String(), nil
		}
	}
	/*err = file.Close()
	if err != nil {
		return "", err
	}*/
	return "", nil
}

func (c *Registry) loadImage(path string, index bool) (partial.WithRawManifest, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		img, err := crane.Load(path)
		if err != nil {
			return nil, fmt.Errorf("loading %s as tarball: %w", path, err)
		}
		return img, nil
	}

	l, err := layout.ImageIndexFromPath(path)
	if err != nil {
		return nil, fmt.Errorf("loading %s as OCI layout: %w", path, err)
	}

	if index {
		return l, nil
	}

	m, err := l.IndexManifest()
	if err != nil {
		return nil, err
	}
	if len(m.Manifests) != 1 {
		return nil, fmt.Errorf("layout contains %d entries, consider --index", len(m.Manifests))
	}

	desc := m.Manifests[0]
	if desc.MediaType.IsImage() {
		return l.Image(desc.Digest)
	} else if desc.MediaType.IsIndex() {
		return l.ImageIndex(desc.Digest)
	}

	return nil, fmt.Errorf("layout contains non-image (mediaType: %q), consider --index", desc.MediaType)
}

func (m *Registry) GetImageNameTag(image_path string) (string, error) {
	man, err := extractManifest(image_path)
	if err != nil {
		return "", err
	}
	var man_obj Manifest
	err = json.Unmarshal([]byte(man), &man_obj)
	if err != nil {
		return "", err
	}
	if len(man_obj) > 0 {
		if len(man_obj[0].RepoTags) > 0 {
			parts := strings.Split(man_obj[0].RepoTags[0], "/")
			return parts[len(parts)-1], nil
		}
	}
	return "", fmt.Errorf("no package info found")
}

// Upload the image to docker Registry
func (m *Registry) Upload(image_path string, log_level int64) (string, error) {

	//check the config

	name_tag, err := m.GetImageNameTag(image_path)
	if err != nil {
		return "", err
	}

	m.Host = strings.TrimSuffix(m.Host, "/")
	library := m.Library
	if library == "" {
		library = m.UserName
	}
	config := fmt.Sprintf("%s/%s/%s", m.Host, library, name_tag)

	basicAuthn := &authn.Basic{
		Username: m.UserName,
		Password: m.Password,
	}

	var options []remote.Option
	if m.UserName != "" && m.Password != "" {
		withAuthOption := remote.WithAuth(basicAuthn)
		//--insecure
		options = []remote.Option{withAuthOption}
	} else {
		options = []remote.Option{}
	}

	//**********************
	//Check if to ignore ssl
	//**********************
	if m.IgnoreSSL {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		options = append(options, remote.WithTransport(transport))
	}
	if log_level > 0 {
		log.Printf("Importing Image %s to %s \n", image_path, config)
	}
	img, err := m.loadImage(image_path, false)
	if err != nil {
		return "", err
	}

	ref, err := name.ParseReference(config)
	if err != nil {
		return "", err
	}
	//var h v1.Hash
	switch t := img.(type) {
	case v1.Image:
		if err := remote.Write(ref, t, options...); err != nil {
			return "", err
		}
		if _, err = t.Digest(); err != nil {
			return "", err
		}
	case v1.ImageIndex:
		if err := remote.WriteIndex(ref, t, options...); err != nil {
			return "", err
		}
		if _, err = t.Digest(); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("cannot push type (%T) to registry", img)
	}
	//fmt.Printf("Finished Importing Tar %s to %s\n", image_path, config)
	return config, nil
}

// Upload the image from the docker Registry
func (m *Registry) Download(config string, save_name string) error {

	basicAuthn := &authn.Basic{
		Username: m.UserName,
		Password: m.Password,
	}

	imageMap := map[string]v1.Image{}
	var options []remote.Option
	if m.UserName != "" && m.Password != "" {
		withAuthOption := remote.WithAuth(basicAuthn)
		options = []remote.Option{withAuthOption}
	} else {
		options = []remote.Option{}
	}
	//**********************
	//Check if to ignore ssl
	//**********************
	if m.IgnoreSSL {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		options = append(options, remote.WithTransport(transport))

	}

	imageName := config

	ref, err := name.ParseReference(imageName)
	if err != nil {
		log.Fatalf("cannot parse reference of the image %s , detail: %v", imageName, err)
	}

	rmt, err := remote.Get(ref, options...)
	if err != nil {
		return err
	}

	img, err := rmt.Image()
	if err != nil {
		return err
	}

	imageMap[imageName] = img

	//**************
	//Save the image
	//**************
	full_img_name := ref.Context().Name()
	parts := strings.Split(full_img_name, "/")
	if len(parts) < 1 {
		return errors.New("error parsing image name")
	}
	img_name := parts[len(parts)-1]
	tag := "latest"
	if len(parts) > 3 {
		tag = parts[3]
	}
	full_path := img_name + "-" + tag + ".tar"
	if save_name != "" {
		full_path = save_name
	}

	file := path.Dir(full_path)
	err = os.MkdirAll(file, os.ModePerm)
	if err != nil {
		return err
	}

	/*for k := range imageMap {
		fmt.Printf("Downloading Image %v Please Wait!\n", k)
	}*/
	if err := crane.MultiSave(imageMap, full_path); err != nil {
		return fmt.Errorf("error Saving tarball %s: %w", full_path, err)
	}
	/*for k := range imageMap {
		fmt.Printf("Finished Downloading Image %v \n", k)
	}*/

	return nil
}
