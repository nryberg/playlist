<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Playlist - Playback</title>
  <link rel="stylesheet" href="/static/css/bootstrap.css">
  <style>
    body { padding-top: 60px; }
  </style>

</head>
<body >
  <nav class="navbar navbar-inverse navbar-fixed-top">
    <div class="container">
      <div class="navbar-header">
        <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
          <span class="sr-only">Toggle navigation</span>
          <span class="icon-bar"></span>
          <span class="icon-bar"></span>
          <span class="icon-bar"></span>
        </button>
        <a class="navbar-brand" href="#">Playback for Playlist</a>
      </div>
      <div id="navbar" class="collapse navbar-collapse">
        <ul class="nav navbar-nav">
          <li class="active"><a href="/songs">Songs</a></li>
          <li><a href="/artists">Artists</a></li>
          <li><a href="/stations">Stations</a></li>
          <li><a href="#">About</a></li>
        </ul>
      </div><!--/.nav-collapse -->
    </div>
  </nav>



  <div class="container">
     {{template "content" .}}
  </div>
</body>
</html>
