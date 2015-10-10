{{define "content"}}
  <h2>Song Information</h2>
  <div class="row">
    <div class="col-md-2"/><strong>ID</strong></div>
    <div class="col-md-10">{{.SongID}}</div>
  </div>
  <div class="row">
    <div class="col-md-2"/><strong>Title</strong></div>
    <div class="col-md-10">{{.Title}}</div>
  </div>
{{end}}
