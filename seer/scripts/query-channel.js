// Run with:
// docker run -i --init --cap-add=SYS_ADMIN --rm ghcr.io/puppeteer/puppeteer:latest node -e "$(cat query-user.js)"
const puppeteer = require('puppeteer');
(async () => {
    const url = process.env.INPUT_URL;
    if (!url) {
        throw new Error("missing INPUT_URL");
    }

    const browser = await puppeteer.launch();
    try {
        const page = await browser.newPage();
        await page.goto(url);

        const avatarEl = await page.waitForSelector('#channel-header-container img');
        const avatar = await page.$eval("#channel-header-container img", element => element.getAttribute("src"));
        await avatarEl.dispose();

        const idEl = await page.waitForSelector('#channel-handle');
        const id = await page.$eval("#channel-handle", element => element.textContent);
        await idEl.dispose();

        const nameEl = await page.waitForSelector('.ytd-channel-name #text');
        const name = await page.$eval(".ytd-channel-name #text", element => element.textContent);
        await nameEl.dispose();

        const subsEl = await page.waitForSelector('#subscriber-count');
        let subs = await page.$eval("#subscriber-count", element => element.textContent);
        await subsEl.dispose();
        const parts = subs.split(" ");
        if (parts.length > 1) {
            subs = parts[0];
        }

        // output json must conform to seer API's ChannelDetails
        console.log(JSON.stringify({
            id,
            avatar,
            name,
            subs,
        }));
    } finally {
        await browser.close();
    }
})();

