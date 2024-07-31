package images_test

import (
	"slices"
	"testing"

	"github.com/joshmeranda/marina/gateway/images"
)

func TestParseReference(t *testing.T) {
	type TestCase struct {
		Name     string
		Image    string
		Expected images.ImageReference
		Err      error
	}

	cases := []TestCase{
		{
			Name:  "WithAll",
			Image: "registry.example.com/repo/image:v0.1.2@sha256:1234567890abcdef",
			Expected: images.ImageReference{
				Registry:   "registry.example.com",
				Repository: "repo/image",
				Tag:        "v0.1.2",
				Sha:        "1234567890abcdef",
			},
		},
		{
			Name:  "WithRegistryNameAndTag",
			Image: "registry.example.com/repo/image:v0.1.2",
			Expected: images.ImageReference{
				Registry:   "registry.example.com",
				Repository: "repo/image",
				Tag:        "v0.1.2",
			},
		},
		{
			Name:  "WithRegistryNameAndSha",
			Image: "registry.example.com/repo/image@sha256:1234567890abcdef",
			Expected: images.ImageReference{
				Registry:   "registry.example.com",
				Repository: "repo/image",
				Sha:        "1234567890abcdef",
			},
		},
		{
			Name:  "WithNameAndTag",
			Image: "repo/image:v0.1.2",
			Expected: images.ImageReference{
				Repository: "repo/image",
				Tag:        "v0.1.2",
			},
		},

		{
			Name:  "WithRegistryAndName",
			Image: "registry.example.com/repo/image",
			Err:   images.ErrMissingTagOrSha,
		},
		{
			Name:  "WithNameNoNamespace",
			Image: "image",
			Err:   images.ErrMissingTagOrSha,
		},
		{
			Name:  "WithMissingTagAndSha",
			Image: "repo/image",
			Err:   images.ErrMissingTagOrSha,
		},
		{
			Name: "Empty",
			Err:  images.ErrEmptyImage,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			actual, err := images.ParseReference(c.Image)

			if err != c.Err {
				t.Errorf("expected error %v, got %v", c.Err, err)
			}

			if actual != c.Expected {
				t.Errorf("expected %v\n  actual %v", c.Expected, actual)
			}
		})
	}
}

func TestParseMatcher(t *testing.T) {
	type TestCase struct {
		Name     string
		Image    string
		Expected images.ImageMatcher
		Err      error
	}

	cases := []TestCase{
		{
			Name:  "WithAll",
			Image: "registry.example.com/repo/image:v0.1.2@sha256:1234567890abcdef",
			Expected: images.ImageMatcher{
				Registry:   "registry.example.com",
				Repository: "repo/image",
				Tag:        "v0.1.2",
				Sha:        "1234567890abcdef",
			},
		},
		{
			Name:  "WithRegistryNameAndTag",
			Image: "registry.example.com/repo/image:v0.1.2",
			Expected: images.ImageMatcher{
				Registry:   "registry.example.com",
				Repository: "repo/image",
				Tag:        "v0.1.2",
			},
		},
		{
			Name:  "WithRegistryNameAndSha",
			Image: "registry.example.com/repo/image@sha256:1234567890abcdef",
			Expected: images.ImageMatcher{
				Registry:   "registry.example.com",
				Repository: "repo/image",
				Sha:        "1234567890abcdef",
			},
		},
		{
			Name:  "WithNameAndTag",
			Image: "repo/image:v0.1.2",
			Expected: images.ImageMatcher{
				Repository: "repo/image",
				Tag:        "v0.1.2",
			},
		},
		{
			Name:  "WithRegistryAndName",
			Image: "registry.example.com/repo/image",
			Expected: images.ImageMatcher{
				Registry:   "registry.example.com",
				Repository: "repo/image",
			},
		},
		{
			Name:  "WithNameNoNamespace",
			Image: "image",
			Expected: images.ImageMatcher{
				Repository: "image",
			},
		},
		{
			Name:  "WithMissingTagAndSha",
			Image: "repo/image",
			Expected: images.ImageMatcher{
				Repository: "repo/image",
			},
		},
		{
			Name:     "Empty",
			Expected: images.ImageMatcher{},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			actual, err := images.ParseMatcher(c.Image)

			if err != c.Err {
				t.Errorf("expected error %v, got %v", c.Err, err)
			}

			if actual != c.Expected {
				t.Errorf("expected %v\n  actual %v", c.Expected, actual)
			}
		})
	}
}

func TestMatcherMatch(t *testing.T) {
	type TestCase struct {
		Name    string
		Matcher images.ImageMatcher
		Match   []images.ImageReference
	}

	refs := []images.ImageReference{
		{"registry.example.com", "repo/image", "v0.1.2", "1234567890abcdef"},
		{"registry.example.com", "repo/image", "v0.1.2", ""},
		{"registry.example.com", "repo/image", "", "1234567890abcdef"},

		{"", "repo/image", "v0.1.2", "1234567890abcdef"},
		{"", "repo/image", "v0.1.2", ""},
		{"", "repo/image", "", "1234567890abcdef"},

		{"", "repo/image", "latest", ""},
	}

	cases := []TestCase{
		{
			Name:    "MatchAll",
			Matcher: images.ImageMatcher{},
			Match: []images.ImageReference{
				{"registry.example.com", "repo/image", "v0.1.2", "1234567890abcdef"},
				{"registry.example.com", "repo/image", "v0.1.2", ""},
				{"registry.example.com", "repo/image", "", "1234567890abcdef"},

				{"", "repo/image", "v0.1.2", "1234567890abcdef"},
				{"", "repo/image", "v0.1.2", ""},
				{"", "repo/image", "", "1234567890abcdef"},

				{"", "repo/image", "latest", ""},
			},
		},
		{
			Name: "OnlyLatest",
			Matcher: images.ImageMatcher{
				Tag: "latest",
			},
			Match: []images.ImageReference{
				{"", "repo/image", "latest", ""},
			},
		},
		{
			Name: "RegistryRequired",
			Matcher: images.ImageMatcher{
				Registry: "registry.example.com",
			},
			Match: []images.ImageReference{
				{"registry.example.com", "repo/image", "v0.1.2", "1234567890abcdef"},
				{"registry.example.com", "repo/image", "v0.1.2", ""},
				{"registry.example.com", "repo/image", "", "1234567890abcdef"}},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			for _, i := range refs {
				if c.Matcher.Match(i) {
					if !slices.Contains(c.Match, i) {
						t.Errorf("ref '%v' should not have matched", i)
					}
				} else {
					if slices.Contains(c.Match, i) {
						t.Errorf("ref '%v' should have matched", i)
					}
				}
			}
		})
	}
}

func TestImagesAccessList(t *testing.T) {
	accessList := images.ImagesAccessList{
		Allowed: []images.ImageMatcher{
			{
				Repository: "repo/image",
			},
		},
		Blocked: []images.ImageMatcher{
			{
				Tag: "latest",
			},
		},
	}

	type TestCase struct {
		Name    string
		Image   images.ImageReference
		Allowed bool
	}

	cases := []TestCase{
		{
			Name: "OnlyAllowed",
			Image: images.ImageReference{
				Repository: "repo/image",
				Tag:        "v1.2.3",
			},
			Allowed: true,
		},
		{
			Name: "OnlyBlocked",
			Image: images.ImageReference{
				Repository: "repo/different-image",
				Tag:        "latest",
			},
		},
		{
			Name: "AllowedAndBlocked",
			Image: images.ImageReference{
				Repository: "repo/image",
				Tag:        "latest",
			},
		},
		{
			Name: "NotAllowedNotBlocked",
			Image: images.ImageReference{
				Repository: "repo/different-image",
				Tag:        "v1.2.3",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if accessList.IsRefAllowed(c.Image) != c.Allowed {
				t.Errorf("expected %v, got %v", c.Allowed, !c.Allowed)
			}
		})
	}
}

func TestEmptyImagesAccessList(t *testing.T) {
	accessList := images.ImagesAccessList{
		Allowed: []images.ImageMatcher{},
		Blocked: []images.ImageMatcher{},
	}

	refs := []images.ImageReference{
		{
			Repository: "repo/image",
			Tag:        "v1.2.3",
		},
		{
			Repository: "repo/different-image",
			Tag:        "latest",
		},
		{
			Repository: "repo/image",
			Tag:        "latest",
		},
		{
			Repository: "repo/different-image",
			Tag:        "v1.2.3",
		},
	}

	for _, r := range refs {
		if !accessList.IsRefAllowed(r) {
			t.Errorf("reference '%v' should be allowed", r)
		}
	}
}
