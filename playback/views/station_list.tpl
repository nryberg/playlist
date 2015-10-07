{{define "content"}}
<p> Station Table </p>
<table class="table table-striped">
  <tr>
    <th>Station</th>
    <th>Frequency</th>
    <th>Plays</th>
  </tr>
  {{range .}}
    <tr>
      <td>
        <a href='/station/{{.ID}}'>
          {{.Location}}
        </a>
      </td>
      <td>{{.Freq}}</td>
      <td>{{.Plays}}</td>
    </tr>
  {{end}}
</table>
{{end}}
