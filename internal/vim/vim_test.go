package vim

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/vimcolorschemes/worker/internal/repository"
	"github.com/vimcolorschemes/worker/internal/test"
)

func TestGetVimColorSchemes(t *testing.T) {
	t.Run("should return vim color scheme file URLs of multiple valid files", func(t *testing.T) {
		fileContent1 := `
			hi clear
			syntax reset
			let g:colors_name = "test"
			hi Normal cterm=97
		`
		fileContent2 := `
			hi clear
			syntax reset
			let g:colors_name = "hello"
			hi Normal cterm=97
		`

		server1 := test.MockServer(fileContent1, http.StatusOK)
		defer server1.Close()

		server2 := test.MockServer(fileContent2, http.StatusOK)
		defer server2.Close()

		colorSchemes, err := GetVimColorSchemes([]string{server1.URL, server2.URL})

		if err != nil {
			t.Error("Incorrect result for GetVimColorSchemes, got error")
		}

		expectedColorSchemes := []repository.ColorScheme{{Name: "test", FileURL: server1.URL}, {Name: "hello", FileURL: server2.URL}}

		if !reflect.DeepEqual(colorSchemes, expectedColorSchemes) {
			t.Errorf("Incorrect result for GetVimColorSchemes, got: %s, want: %s", colorSchemes, expectedColorSchemes)
		}
	})

	t.Run("should handle duplicate vim color scheme file URLs", func(t *testing.T) {
		fileContent := `
			hi clear
			syntax reset
			let g:colors_name = "hello"
			hi Normal cterm=97
		`

		server := test.MockServer(fileContent, http.StatusOK)
		defer server.Close()

		colorSchemes, err := GetVimColorSchemes([]string{server.URL, server.URL})

		if err != nil {
			t.Error("Incorrect result for GetVimColorSchemes, got error")
		}

		expectedColorSchemes := []repository.ColorScheme{{Name: "hello", FileURL: server.URL}}

		if !reflect.DeepEqual(colorSchemes, expectedColorSchemes) {
			t.Errorf("Incorrect result for GetVimColorSchemes, got: %s, want: %s", colorSchemes, expectedColorSchemes)
		}
	})

	t.Run("should return empty array and error on invalid vim color scheme files", func(t *testing.T) {
		fileContent1 := `
			hi clear
			syntax reset
			let g:color='hello'
		`
		fileContent2 := `
			hi clear
			syntax reset
			hi Normal cterm=97
		`
		fileContent3 := `
			hi clear
			syntax reset
			let g:color='test'
			hi Normal cterm=97
		`

		server1 := test.MockServer(fileContent1, http.StatusOK)
		defer server1.Close()

		server2 := test.MockServer(fileContent2, http.StatusOK)
		defer server2.Close()

		server3 := test.MockServer(fileContent3, http.StatusOK)
		defer server3.Close()

		colorSchemes, err := GetVimColorSchemes([]string{server1.URL, server2.URL, server3.URL})

		expectedColorSchemes := []repository.ColorScheme{}

		if err == nil {
			t.Error("Incorrect result for GetVimColorSchemes, got no error")
		}

		if !reflect.DeepEqual(colorSchemes, expectedColorSchemes) {
			t.Errorf("Incorrect result for GetVimColorSchemes, got: %s, want: %s", colorSchemes, expectedColorSchemes)
		}
	})

	t.Run("should ignore invalid file URLs", func(t *testing.T) {
		fileContent := `
			hi clear
			syntax reset
			let g:colors_name = "test"
			hi Normal cterm=97
		`

		server := test.MockServer(fileContent, http.StatusOK)
		defer server.Close()

		colorSchemes, err := GetVimColorSchemes([]string{server.URL, "wrong url"})

		if err != nil {
			t.Error("Incorrect result for GetVimColorSchemes, got error")
		}

		expectedColorSchemes := []repository.ColorScheme{{Name: "test", FileURL: server.URL}}

		if !reflect.DeepEqual(colorSchemes, expectedColorSchemes) {
			t.Errorf("Incorrect result for GetVimColorSchemes, got: %s, want: %s", colorSchemes, expectedColorSchemes)
		}
	})
}

var validTests = []struct {
	fileContent string
	name        string
}{
	{fileContent: `
		hi clear
		let g:colors_name = "test"
		syntax reset`, name: "test"},
	{fileContent: `
		hi clear
		let g:color_name = "hello-world"
		syntax reset`, name: "hello-world"},
	{fileContent: `
		hi clear
		let g:colors_name="hello_world"
		syntex reset`, name: "hello_world"},
	{fileContent: `
		hi clear
		let g:color_name="hello (world)"
		syntex reset`, name: "hello (world)"},
	{fileContent: `
		hi clear
		let colors_name = "abcd1234"
		syntex reset`, name: "abcd1234"},
	{fileContent: `
		hi clear
		let color_name = "TEST"
		syntex reset`, name: "TEST"},
	{fileContent: `
		hi clear
		let colors_name="TEst"
		syntex reset`, name: "TEst"},
	{fileContent: `
		hi clear
		let color_name="test"
		syntex reset`, name: "test"},
	{fileContent: `
		hi clear
		let g:colors_name = 'test'
		syntex reset`, name: "test"},
	{fileContent: `
		hi clear
		let g:color_name = 'test'
		syntex reset`, name: "test"},
	{fileContent: `
		hi clear
		let g:colors_name='test'
		syntex reset`, name: "test"},
	{fileContent: `
		hi clear
		let g:color_name='test'
		syntex reset`, name: "test"},
	{fileContent: `
		hi clear
		let colors_name = 'test'
		syntex reset`, name: "test"},
	{fileContent: `
		hi clear
		let color_name = 'test'
		syntex reset`, name: "test"},
	{fileContent: `
		hi clear
		let colors_name='test'
		syntex reset`, name: "test"},
	{fileContent: `
		hi clear
		let color_name='test'
		syntex reset`, name: "test"},
}

func TestGetVimColorSchemeName(t *testing.T) {
	t.Run("should return the vim color scheme name if the file is valid", func(t *testing.T) {
		for _, item := range validTests {
			name, err := GetVimColorSchemeName(&item.fileContent)
			if err != nil {
				t.Error("Incorrect result for GetVimColorSchemeName, got error")
			}
			if name != item.name {
				t.Errorf("Incorrect result for GetVimColorSchemeName, got: %s, want: %s", name, item.name)
			}
		}
	})

	invalid := []string{
		`
			hi clear
			let g:colors_name = "test
			syntax reset
		`,
		`
			hi clear
			let g:color_names = "test"
			syntax reset
		`,
		`
			hi clear
			g:colors_name="test"
			syntax reset
		`,
		`
			hi clear
			let g:color_name="'test'"
			syntax reset
		`,
		`
			hi clear
			let colors_name = "{test}"
			syntax reset
		`,
		`
			hi clear
			let color_name = expand("test")
			syntax reset
		`,
	}
	t.Run("should return an error if the file is invalid", func(t *testing.T) {
		for _, fileContent := range invalid {
			name, err := GetVimColorSchemeName(&fileContent)
			if err == nil {
				t.Error("Incorrect result for GetVimColorSchemeName, got no error")
			}
			if name != "" {
				t.Errorf("Incorrect result for GetVimColorSchemeName, got: %s, want: %s", name, "")
			}
		}
	})
}

func TestIsVimColorScheme(t *testing.T) {
	valid := []string{
		`
			exe "hi! Normal"        .s:fg_foreground  .s:bg_normal      .s:fmt_none
		`,
		`
			hi Normal ctermbg=254 ctermfg=237 guibg=#e8e9ec guifg=#33374c
		`,
	}
	t.Run("should return true if the file has necessary content", func(t *testing.T) {
		for _, fileContent := range valid {
			if !IsVimColorScheme(&fileContent) {
				t.Errorf("Incorrect result for IsVimColorScheme, got false for: %s", fileContent)
			}
		}
	})

	invalid := []string{
		`
			Normal ctermbg=254 ctermfg=237 guibg=#e8e9ec guifg=#33374c
		`,
	}
	t.Run("should return false if the file does not have necessary content", func(t *testing.T) {
		for _, fileContent := range invalid {
			if IsVimColorScheme(&fileContent) {
				t.Errorf("Incorrect result for IsVimColorScheme, got true for: %s", fileContent)
			}
		}
	})
}
