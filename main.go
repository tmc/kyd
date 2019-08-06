// Command kyd diffs two yaml documents.
// Only documents that appear in the second file are emitted.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type yamldiffer struct {
	// references to yaml documents
	a, b []*yaml.Node
}

// Diff emits the list of yaml nodes that are present in the second file but not the first.
func (y *yamldiffer) Diff() ([]*yaml.Node, error) {
	result := []*yaml.Node{}

	lhs := map[string]bool{}
	for _, n := range y.a {
		h, err := hash(n)
		if err != nil {
			return nil, err
		}
		lhs[h] = true
	}
	for _, n := range y.b {
		h, err := hash(n)
		if err != nil {
			return nil, err
		}
		if !lhs[h] {
			result = append(result, n)
		}
	}
	return result, nil
}

// hash generates a key for a yaml node.
func hash(y *yaml.Node) (string, error) {
	buf := &bytes.Buffer{}
	enc := yaml.NewEncoder(buf)
	if err := enc.Encode(y); err != nil {
		return "", err
	}
	if err := enc.Close(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func decode(r io.Reader) ([]*yaml.Node, error) {
	result := []*yaml.Node{}
	dec := yaml.NewDecoder(r)
	for {
		d := &yaml.Node{}
		err := dec.Decode(d)
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, err
		}
		result = append(result, d)
	}
	return result, nil
}

var cmdDiff = &cobra.Command{
	Use:  "kyd [file1] [file2]",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		f1, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "issue opening first file:", err)
			os.Exit(1)
		}

		f2, err := os.Open(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "issue opening second file:", err)
			os.Exit(1)
		}

		a, err := decode(f1)
		if err != nil {
			fmt.Fprintln(os.Stderr, "issue decoding first file:", err)
			os.Exit(1)
		}
		b, err := decode(f2)
		if err != nil {
			fmt.Fprintln(os.Stderr, "issue decoding second file:", err)
			os.Exit(1)
		}

		d := &yamldiffer{
			a: a,
			b: b,
		}

		diffSet, err := d.Diff()

		if err != nil {
			fmt.Fprintln(os.Stderr, "issue diffing:", err)
			os.Exit(1)
		}
		enc := yaml.NewEncoder(os.Stdout)
		for _, doc := range diffSet {
			if err := enc.Encode(doc); err != nil {
				fmt.Fprintln(os.Stderr, "issue encoding:", err)
				os.Exit(1)
			}
		}
		if err := enc.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "issue closing:", err)
			os.Exit(1)
		}
	},
}

func main() {
	cmdDiff.Execute()
}
