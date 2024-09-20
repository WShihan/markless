package page

import (
	"embed"
	"fmt"
	"net/http"
	"strings"

	"markless/assets"
	"markless/util"

	"github.com/julienschmidt/httprouter"
)

func AssetsFinder(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	assettype := params.ByName("assettype")
	fileName := params.ByName("filename")
	assetDir := fmt.Sprintf("static/%s/%s", assettype, fileName)
	var fs embed.FS
	if strings.HasSuffix(fileName, ".js") || strings.HasSuffix(fileName, ".map") {
		fs = assets.JS
		w.Header().Set("Content-Type", "text/js")
	} else if strings.HasSuffix(fileName, ".png") {
		fs = assets.IMG
		w.Header().Set("Content-Type", "image/x-icon")
	} else if strings.HasSuffix(fileName, ".ico") {
		fs = assets.IMG
		w.Header().Set("Content-Type", "image/x-icon")
	} else {
		fs = assets.CSS
		w.Header().Set("Content-Type", "text/css")
	}
	content, err := fs.ReadFile(assetDir)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	w.Write(content)
}
