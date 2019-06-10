
	package predicator

	import (
		"testing"
	)

	// Test_load_boolean -- AUTOGENERATED DO NOT EDIT
	func Test_load_boolean(t *testing.T) {
		instructions := [][]interface {}{[]interface {}{"load", "bool_var"}, []interface {}{"to_bool"}}
		tt := []struct{
			Name string
			Result bool
			Data   map[string]interface{} 
		}{
		
			{
				"with_no_context",
				false,
				map[string]interface {}(nil),
			},
		
			{
				"with_false",
				false,
				map[string]interface {}{"bool_var":false},
			},
		
			{
				"with_true",
				true,
				map[string]interface {}{"bool_var":true},
			},
		
		}
		for _, test := range tt {
			e := NewEvaluator(instructions, test.Data)
			got := e.result()
			if  got != test.Result {
				t.Logf("FAILED %s_%s expected %v got %v", "load_boolean", test.Name, test.Result, got)
				t.Fail()
			}
			if e.stack.count > 0 {
				t.Logf("FAILED %s_%s expected empty stack",  "load_boolean", test.Name)
				t.Fail()
			}
		}
	}
	