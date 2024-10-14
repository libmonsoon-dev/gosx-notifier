package gosxnotifier

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
)

type Sound string

const (
	Default   Sound = "'default'"
	Basso     Sound = "Basso"
	Blow      Sound = "Blow"
	Bottle    Sound = "Bottle"
	Frog      Sound = "Frog"
	Funk      Sound = "Funk"
	Glass     Sound = "Glass"
	Hero      Sound = "Hero"
	Morse     Sound = "Morse"
	Ping      Sound = "Ping"
	Pop       Sound = "Pop"
	Purr      Sound = "Purr"
	Sosumi    Sound = "Sosumi"
	Submarine Sound = "Submarine"
	Tink      Sound = "Tink"
)

type Notification struct {
	Message      string //required
	Title        string //optional
	Subtitle     string //optional
	Sound        Sound  //optional
	Link         string //optional
	Sender       string //optional
	Group        string //optional
	AppIcon      string //optional
	ContentImage string //optional
}

func NewNotification(message string) *Notification {
	n := &Notification{Message: message}
	return n
}

func (n *Notification) Push() error {
	cmd, err := n.BuildCommand(nil)
	if err != nil {
		return err
	}

	return cmd.Run()
}

// BuildCommand works like Push, but instead of running *exec.Cmd it just returns it.
func (n *Notification) BuildCommand(ctx context.Context) (*exec.Cmd, error) {
	err := check()
	if err != nil {
		return nil, err
	}

	commandTuples := make([]string, 0)

	//check required commands
	if n.Message == "" {
		return nil, errors.New("please specify a proper message argument")
	} else {
		commandTuples = append(commandTuples, []string{"-message", n.Message}...)
	}

	//add title if found
	if n.Title != "" {
		commandTuples = append(commandTuples, []string{"-title", n.Title}...)
	}

	//add subtitle if found
	if n.Subtitle != "" {
		commandTuples = append(commandTuples, []string{"-subtitle", n.Subtitle}...)
	}

	//add sound if specified
	if n.Sound != "" {
		commandTuples = append(commandTuples, []string{"-sound", string(n.Sound)}...)
	}

	//add group if specified
	if n.Group != "" {
		commandTuples = append(commandTuples, []string{"-group", n.Group}...)
	}

	//add appIcon if specified
	if n.AppIcon != "" {
		img, err := normalizeImagePath(n.AppIcon)
		if err != nil {
			return nil, fmt.Errorf("normalize path %s: %w", n.AppIcon, err)
		}

		commandTuples = append(commandTuples, []string{"-appIcon", img}...)
	}

	//add contentImage if specified
	if n.ContentImage != "" {
		img, err := normalizeImagePath(n.ContentImage)
		if err != nil {
			return nil, fmt.Errorf("normalize path %s: %w", n.ContentImage, err)
		}

		commandTuples = append(commandTuples, []string{"-contentImage", img}...)
	}

	//add url if specified
	url, err := url.Parse(n.Link)
	if err != nil {
		n.Link = ""
	}
	if url != nil && n.Link != "" {
		commandTuples = append(commandTuples, []string{"-open", n.Link}...)
	}

	//add bundle id if specified
	if strings.HasPrefix(strings.ToLower(n.Link), "com.") {
		commandTuples = append(commandTuples, []string{"-activate", n.Link}...)
	}

	//add sender if specified
	if strings.HasPrefix(strings.ToLower(n.Sender), "com.") {
		commandTuples = append(commandTuples, []string{"-sender", n.Sender}...)
	}

	if len(commandTuples) == 0 {
		return nil, errors.New("please provide a Message and Type at a minimum")
	}

	return exec.CommandContext(ctx, binPath, commandTuples...), nil
}

func normalizeImagePath(image string) (string, error) {
	if imagePath, err := filepath.Abs(image); err != nil {
		return "", errors.New("Could not resolve image path of image: " + image)
	} else {
		return imagePath, nil
	}
}
