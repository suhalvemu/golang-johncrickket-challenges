package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	version   = "1.0.0" // can be set at build time using -ldflags
	verbose   bool
	showChars bool
	showLines bool
	showWords bool
	showBytes bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			version = "dev" // fallback if not set
		}
		fmt.Println("ccwc version:", version)
	},
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "ccwc [flags] [file]",
	Short: "ccwc is a word/line/byte/character counter (like wc) in Go",
	Long: `ccwc is a Go implementation of the Unix 'wc' command.
It can count lines, words, characters, and bytes from a file or standard input.`,
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		// No file provided → read from stdin
		if len(args) == 0 {
			return processInput(os.Stdin, "stdin")
		}

		// Iterate over multiple files
		for _, file := range args {
			f, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("failed to open %s: %w", file, err)
			}
			defer f.Close()
			if err := processInput(f, file); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	// Flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVarP(&showLines, "lines", "l", false, "Print line count")
	rootCmd.Flags().BoolVarP(&showWords, "words", "w", false, "Print word count")
	rootCmd.Flags().BoolVarP(&showBytes, "bytes", "c", false, "Print byte count")
	rootCmd.Flags().BoolVarP(&showChars, "chars", "m", false, "Print character count (may differ from bytes in UTF-8)")

	// Default: if no flags, behave like wc (-c -l -w)
	rootCmd.Flags().Lookup("lines").NoOptDefVal = "true"
	rootCmd.Flags().Lookup("words").NoOptDefVal = "true"
	rootCmd.Flags().Lookup("bytes").NoOptDefVal = "true"
}

// processInput reads from a file/stdin and prints stats
func processInput(r io.Reader, label string) error {
	var lines, words, chars, bytes int

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		lines++
		words += len(splitWords(line))
		chars += len([]rune(line)) // rune count = characters
		bytes += len(line) + 1     // +1 for newline
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", label, err)
	}

	// Decide what to print
	if !showLines && !showWords && !showBytes && !showChars {
		showLines, showWords, showBytes = true, true, true
	}

	if showLines {
		fmt.Printf("%8d ", lines)
	}
	if showWords {
		fmt.Printf("%8d ", words)
	}
	if showBytes {
		fmt.Printf("%8d ", bytes)
	}
	if showChars {
		fmt.Printf("%8d ", chars)
	}
	fmt.Printf("%s\n", label)

	return nil
}

// splitWords is a helper for word counting
func splitWords(s string) []string {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)

	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words
}

// Execute runs the CLI
func Execute() {
	rootCmd.AddCommand(versionCmd)
	if err := rootCmd.Execute(); err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
