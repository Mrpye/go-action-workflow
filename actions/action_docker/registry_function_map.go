package action_docker

import (
	"fmt"
	"strings"

	"github.com/Mrpye/go-action-workflow/workflow"
	dockerparser "github.com/novln/docker-parser"
)

func RemapImage2(w *workflow.Workflow, image string, target_name string, no_tag bool) string {
	return RemapImage(w, image, target_name, no_tag, false, "")
}

// Remap the docker image
func RemapImage(w *workflow.Workflow, image string, target_name string, no_tag bool, use_original bool, original string) string {

	if use_original && original != "" {
		return original
	}

	docker_reg_obj, _ := w.MapTargetConfigValue(target_name, &Registry{})
	if docker_reg_obj == nil {
		return ""
	}
	docker_reg := docker_reg_obj.(*Registry)

	// If the host is not set, use the default
	if docker_reg.Library == "" {
		docker_reg.Library = "library"
		if docker_reg.Host == "docker.io" {
			docker_reg.Library = docker_reg.UserName
		}
	}

	//See if to prefix with the host
	url := ""
	no_host := false
	if docker_reg.Host == "docker.io" {
		no_host = true
	}

	url = fmt.Sprintf("%s/%s", docker_reg.Library, ImageName(image))
	if !no_tag {
		url = fmt.Sprintf("%s/%s", docker_reg.Library, ImageNameTag(image))
	}
	if no_host {
		return url
	}
	return fmt.Sprintf("%s/%s", docker_reg.Host, url)
}

// docker.io/circleci/[slim-base]:latest
func ImageName(image string) string {
	if image != "" {
		parse, _ := dockerparser.Parse(image)
		parts := strings.Split(parse.Repository(), "/")
		return parts[len(parts)-1]
	}
	return ""
}

// docker.io/circleci/[slim-base:latest]
func ImageNameTag(image string) string {
	if image != "" {
		return fmt.Sprintf("%s:%s", ImageName(image), ImageTag(image))
	}
	return ""
}

// docker.io/[circleci]/slim-base:latest
func ImageAccount(image string) string {
	if image != "" {
		parse, _ := dockerparser.Parse(image)
		parts := strings.Split(parse.ShortName(), "/")
		return parts[0]
	}
	return ""
}

// docker.io/[circleci/slim-base]:latest
func ImageShortName(image string) string {
	if image != "" {
		parse, _ := dockerparser.Parse(image)
		return parse.ShortName()
	}
	return ""
}

// [docker.io/circleci/slim-base]:latest
func ImageRegistry(image string) string {
	if image != "" {
		parse, _ := dockerparser.Parse(image)
		return parse.ShortName()
	}
	return ""
}

// docker.io/circleci/slim-base:[latest]
func ImageTag(image string) string {
	if image != "" {
		parse, _ := dockerparser.Parse(image)
		return parse.Tag()
	}
	return ""
}
