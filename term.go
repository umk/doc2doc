package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"unicode"

	"golang.org/x/term"
)

// readKeySilent writes the prompt to the terminal and waits for a single key press silently.
// It returns the key pressed as a rune, handling multi-byte characters.
func readKeySilent(prompt ...string) (rune, error) {
	fd := int(os.Stdin.Fd())

	// Save the original terminal state to restore later
	oldState, err := term.GetState(fd)
	if err != nil {
		return 0, err
	}
	defer term.Restore(fd, oldState)

	// Put terminal into raw mode to read single characters
	if _, err := term.MakeRaw(fd); err != nil {
		return 0, err
	}

	// Write the prompt
	for _, p := range prompt {
		fmt.Print(p)
	}

	// Create a buffered reader for stdin
	reader := bufio.NewReader(os.Stdin)

	// Read a single UTF-8 encoded rune
	r, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}

	return r, nil
}

// readKeyOfSilent writes the prompt and waits for a key press that matches one of the provided runes silently.
// It keeps prompting until a valid key is pressed and returns that rune.
func readKeyOfSilent(prompt string, keys ...rune) (rune, error) {
	fmt.Print(prompt)
	for {
		key, err := readKeySilent()
		if err != nil {
			return 0, err
		}
		// Convert to lower case if it's a letter for case-insensitive matching
		keyLower := unicode.ToLower(key)
		if slices.Index(keys, keyLower) >= 0 {
			return key, nil
		}
	}
}

// readKeyOrDefaultOfSilent writes the prompt and waits for a key press that matches one of the provided runes silently.
// If the Enter key (rune '\r') is pressed, it returns the default rune.
func readKeyOrDefaultOfSilent(prompt string, def rune, others ...rune) (rune, error) {
	// Include '\r' (Enter key) as a valid key
	keys := append([]rune{'\r', def}, others...)
	result, err := readKeyOfSilent(prompt, keys...)
	if err != nil {
		return 0, err
	}
	if result == '\r' {
		return def, nil
	}
	return result, nil
}

// readKeyOrDefaultOf writes the prompt and waits for a key press that matches one of the provided runes.
// If the Enter key (rune '\r') is pressed, it echoes and returns the default rune.
func readKeyOrDefaultOf(prompt string, def rune, others ...rune) (rune, error) {
	result, err := readKeyOrDefaultOfSilent(prompt, def, others...)
	if err != nil {
		fmt.Println() // Add a newline
		return 0, err
	}
	fmt.Println(string(result)) // Echo the character and add a newline
	return result, nil
}
