package step

// Error is an error occuring top level in a step
type Error struct {
	StepID, Tag, ShortMsg string
	Err                   error
	Recommendations       Recommendation
}

// BranchReccomendations ...
type BranchReccomendations struct {
	AvailableBranches []string
}

// Recommendation interface
type Recommendation interface {
	GetRecommendationID() string
}

// GetRecommendationID ...
func (br BranchReccomendations) GetRecommendationID() string {
	return "branchRecommendation"
}

// NewError constructs a step.Error
func NewError(stepID, tag string, err error, shortMsg string) *Error {
	return &Error{
		StepID:   stepID,
		Tag:      tag,
		Err:      err,
		ShortMsg: shortMsg,
	}
}

// NewErrorWithRecommendations constructs a step.Error
func NewErrorWithRecommendations(stepID, tag string, err error, shortMsg string, recommendations Recommendation) *Error {
	return &Error{
		StepID:          stepID,
		Tag:             tag,
		Err:             err,
		ShortMsg:        shortMsg,
		Recommendations: recommendations,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
