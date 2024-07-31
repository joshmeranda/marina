package images

import (
	"fmt"
	"strings"
)

var (
	ErrMissingRepository = fmt.Errorf("missing repository")
	ErrMissingTagOrSha   = fmt.Errorf("missing tag or sha")
	ErrEmptyImage        = fmt.Errorf("image may not be empty")
)

type rawRef struct {
	Registry   string
	Repository string
	Tag        string
	Sha        string
}

func parse(s string) (ImageReference, error) {
	m := ImageReference{}

	digestSplit := strings.Split(s, "@sha256:")
	if len(digestSplit) == 2 {
		m.Sha = digestSplit[1]
		s = digestSplit[0]
	}

	tagSplit := strings.Split(s, ":")
	if len(tagSplit) == 2 {
		m.Tag = tagSplit[1]
		s = tagSplit[0]
	}

	repoSplit := strings.Split(s, "/")
	if len(repoSplit) == 3 {
		m.Registry = repoSplit[0]
		m.Repository = repoSplit[1] + "/" + repoSplit[2]
	} else if len(repoSplit) == 2 {
		m.Repository = repoSplit[0] + "/" + repoSplit[1]
	} else {
		m.Repository = repoSplit[0]
	}

	return m, nil
}

type ImageReference rawRef

func ParseReference(s string) (ImageReference, error) {
	if s == "" {
		return ImageReference{}, ErrEmptyImage
	}

	ref, err := parse(s)
	if err != nil {
		return ImageReference{}, err
	}

	if ref.Repository == "" {
		return ImageReference{}, ErrMissingRepository
	}

	if ref.Tag == "" && ref.Sha == "" {
		return ImageReference{}, ErrMissingTagOrSha
	}

	return ref, nil
}

func (r ImageReference) String() string {
	s := strings.Builder{}

	if r.Registry != "" {
		s.WriteString(r.Registry + "/")
	}

	s.WriteString(r.Repository)

	if r.Tag != "" {
		s.WriteString(":" + r.Tag)
	}

	if r.Sha != "" {
		s.WriteString("@sha256:" + r.Sha)
	}

	return s.String()
}

type ImageMatcher rawRef

func ParseMatcher(s string) (ImageMatcher, error) {
	ref, err := parse(s)
	if err != nil {
		return ImageMatcher{}, err
	}

	return ImageMatcher(ref), nil
}

func (m ImageMatcher) Match(image ImageReference) bool {
	return !(m.Registry != "" && m.Registry != image.Registry) &&
		!(m.Repository != "" && m.Repository != image.Repository) &&
		!(m.Tag != "" && m.Tag != image.Tag) &&
		!(m.Sha != "" && m.Sha != image.Sha)
}

type ImagesAccessList struct {
	Allowed []ImageMatcher
	Blocked []ImageMatcher
}

func (l ImagesAccessList) IsRefAllowed(ref ImageReference) bool {
	for _, m := range l.Blocked {
		if m.Match(ref) {
			return false
		}
	}

	if len(l.Allowed) == 0 {
		return true
	}

	for _, m := range l.Allowed {
		if m.Match(ref) {
			return true
		}
	}

	return false
}

func (l ImagesAccessList) IsImageAllowed(image string) (bool, error) {
	ref, err := ParseReference(image)
	if err != nil {
		return false, fmt.Errorf("failed to parse image reference: %w", err)
	}

	return l.IsRefAllowed(ref), nil
}
