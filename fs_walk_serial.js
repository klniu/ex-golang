var fs = require('fs');

var walk = function(baseDir, dir, done) {
  var results = [];
  
  fs.readdirSync(baseDir + dir, function(err, list) {
    if (err) return done(err);
    var i = 0;
    
    (function next() {
      var file = list[i++];
      if (!file) return done(null, results);
      var path = dir + '/' + file;
      var filePath = baseDir + '/' + path;
      fs.statSync(filePath, function(err, stat) {
        if (stat && stat.isDirectory()) {
          walk(baseDir, path, function(err, res) {
            results = results.concat(res);
            next();
          });
        } else {
          results.push(path);
          next();
        }
      });
    })();

  });
};

console.log("begin");

walk("D:/home/gocode/src/github.com/zhgo/example/frontend/src", "", function(err, results) {
  if (err) throw err;
  console.log(results);
});

/*

var promise = new Promise(function(resolve, reject) {
    walk("D:/home/gocode/src/github.com/zhgo/example/frontend/src", "", function(err, entries) {
        if (err) {
            reject(err);
        } else {
            resolve(entries);
        }
    });
});

promise
.then(function(entries){
  console.log(entries);
})
.catch(function(err){
  console.log(err);
});

*/

console.log("finiish");

