    {{define "content"}}
    <table class="table table-striped">
      <tr>
        <th>Song</th>
        <th>Plays</th>
      </tr>
      {{range .}}
        <tr>
          <td><a href="/song/{{.SongID}}">{{.Title}}</a></td>
          <td>{{.Plays}}</td>
        </tr>{{end}}
    </table>{{end}}
