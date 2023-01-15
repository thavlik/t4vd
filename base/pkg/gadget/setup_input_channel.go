package gadget

import (
	"fmt"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// SetupInputChannel creates handlers for the named
// input channel that retrieves labels/data from the
// given data reference.
func SetupInputChannel(
	mux *mux.Router,
	name string,
	ref *DataRef,
	log *zap.Logger,
) {
	// retrieve untransformed data by id
	// (data NOT transformed by this gadget)
	mux.HandleFunc(
		fmt.Sprintf("/input/%s/x", name),
		HandleGetOutputDataFromRef(
			ref,
			log,
		))

	// retrieve a specific label by id
	// (used to resolve parent labels)
	mux.HandleFunc(
		fmt.Sprintf("/input/%s/y", name),
		HandleGetOutputLabelFromRef(
			ref,
			log,
		))

	// sample random labels from the input channel
	// used by the gui of this gadget to gather a
	// stack of input examples for the user to label
	mux.HandleFunc(
		fmt.Sprintf("/sample/input/%s/y", name),
		HandleSampleOutputLabelsFromRef(
			ref,
			log,
		))
}
