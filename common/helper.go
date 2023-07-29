package common

import (
	"context"
	"fmt"
	"github.com/realmicro/realmicro/metadata"
	"net/http"
	"strings"
)

func RequestToContext(r *http.Request) context.Context {
	ctx := context.Background()
	md := make(metadata.Metadata)
	for k, v := range r.Header {
		fmt.Println(k, v)
		md[k] = strings.Join(v, ",")
	}
	return metadata.NewContext(ctx, md)
}
