package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	inFilePath := os.Args[1]
	outFilePath := os.Args[2]

	p := newCodeParser()
	parsedFile, err := p.Parse(inFilePath)
	if err != nil {
		log.Fatalf("failed to parse file: %s", err)
	}

	out, err := os.Create(outFilePath)
	if err != nil {
		log.Fatalf("failed to create file %s: %s", outFilePath, err)
	}

	g := newCodeGenerator(out, parsedFile)
	g.Generate()
}

// parsedFile результат парсинга файла
type parsedFile struct {
	PackageName string
	ApiHandlers map[string]apiHandler
	ValStructs  map[string]valStruct
}

// apiHandler обработчик запросов API
type apiHandler struct {
	Name    string
	Methods []apiMethod
}

// apiMethod обработчик метода API
type apiMethod struct {
	Name        string
	HandlerName string
	RequestName string
	Api         apiMeta
}

// apiMeta мета-информация о методе API
type apiMeta struct {
	Url    string
	Auth   bool
	Method string
}

// valStruct структура с валидацией
type valStruct struct {
	Name   string
	Fields []valField
}

// valField поле структуры с валидацией
type valField struct {
	Name  string
	Type  string
	Rules valRules
}

// valRules правила валидации
type valRules struct {
	ParamName string
	Required  bool
	Min       bool
	MinValue  int
	Max       bool
	MaxValue  int
	Enum      []string
	Default   string
}

// codeParser собирает данные для кодогенерации
type codeParser struct {
	apigenPrefix   string
	matchFirstCap  *regexp.Regexp
	matchAllCap    *regexp.Regexp
	matchValidator *regexp.Regexp
}

// newCodeParser создает парсер
func newCodeParser() *codeParser {
	return &codeParser{
		apigenPrefix:   "// apigen:api",
		matchFirstCap:  regexp.MustCompile("(.)([A-Z][a-z]+)"),
		matchAllCap:    regexp.MustCompile("([a-z0-9])([A-Z])"),
		matchValidator: regexp.MustCompile("`apivalidator:\"(.*)\"`"),
	}
}

// Parse парсит файл
func (p *codeParser) Parse(filePath string) (*parsedFile, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	result := &parsedFile{
		PackageName: node.Name.Name,
		ApiHandlers: make(map[string]apiHandler),
		ValStructs:  make(map[string]valStruct),
	}

	for _, decl := range node.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			p.parseFunc(result, decl)
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				if t, ok := spec.(*ast.TypeSpec); ok {
					if st, ok := t.Type.(*ast.StructType); ok {
						p.parseStruct(result, t.Name.Name, st)
					}
				}
			}
		}
	}

	return result, nil
}

// parseFunc парсит функции для обработки запросов
func (p *codeParser) parseFunc(result *parsedFile, f *ast.FuncDecl) {
	if f.Doc != nil {
		api := apiMeta{}
		for _, c := range f.Doc.List {
			if strings.HasPrefix(c.Text, p.apigenPrefix) {
				jsonStr := c.Text[len(p.apigenPrefix):]
				if err := json.Unmarshal([]byte(jsonStr), &api); err == nil {
					break
				}
			}
		}

		if api.Url != "" {
			if handlerName := p.getFuncRecvName(f); handlerName != "" {
				if _, exists := result.ApiHandlers[handlerName]; !exists {
					result.ApiHandlers[handlerName] = apiHandler{
						Name: handlerName,
					}
				}

				if requestType, ok := f.Type.Params.List[1].Type.(*ast.Ident); ok {
					h := result.ApiHandlers[handlerName]
					h.Methods = append(h.Methods, apiMethod{
						Name:        f.Name.Name,
						HandlerName: handlerName,
						RequestName: requestType.Name,
						Api:         api,
					})
					result.ApiHandlers[handlerName] = h
				}
			}
		}
	}
}

// parseStruct парсит структуры с валидацией параметров
func (p *codeParser) parseStruct(result *parsedFile, structName string, structType *ast.StructType) {
	for _, f := range structType.Fields.List {
		if f.Tag != nil {
			if match := p.matchValidator.FindStringSubmatch(f.Tag.Value); len(match) > 0 {
				if _, exists := result.ValStructs[structName]; !exists {
					result.ValStructs[structName] = valStruct{
						Name: structName,
					}
				}

				rules := valRules{
					ParamName: p.toSnakeCase(f.Names[0].Name),
				}

				for _, rule := range strings.Split(match[1], ",") {
					ruleParts := strings.Split(rule, "=")
					switch ruleParts[0] {
					case "required":
						rules.Required = true
					case "paramname":
						rules.ParamName = ruleParts[1]
					case "min":
						rules.Min = true
						rules.MinValue, _ = strconv.Atoi(ruleParts[1])
					case "max":
						rules.Max = true
						rules.MaxValue, _ = strconv.Atoi(ruleParts[1])
					case "enum":
						rules.Enum = strings.Split(ruleParts[1], "|")
					case "default":
						rules.Default = ruleParts[1]
					}
				}

				v := result.ValStructs[structName]
				v.Fields = append(v.Fields, valField{
					Name:  f.Names[0].Name,
					Type:  strings.Title(f.Type.(*ast.Ident).Name),
					Rules: rules,
				})
				result.ValStructs[structName] = v
			}
		}
	}
}

// toSnakeCase переводит строку в snake_case
func (p *codeParser) toSnakeCase(str string) string {
	snake := p.matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = p.matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// getFuncRecvName получает имя структуры, которой принадлежит метод
func (p *codeParser) getFuncRecvName(f *ast.FuncDecl) string {
	if f.Recv != nil {
		for _, fi := range f.Recv.List {
			// if pointer
			if fi, ok := fi.Type.(*ast.StarExpr); ok {
				if fi, ok := fi.X.(*ast.Ident); ok {
					return fi.Name
				}
			}

			if fi, ok := fi.Type.(*ast.Ident); ok {
				return fi.Name
			}
		}
	}

	return ""
}

// codeGenerator генерирует код для обработчиков API и валидации
type codeGenerator struct {
	out  io.Writer
	data *parsedFile
}

// newCodeGenerator создает кодогенератор
func newCodeGenerator(out io.Writer, data *parsedFile) *codeGenerator {
	return &codeGenerator{
		out:  out,
		data: data,
	}
}

// Generate генерация кода
func (g *codeGenerator) Generate() {
	fmt.Fprintf(g.out, "package %s\n", g.data.PackageName)

	g.writeImports()
	g.writeHelpers()
	g.writeValidators()

	serveTpl := g.newServeTpl()
	wrapperTpl := g.newWrapperTpl()
	validatorTpl := g.newValidatorTpl()

	for _, h := range g.data.ApiHandlers {
		serveTpl.Execute(g.out, h)

		for _, m := range h.Methods {
			wrapperTpl.Execute(g.out, m)
		}
	}

	for _, v := range g.data.ValStructs {
		validatorTpl.Execute(g.out, v)
	}
}

// writeImports пишет импорты
func (g *codeGenerator) writeImports() {
	fmt.Fprint(g.out, `
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)
`)
}

// writeHelpers пишет вспомогательные функции
func (g *codeGenerator) writeHelpers() {
	fmt.Fprint(g.out, `
type ApiResponse struct {
	Data       interface{} `+"`"+`json:"response,omitempty"`+"`"+`
	Error      string      `+"`"+`json:"error"`+"`"+`
	statusCode int
}

func sendResponse(w http.ResponseWriter, response ApiResponse) {
	if response.statusCode != 0 {
		w.WriteHeader(response.statusCode)
	}

	js, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getParams(r *http.Request) (params url.Values) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		params, _ = url.ParseQuery(string(body))
	} else {
		params = r.URL.Query()
	}

	return params
}`)
}

// writeValidators пишет функции-валидаторы
func (g *codeGenerator) writeValidators() {
	fmt.Fprint(g.out, `
func validateEmpty(field string, val string) error {
	if val == "" {
		return ApiError{http.StatusBadRequest, fmt.Errorf("%s must me not empty", field)}
	}

	return nil
}

func validateMinInt(field string, val int, minVal int) error {
	if val < minVal {
		return ApiError{http.StatusBadRequest, fmt.Errorf("%v must be >= %v", field, minVal)}
	}

	return nil
}

func validateMinString(field string, val string, minVal int) error {
	if len(val) < minVal {
		return ApiError{http.StatusBadRequest, fmt.Errorf("%v len must be >= %v", field, minVal)}
	}

	return nil
}

func validateMaxInt(field string, val int, minVal int) error {
	if val > minVal {
		return ApiError{http.StatusBadRequest, fmt.Errorf("%v must be <= %v", field, minVal)}
	}

	return nil
}

func validateMaxString(field string, val string, minVal int) error {
	if len(val) > minVal {
		return ApiError{http.StatusBadRequest, fmt.Errorf("%v len must be <= %v", field, minVal)}
	}

	return nil
}

func validateEnum(field string, val string, enum []string) error {
	for _, valid := range enum {
		if valid == val {
			return nil
		}
	}

	return ApiError{http.StatusBadRequest, fmt.Errorf("%v must be one of [%v]", field, strings.Join(enum, ", "))}
}
`)
}

// newServeTpl cоздает шаблон для генерации кода обработчика запросов API
func (g *codeGenerator) newServeTpl() *template.Template {
	return template.Must(template.New("serveTpl").Parse(`func (h *{{ .Name }}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		out interface{}
		err error
	)

	switch r.URL.Path {
	{{ range .Methods }}case "{{ .Api.Url }}":
		out, err = h.wrapper{{ .Name }}(w, r)
	{{ end }}default:
		err = ApiError{Err: fmt.Errorf("unknown method"), HTTPStatus: http.StatusNotFound}
	}

	if err == nil {
		sendResponse(w, ApiResponse{Data: out})
	} else if errApi, ok := err.(ApiError); ok {
		sendResponse(w, ApiResponse{Error: errApi.Error(), statusCode: errApi.HTTPStatus})
	} else {
		sendResponse(w, ApiResponse{Error: err.Error(), statusCode: http.StatusInternalServerError})
	}
}

`))
}

// newWrapperTpl cоздает шаблон для генерации кода обработчика метода API
func (g *codeGenerator) newWrapperTpl() *template.Template {
	return template.Must(template.New("wrapperTpl").Parse(`func (h *{{ .HandlerName }}) wrapper{{ .Name }}(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	{{ if .Api.Auth -}}
	if r.Header.Get("X-Auth") != "100500" {
		return nil, ApiError{http.StatusForbidden, fmt.Errorf("unauthorized")}
	}

	{{ end -}}

	{{ if .Api.Method -}}
	if r.Method != "{{ .Api.Method }}" {
		return nil, ApiError{http.StatusNotAcceptable, fmt.Errorf("bad method")}
	}

	{{ end -}}

	in, err := new{{ .RequestName }}(getParams(r));
	if err != nil {
		return nil, err
	}
	
	return h.{{ .Name }}(r.Context(), in)
}

`))
}

// newValidatorTpl cоздает шаблон для генерации кода валидации структуры
func (g *codeGenerator) newValidatorTpl() *template.Template {
	return template.Must(template.New("validatorTpl").Parse(`func new{{ .Name }}(v url.Values) ({{ .Name }}, error) {
	var err error
	s := {{ .Name }}{}

	{{ range .Fields }}// {{ .Name }}
	
	{{- if eq .Type "Int" }}
	s.{{ .Name }}, err = strconv.Atoi(v.Get("{{ .Rules.ParamName }}"))
	if err != nil {
		return s, ApiError{http.StatusBadRequest, fmt.Errorf("{{ .Rules.ParamName }} must be int")}
	}
	{{- else }}
	s.{{ .Name }} = v.Get("{{ .Rules.ParamName }}")
	{{- end -}}

	{{- if .Rules.Default }}
	if s.{{ .Name }} == "" {
		s.{{ .Name }} = "{{ .Rules.Default }}"
	}
	{{- end -}}

	{{- if .Rules.Required }}
	if err = validateEmpty("{{ .Rules.ParamName }}", s.{{ .Name }}); err != nil {
		return s, err
	}
	{{- end -}}

	{{- if .Rules.Min }}
	if err = validateMin{{ .Type }}("{{ .Rules.ParamName }}", s.{{ .Name }}, {{ .Rules.MinValue }}); err != nil {
		return s, err
	}
	{{- end -}}

	{{- if .Rules.Max }}
	if err = validateMax{{ .Type }}("{{ .Rules.ParamName }}", s.{{ .Name }}, {{ .Rules.MaxValue }}); err != nil {
		return s, err
	}
	{{- end -}}

	{{- if .Rules.Enum }}
	enum := []string{ {{ range $index, $element := .Rules.Enum }}{{ if $index }}, {{ end }}"{{ $element }}"{{ end }} }
	if err = validateEnum("{{ .Rules.ParamName }}", s.{{ .Name }}, enum); err != nil {
		return s, err
	}
	{{- end }}

	{{ end -}}
	return s, nil
}

`))
}
