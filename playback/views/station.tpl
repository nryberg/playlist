{{define "content"}}
  <h2> Station Information</h2>
  <div class="row">
    <div class="col-md-2"/><strong>ID</strong></div>
    <div class="col-md-10">{{.ID}}</div>
  </div>
  <div class="row">
    <div class="col-md-2"/><strong>Location</strong></div>
    <div class="col-md-10">{{.Location}}</div>
  </div>
  <div class="row">
    <div class="col-md-2"/><strong>Frequency</strong></div>
    <div class="col-md-10">{{.Freq}}</div>
  </div>
{{end}}
