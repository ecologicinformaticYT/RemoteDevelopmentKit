const fs = require('fs');
const path = require('path');

function exists(path){
  return fs.existsSync(path);
}

function read(path){
  if (exists(path)){
    return fs.readFileSync(path, 'utf-8');
  } else {
    return null;
  }
}

function write(path, content){
    fs.writeFileSync(path, content);
}

function makedir(path){
    fs.mkdirSync(path);
}

function copyFile(path, destination){
    try {
        fs.copyFileSync(path, destination);
    } catch (err){
        console.log(err);
    }
}

function copyFolder(sourceDir, destinationDir) {
    try {
      // Read the content of the source folder
      const files = fs.readdirSync(sourceDir);
  
      // Read the content of the sub folders
      files.forEach(file => {
        const sourceFile = path.join(sourceDir, file);
        const destinationFile = path.join(destinationDir, file);
  
        try {
          // check if it's a folder
          const stats = fs.statSync(sourceFile);
  
          if (stats.isDirectory()) {
            // If it's a folder, create the destination folder & copy it's content
            fs.mkdirSync(destinationFile, { recursive: true });
            copyFolder(sourceFile, destinationFile);
          } else {
            // If it's a file, copy it
            fs.copyFileSync(sourceFile, destinationFile);
          }
        } catch (err) {
          console.error(`Error while copying ${sourceFile}: `, err);
        }
      });
    } catch (err) {
      console.error('Error while reading the source folder : ', err);
    }
}

function listChilds(parent){
  const result = fs.readdirSync(parent);

  return result
}

function listAllChildsSync(dir) {
  let results = [];

  function readDirRecursive(currentDir) {
      const files = fs.readdirSync(currentDir);

      for (const file of files) {
          const fullPath = path.join(currentDir, file);
          const stat = fs.statSync(fullPath);

          if (stat.isDirectory()) {
              results.push(fullPath);
              readDirRecursive(fullPath);
          } else {
              results.push(fullPath);
          }
      }
  }

  readDirRecursive(dir);
  return results;
}

function Delete(path){
  fs.unlinkSync(path);
}

function reWriteDat(path, toDeleteData){
  let toReWrite = '';
  let txt = read(path);
  txt = txt.split('\n');
  txt.forEach((line) => {
    if (line != toDeleteData){
      toReWrite += line+'\n';
    }
  });
  write(toReWrite);
}

function isDir(path){
  try {
    const stats = fs.statSync(path);
    return stats.isDirectory();
  } catch (err) {
    console.error(`Error while checking if ${path} is a folder : ${err}`);
    return false;
  }
}

module.exports = {
    read,
    write,
    makedir,
    exists,
    copyFile,
    copyFolder,
    listChilds,
    Delete,
    reWriteDat, 
    isDir,
    listAllChildsSync
}