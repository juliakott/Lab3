'use strict';
const fs = require('fs');
const crypto = require('crypto');

const args = process.argv.slice(2);

const inputDir = args[0];
const outputDir = args[1];

//const inputDir = './input';
//const outputDir = './output';
const readWriteAsync = (dirIn, dirOut) => {
  if (!fs.existsSync(dirIn)) return null;
  if (!fs.existsSync(dirOut)) fs.mkdirSync(dirOut);

  fs.readdir(dirIn, 'utf-8',(err, files) => {
    if(err) throw err;
    Promise.all(files.map(file => {
      return new Promise((resolve,reject) => {
        const md5sum = crypto.createHash('md5');
        const stream = fs.createReadStream(`./${dirIn}/${file}`, err => { if(err) reject(err);});
        stream
          .on('data', chunk => {
            md5sum.update(chunk);
          })
          .on('end',() => {
            fs.writeFile(`./${dirOut}/${file}.res`, md5sum.digest('hex'),'utf-8', err => {if(err) reject(err);});
          });
        resolve(true);
      });
    }))
      .then(res => {
        console.log(`Total number of processed files: ${res.length}.`);
      }, err => { if(err) console.log(err.stack); });
  });
};

readWriteAsync(inputDir, outputDir);
