<!DOCTYPE html>
<html>
<head>
  <title>Playlist - Playback</title>
  <link rel="stylesheet" href="/static/css/bootstrap.css">
</head>
<body>
<div class="container">
  <h3>Tracks Sample</h3>
<p>{{.Timestamp}}</p>
<h4>{{.Station.Location}}</h4>

{{template "tracks" .}}
<!--
<form action="/new" method="post">
<p><label for="title">Title</label>
<input type="text" name="title" required /></p>
<p><label for="content">Content</label>
<textarea name="content" required></textarea></p>
<p><input type="submit" value="Create Pastebin" /></p>
</form>
</div>
<ul>
{{range .}}
<li><a href="/paste/{{.Id}}">{{.Title}}</a></li>
{{end}}
--> 
</ul>
</body>
</html>
