package stats

import (
	"encoding/json"
	"fmt"
	statsapi "github.com/fukata/golang-stats-api-handler"
	"html/template"
	"net/http"
	"raycat/internal/pkg/stats/templates"
	"strconv"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	stats := statsapi.GetStats()

	jsonData, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	memoryFields := []string{"memory_alloc", "memory_total_alloc", "memory_sys", "heap_alloc", "heap_sys", "heap_idle", "heap_inuse", "heap_released"}
	for _, field := range memoryFields {
		if value, ok := data[field].(float64); ok {
			data[field] = formatBytes(value)
		}
	}

	intFields := []string{"memory_lookups", "memory_mallocs", "memory_frees", "heap_objects"}
	for _, field := range intFields {
		if value, ok := data[field].(float64); ok {
			data[field] = addCommas(strconv.FormatInt(int64(value), 10))
		}
	}
	if value, ok := data["gc_next"].(float64); ok {
		data["gc_next"] = formatBytes(value)
	}
	if value, ok := data["gc_last"].(float64); ok {
		t := time.Unix(0, int64(value))
		data["gc_last"] = t.Format("2006-01-02 15:04:05")
	}
	if gcPause, ok := data["gc_pause"].([]interface{}); ok {
		formattedPauses := make([]string, len(gcPause))
		for i, pause := range gcPause {
			if p, ok := pause.(float64); ok {
				formattedPauses[i] = formatDuration(time.Duration(p * float64(time.Second)))
			}
		}
		data["gc_pause"] = formattedPauses
	}
	tmpl, err := template.New("stats.html").Funcs(template.FuncMap{
		"formatBytes": formatBytes,
		"addCommas":   addCommas,
	}).ParseFS(templates.TpFs, "stats.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func formatBytes(bytes float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	for ; bytes >= 1024 && i < len(units)-1; i++ {
		bytes /= 1024
	}
	return fmt.Sprintf("%.2f %s", bytes, units[i])
}

func addCommas(s string) string {
	n, _ := strconv.Atoi(s)
	in := strconv.FormatInt(int64(n), 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
	if in[0] == '-' {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%.2f ns", float64(d.Nanoseconds()))
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.2f µs", float64(d.Microseconds()))
	} else if d < time.Second {
		return fmt.Sprintf("%.2f ms", float64(d.Milliseconds()))
	}
	return d.String()
}
