
{{define "content"}}
<table class="table table-striped">
  <tr>
    <th>Artist</th>
  </tr>
  {{ range $key, $value := . }}
   <!-- <li><strong>{{ $key }}</strong>: {{ $value }}</li> -->
    <tr>
      <td>{{$key}}</td>
      <td>{{$value}}</td>
    </tr>
{{ end }}
<!--
  {{range .Artists}}
    <tr>
      <td>{{.}}</td>
      <td>{{.Track.Title}}</td>
    </tr>
  {{end}}
  -->
</table>
{{end}}
