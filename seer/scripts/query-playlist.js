// Run with:
// docker run -i --init --cap-add=SYS_ADMIN --rm ghcr.io/puppeteer/puppeteer:latest node -e "$(cat query-user.js)"
const puppeteer = require('puppeteer');
const querystring = require('node:querystring');
(async () => {
    const url = process.env.INPUT_URL;
    if (!url) {
        throw new Error("missing INPUT_URL");
    }
    const startOfQuery = url.indexOf("?");
    const query = url.slice(startOfQuery + 1);
    const queryParts = querystring.parse(query);
    const id = queryParts['list'];

    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.goto(url);

    const titleEl = await page.waitForSelector('h3.ytd-playlist-panel-renderer a');
    const title = await page.$eval("h3.ytd-playlist-panel-renderer a", element => element.textContent);
    await titleEl.dispose();

    const channelEl = await page.waitForSelector('#publisher-container yt-formatted-string.publisher a');
    const channel = await page.$eval("#publisher-container yt-formatted-string.publisher a", element => element.textContent);
    const channelID = (await page.$eval("#publisher-container yt-formatted-string.publisher a", element => element.getAttribute('href'))).slice(1);
    const numVideos = Number(await page.$eval("#publisher-container > div.index-message-wrapper yt-formatted-string span:last-child", element => element.textContent));
    await channelEl.dispose();


    // output json must conform to seer API's PlaylistDetails
    console.log(JSON.stringify({
        id,
        title,
        channel,
        channelID,
        numVideos,
    }));

    await browser.close();
})();

