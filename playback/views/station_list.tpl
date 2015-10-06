{{define "content"}}
<p> Station Table </p>
<table class="table table-striped">
  <tr>
    <th>Station</th>
    <th>Plays</th>
  </tr>
  {{range .}}
    <tr>
      <td>
        <a href='/station/{{.StationID}}'>
          {{.Name}}
        </a>
      </td>
      <td>{{.Plays}}</td>
    </tr>
  {{end}}
</table>
{{end}}
