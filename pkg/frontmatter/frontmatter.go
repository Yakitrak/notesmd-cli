package frontmatter

import (
	"bytes"
	"errors"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

const (
	Delimiter               = "---"
	NoFrontmatterError      = "note does not contain frontmatter"
	InvalidFrontmatterError = "frontmatter contains invalid YAML"
)

// Parse extracts and parses frontmatter from note content.
// Returns the frontmatter as a map, the body content, and any error.
func Parse(content string) (map[string]interface{}, string, error) {
	var fm map[string]interface{}
	rest, err := frontmatter.Parse(strings.NewReader(content), &fm)
	if err != nil {
		return nil, "", errors.New(InvalidFrontmatterError)
	}
	return fm, string(rest), nil
}

// Format converts a frontmatter map to a YAML string.
func Format(fm map[string]interface{}) (string, error) {
	if len(fm) == 0 {
		return "", nil
	}
	data, err := yaml.Marshal(fm)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// HasFrontmatter checks if content starts with frontmatter delimiters.
func HasFrontmatter(content string) bool {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return false
	}
	return strings.TrimSpace(lines[0]) == Delimiter
}

// SetKey updates or adds a key in the frontmatter, returning the full updated content.
// If no frontmatter exists, it creates new frontmatter with the key.
func SetKey(content, key, value string) (string, error) {
	parsedValue := parseValue(value)

	if !HasFrontmatter(content) {
		// Create new frontmatter
		fm := map[string]interface{}{key: parsedValue}
		fmStr, err := yaml.Marshal(fm)
		if err != nil {
			return "", err
		}
		return Delimiter + "\n" + string(fmStr) + Delimiter + "\n" + content, nil
	}

	var valueNode yaml.Node
	if err := valueNode.Encode(parsedValue); err != nil {
		return "", err
	}

	fm, body, err := parseNode(content)
	if err != nil {
		return "", err
	}

	if i := keyIndex(fm, key); i >= 0 {
		old := fm.Content[i+1]
		valueNode.HeadComment, valueNode.LineComment, valueNode.FootComment = old.HeadComment, old.LineComment, old.FootComment
		fm.Content[i+1] = &valueNode
		resolveAliases(fm, old)
		removeKey(fm, key, i+2)
	} else {
		fm.Content = append(fm.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: key}, &valueNode)
	}

	return formatNode(fm, body)
}

// DeleteKey removes a key from the frontmatter, returning the full updated content.
func DeleteKey(content, key string) (string, error) {
	if !HasFrontmatter(content) {
		return "", errors.New(NoFrontmatterError)
	}

	fm, body, err := parseNode(content)
	if err != nil {
		return "", err
	}

	removeKey(fm, key, 0)

	// If no keys left, return just the body
	if len(fm.Content) == 0 {
		return strings.TrimPrefix(body, "\n"), nil
	}

	return formatNode(fm, body)
}

// parseNode parses the frontmatter as a YAML mapping node, so that edits keep
// the other properties as written (values, order and comments).
func parseNode(content string) (*yaml.Node, string, error) {
	var raw []byte
	format := frontmatter.NewFormat(Delimiter, Delimiter, func(data []byte, v interface{}) error {
		*v.(*[]byte) = data
		return nil
	})
	rest, err := frontmatter.Parse(strings.NewReader(content), &raw, format)
	if err != nil {
		return nil, "", errors.New(InvalidFrontmatterError)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return nil, "", errors.New(InvalidFrontmatterError)
	}
	if len(doc.Content) == 0 {
		return &yaml.Node{Kind: yaml.MappingNode}, string(rest), nil
	}
	if doc.Content[0].Kind != yaml.MappingNode {
		return nil, "", errors.New(InvalidFrontmatterError)
	}
	return doc.Content[0], string(rest), nil
}

// removeKey removes every entry for key at or after index from, so that a
// duplicated key does not keep its old value. Aliases to a removed value are
// replaced by a copy of it, so the frontmatter stays valid.
func removeKey(mapping *yaml.Node, key string, from int) {
	for i := from; i+1 < len(mapping.Content); {
		if mapping.Content[i].Value != key {
			i += 2
			continue
		}
		removed := mapping.Content[i+1]
		mapping.Content = append(mapping.Content[:i], mapping.Content[i+2:]...)
		resolveAliases(mapping, removed)
	}
}

// resolveAliases replaces aliases in node that point into the removed subtree
// with copies of the nodes they point to.
func resolveAliases(node, removed *yaml.Node) {
	anchors := map[*yaml.Node]bool{}
	var collect func(*yaml.Node)
	collect = func(n *yaml.Node) {
		if n.Anchor != "" {
			anchors[n] = true
		}
		for _, c := range n.Content {
			collect(c)
		}
	}
	collect(removed)
	if len(anchors) == 0 {
		return
	}

	var walk func(*yaml.Node)
	walk = func(n *yaml.Node) {
		for i, c := range n.Content {
			if c.Kind == yaml.AliasNode && anchors[c.Alias] {
				copied := *c.Alias
				copied.Anchor = ""
				n.Content[i] = &copied
				continue
			}
			walk(c)
		}
	}
	walk(node)
}

// keyIndex returns the index of key in a mapping node's content, or -1.
func keyIndex(mapping *yaml.Node, key string) int {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return i
		}
	}
	return -1
}

func formatNode(fm *yaml.Node, body string) (string, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(fm); err != nil {
		return "", err
	}
	if err := enc.Close(); err != nil {
		return "", err
	}
	return Delimiter + "\n" + buf.String() + Delimiter + "\n" + body, nil
}

// parseValue attempts to parse the value into appropriate Go types.
// Supports: booleans, arrays (comma-separated in brackets), strings.
func parseValue(value string) interface{} {
	// Try boolean
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	// Try array (comma-separated in brackets)
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		inner := value[1 : len(value)-1]
		if inner == "" {
			return []string{}
		}
		parts := strings.Split(inner, ",")
		result := make([]string, 0, len(parts))
		for _, p := range parts {
			result = append(result, strings.TrimSpace(p))
		}
		return result
	}

	// Default to string
	return value
}
