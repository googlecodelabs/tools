const process = require('process');
const fs = require('fs');
const merge = require('easy-pdf-merge');
const puppeteer = require('puppeteer');
const url = 'https://codelabs.developers.google.com/codelabs/'

const get = async (lab, first, last) => {
  const browser = await puppeteer.launch();
  var pdfFiles=[];
  for (i = first; i <= last; i++) {
    console.log('fetching step ' + i);
    const page = await browser.newPage();
    const file = lab + '.' + i + '.pdf';
    await page.goto(url+lab+'/index.html#' + i, {waitUntil: 'networkidle0'});
    let height = await page.evaluate(i => {return document.getElementsByTagName('google-codelab-step')[i].scrollHeight}, i);
    await page.pdf({path: file, height: height});
    pdfFiles.push(file);
  }
  if (pdfFiles.length > 1) {
    await merge(pdfFiles, lab+'.pdf', (err)=>{
      if (err) return console.log(err);
      else console.log('Successfully merged!');
    });
    pdfFiles.forEach((file) => {fs.unlinkSync(file)});
  }
  await browser.close();
};

get(process.argv[2], process.argv[3], process.argv[4]);
