package interceptor

//go:generate embed file --source template.tmpl.go --var template
//go:generate embed file --source template_test.tmpl.go --var templateTest
var (
    template string
    templateTest string
)
