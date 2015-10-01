{{define "content"}}
<p> Artist Table </p>
<table class="table table-striped">
  <tr>
    <th>Artist</th>
    <th>Plays</th>
  </tr>
  {{range .}}
    <tr>
      <td>
        <a href='/artist/{{.ArtistID}}'>
          {{.Name}}
        </a>
      </td>
      <td>{{.Plays}}</td>
    </tr>
  {{end}}
</table>
{{end}}
