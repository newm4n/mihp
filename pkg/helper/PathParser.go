package helper

import (
	"fmt"
	"strings"
)

func IsTemplateCompatible(template, path string) bool {
	if template == path {
		return true
	}
	if !strings.Contains(template, "{") {
		return false
	}
	templatePaths := strings.Split(template, "/")
	pathPaths := strings.Split(path, "/")
	if len(templatePaths) != len(pathPaths) {
		return false
	}
	for idx, templateElement := range templatePaths {
		pathElement := pathPaths[idx]
		if len(templateElement) > 0 && len(pathElement) > 0 {
			if templateElement[:1] == "{" && templateElement[len(templateElement)-1:] == "}" {
				continue
			} else if templateElement != pathElement {
				return false
			}
		}
	}
	return true
}

// ParsePathParams parse request path param according to path template and extract its values.
func ParsePathParams(template, path string) (map[string]string, error) {
	templatePaths := strings.Split(template, "/")
	pathPaths := strings.Split(path, "/")
	if len(templatePaths) != len(pathPaths) {
		return nil, fmt.Errorf("pathElement length not equals to templateElement length")
	}
	ret := make(map[string]string)
	for idx, templateElement := range templatePaths {
		pathElement := pathPaths[idx]
		if len(templateElement) > 0 && len(pathElement) > 0 {
			if templateElement[:1] == "{" && templateElement[len(templateElement)-1:] == "}" {
				tKey := templateElement[1 : len(templateElement)-1]
				ret[tKey] = pathElement
			} else if templateElement != pathElement {
				return nil, fmt.Errorf("template %s not compatible with path %s", template, path)
			}
		}
	}
	return ret, nil
}
