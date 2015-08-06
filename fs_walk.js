var fs = require('fs');
var path = require('path');

var scriptsPath = 'D:/home/gocode/src/github.com/zhgo/example/frontend/src';

function getFolders(dir) {
  return fs.readdirSync(dir)
    .filter(function(file) {
      return fs.statSync(path.join(dir, file)).isDirectory();
    });
}

var folders = getFolders(scriptsPath);

var files = folders.map(function(folder) {
  console.log(folder);

  var files = fs.readdirSync(path.join(scriptsPath, folder));
  for (var i in files) {
    console.log(path.join(folder, files[i]));
  }
});