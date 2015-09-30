
{{define "content"}}
<table class="table table-striped">
  <tr>
    <th>Artist</th>
  </tr>
  {{range .Artists}}
    <tr>
      <td>{{.}}</td>
      <td>{{.Track.Title}}</td>
    </tr>
  {{end}}
</table>
{{end}}
