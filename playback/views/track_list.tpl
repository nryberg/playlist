    {{define "content"}}
    <table class="table table-striped">
      <tr>
        <th>Artist</th>
        <th>Title</th>
      </tr>
      {{range .Tracks}}
        <tr>
          <td>{{.Track.Artist}}</td>
          <td>{{.Track.Title}}</td>
        </tr>{{end}}
    </table>{{end}}
