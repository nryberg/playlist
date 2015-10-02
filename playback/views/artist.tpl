{{define "content"}}
  <h2> Artist Information</h2>
  <div class="row">
    <div class="col-md-2"/><strong>ID</strong></div>
    <div class="col-md-10">{{.ArtistID}}</div>
  </div>
  <div class="row">
    <div class="col-md-2"/><strong>Name</strong></div>
    <div class="col-md-10">{{.Name}}</div>
  </div>

{{end}}
