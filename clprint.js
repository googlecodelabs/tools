const process = require('process');
const fs = require('fs');
const merge = require('easy-pdf-merge');
const puppeteer = require('puppeteer');
const url = 'https://codelabs.developers.google.com/codelabs/'

/*
async function scrollToBottom() {
  await new Promise(resolve => {
    const distance = 100; // should be less than or equal to window.innerHeight
    const delay = 100;
    const timer = setInterval(() => {
      document.scrollingElement.scrollBy(0, distance);
      if (document.scrollingElement.scrollTop + window.innerHeight >= document.scrollingElement.scrollHeight) {
        clearInterval(timer);
        resolve();
      }
    }, delay);
  });
}
*/

const get = async (lab, first, last) => {
  const browser = await puppeteer.launch();
  var pdfFiles=[];
  for (i = first; i <= last; i++) {
    console.log('fetching step ' + i);
    const page = await browser.newPage();
    await page.setViewport({ width: 1920, height: 1080 });
    const file = lab + '.' + i + '.pdf';
    await page.goto(url+lab+'/index.html#' + i, {waitUntil: 'networkidle0'});
    //await page.evaluate(scrollToBottom);
    //await page.waitFor(1000);
await page.addStyleTag({
    content: '@page { size: auto; }',
})
    let height = await page.evaluate(() => document.documentElement.offsetHeight);
    console.log(height);
    await page.pdf({path: file, fullPage: true}); //height: height+'px'});
    pdfFiles.push(file);
  }
  if (pdfFiles.length > 1) {
    merge(pdfFiles, lab+'.pdf', (err)=>{
      if (err) return console.log(err);
      else console.log('Successfully merged!');
    });
    pdfFiles.forEach((file) => {fs.unlinkSync(file)});
  }
  await browser.close();
};

get(process.argv[1], process.argv[2], process.argv[3]);
