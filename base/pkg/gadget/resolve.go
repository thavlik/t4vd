package gadget

import "fmt"

func ResolveBaseURL(gadgetName string) string {
	return fmt.Sprintf("http://%s:80", gadgetName)
}
