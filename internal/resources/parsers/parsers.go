package parsers

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	valuesctx "go.mws.cloud/go-sdk/pkg/context/values"
)

type Template []PathElement

func (t Template) String() string {
	realValues := t.RealValues()
	slices.Reverse(realValues)
	return strings.Join(realValues, "/")
}

func (t Template) RealValues() []string {
	result := make([]string, 0, len(t))
	for _, pathElement := range t {
		result = append(result, pathElement.RealValue())
	}
	return result
}

func (t Template) AsID() Template {
	result := make(Template, 0, len(t))
	for _, pathElement := range t {
		pathElement.SearchAfter = true
		result = append(result, pathElement)
	}
	return result
}

type PathElement struct {
	Value       string
	IsConstant  bool
	SearchAfter bool
	Pattern     *regexp.Regexp
}

func (p PathElement) Matches(value string) bool {
	if p.Pattern == nil {
		return true
	}
	return p.Pattern.MatchString(value)
}

func (p PathElement) RealValue() string {
	if p.IsConstant {
		return p.Value
	}
	return "{" + p.Value + "}"
}

// Reference compares the reference string with the template and extracts the key fields.
// Key fields in the context are taken into account if the reference is smaller than the template.
// Input: ref - "service/foo1/foo2/barName/hello/worldName", template.String() - "foo1/foo2/{bar}/hello/{world}",
// serviceName - service.
// Result: map[string]string{"bar":"barName", "world":"worldName"}
func Reference(ctx context.Context, ref string, template Template) (map[string]string, error) {
	splitRef := strings.Split(ref, "/")
	slices.Reverse(splitRef)

	splitLen := len(splitRef)
	templateLen := len(template)

	if splitLen > templateLen {
		return nil, fmt.Errorf("%w: many parts of the reference", ErrReferenceParsing)
	}

	searchAfter := false
	recoveryMode := false
	result := make(map[string]string)

	for index, templateValue := range template {
		inputRefIsHealthy := index <= splitLen-1

		if !inputRefIsHealthy && searchAfter && !recoveryMode {
			return nil, fmt.Errorf("%w: is not enough element - %s", ErrReferenceParsing, templateValue.RealValue())
		}

		if inputRefIsHealthy {
			searchAfter = templateValue.SearchAfter

			if templateValue.IsConstant {
				if templateValue.Value == splitRef[index] {
					continue
				}
				return nil, fmt.Errorf("%w: invalid reference part - '%s'", ErrReferenceParsing, splitRef[index])
			}

			if splitRef[index] == "" {
				return nil, fmt.Errorf("%w: empty part of the reference", ErrReferenceParsing)
			}
			if !templateValue.Matches(splitRef[index]) {
				return nil, fmt.Errorf("%w %s: %s", ErrPatternMatches, templateValue.Value, splitRef[index])
			}
			result[templateValue.Value] = splitRef[index]

			continue
		}

		recoveryMode = true
		searchAfter = false
		if templateValue.IsConstant {
			continue
		}

		val, ok := valuesctx.From(ctx, templateValue.Value)
		if !ok {
			return nil, fmt.Errorf("%w: value '%s' was not found in the context", ErrReferenceParsing, templateValue.Value)
		}
		if !templateValue.Matches(val) {
			return nil, fmt.Errorf("%w %s: %s", ErrPatternMatches, templateValue.Value, val)
		}
		result[templateValue.Value] = val
	}

	return result, nil
}

func CheckOneOfReference(ref string) error {
	if !strings.Contains(ref, "/") {
		return fmt.Errorf("%w: oneOf reference should not consist of a single element", ErrReferenceParsing)
	}
	return nil
}
