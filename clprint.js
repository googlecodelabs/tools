const process = require('process');
const fs = require('fs');
const merge = require('easy-pdf-merge');
const puppeteer = require('puppeteer');
const url = 'https://codelabs.developers.google.com/codelabs/'

var files=[];

const exit = () => {files.forEach((f) => {fs.unlinkSync(f)})};

const get = async (lab, first, last) => {
  let browser = await puppeteer.launch();
  for (i = first; i <= last; i++) {
    console.log('fetching step ' + i);
    let page = await browser.newPage();
    let file = lab + '.' + i + '.pdf';
    await page.goto(url+lab+'/index.html#'+i, {waitUntil: 'networkidle0'});
    let height = await page.evaluate(i => {return document.getElementsByTagName('google-codelab-step')[i].scrollHeight}, i);
    await page.pdf({path: file, height: height});
    files.push(file);
  }
  await browser.close();
  if (files.length > 1) {
    merge(files, lab+'.pdf', (err)=>{
      if (err) return console.log(err);
      else console.log('Successfully merged!');
    });
  }
};

process.on('exit', exit);
get(process.argv[2], process.argv[3], process.argv[4]);
