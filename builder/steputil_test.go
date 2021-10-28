package builder

/*
func TestDefaultPrepareStepsTemplate(t *testing.T) {
	tests := []struct {
		name           string
		isIncludeCache bool
		want           *models.OptionNode
	}{
		{
			name:           "no cache",
			isIncludeCache: false,
		},
		{
			name:           "cache",
			isIncludeCache: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultPrepareStepsTemplate(tt.isIncludeCache)
			output, err := got.ExecuteAll(nil, nil)

			require.NoError(t, err)

			gotOrig := steps.DefaultPrepareStepList(tt.isIncludeCache)
			require.Equal(t, gotOrig, output.EvaluatedTemplate.Steps)
		})
	}
}*/
