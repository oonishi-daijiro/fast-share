const express = require('express');
const fs = require('fs');
const path = require('path');
const bodyParser = require('body-parser');
const argv = require('argv');
const option = require('./options.json');
const readline = require('readline');
let publishingDirectory
let server = express()
server.use(bodyParser.urlencoded({
  extended: true
}))
server.listen(80)

server.get('/', (req, res) => {
  res.write("This is the http file dister")
  res.end()
});

function parseArg() {
  return argv.option(option).run()
}

function getCurrentTime() {
  let now = new Date();
  return `${now.getFullYear()}/${now.getMonth()}/${now.getDate()}  ${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`
}

(() => {
  if (fs.existsSync(parseArg().options.directory) && (fs.statSync(parseArg().options.directory).isDirectory())) {
    console.log(path.dirname(parseArg().options.directory));
    publishingDirectory = parseArg().options.directory
    console.log(`${getCurrentTime()}: Publish: ${parseArg().options.directory}`);
    require('dns').lookup(require('os').hostname(), function (err, add) {
      console.log(`ipv4: ${add}`);
    })
  } else {
    console.log(`No such as directory that you specfied "${parseArg().options.directory}"`);
    process.exit()
  }
})()

server.post('/', (req, res) => {
  console.log(req.headers['user-agent'])
  if (req.body.state === 'requestDirInfo') {
    console.log("Requested connection form", req.ip);
    reqAnserToUser("Do you accept?[y/n]").then((ans) => {
      if (ans === "y") {
        const dirName = path.basename(req.body.dirName)
        console.log(`${getCurrentTime()} : Receive the directory request "${req.body.dirName} from ${req.ip}"`);
        if (!searchDir(dirName)) {
          res.status(404)
          res.write('You requested the directory "' + dirName + '", but no such directory.')
          res.end()
          return
        }
        const pathes = readFilePathOfDir(dirName).map(e => {
          return path.relative(publishingDirectory, e)
        })
        console.log(getDirSize(pathes));
        res.json({
          pathes: pathes,
          packageName: dirName,
        })
        res.end()
      } else if (ans === "n") {
        res.status(403)
        res.write("Access denied.")
        res.end()
        return
      }
    })

  } else if (req.body.state === 'requestDirData') {
    const filePath = `${publishingDirectory}/${req.body.dirName}`
    const fileStream = streamFile(filePath, res)
    fileStream.then(() => {
      res.end()
    })
    fileStream.catch(err => {
      console.log(err);
    })
  } else {
    res.status(404)
    res.write("Invalid POST request.")
    res.end()
  }
})

function reqAnserToUser(question) {
  question = question + "\n"
  const stdinReader = readline.createInterface({
    input: process.stdin,
    output: process.stdout
  })
  return new Promise((resolve, reject) => {
    stdinReader.question(question, (ans) => {
      resolve(ans)
      stdinReader.close()
    })
  })
}

function searchDir(name) {
  const content = fs.readdirSync(publishingDirectory)
  const target = content.find((e) => {
    if (e === name) {
      return e
    }
  })
  return target ? publishingDirectory + "/" + target : false
}

function streamFile(path, res) {
  return new Promise((resolve, reject) => {
    const stream = fs.createReadStream(path)
    stream.on('data', (chunk) => {
      res.write(chunk)
    })
    stream.on('error', (err) => {
      console.log(err);
      reject()
    })
    stream.on('end', () => {
      resolve()
    })
  })
}

function getDirSize(pathes) {
  let totalSize = 0
  pathes.forEach(eachPath => {
    totalSize += fs.statSync(`${publishingDirectory}/${eachPath}`).size
  });
  return totalSize
}

function readFilePathOfDir(target) {
  const contents = []
  const packagePath = searchDir(target)

  function isDir(path) {
    try {
      const stat = fs.statSync(path)
      return stat.isDirectory() ? true : false
    } catch (err) {
      console.log(err);
      return 'error'
    }
  }

  function getEachPath(path) {
    const dir = fs.readdirSync(path)
    dir.forEach((e) => {
      if (isDir(path + '/' + e)) {
        getEachPath(path + '/' + e)
        return
      } else {
        contents.push(path + '/' + e)
      }
    })
  }
  getEachPath(packagePath)
  return contents
}
