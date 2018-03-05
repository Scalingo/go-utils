package publicgen

import "reflect"

type GeneratorParams struct {
	OutputFile    string
	OutputPackage string
	Types         []interface{}
}

func GeneratePublicModels(params GeneratorParams) error {
	astFields := make(map[string][]field)

	for _, typ := range params.Types {
		st := reflect.TypeOf(typ)
		fields := listFields(st)
		astFields[st.Name()] = fields
	}

	ast := newAST(params.OutputPackage, astFields)
	err := writeAst(ast, params.OutputFile)
	return err
}
