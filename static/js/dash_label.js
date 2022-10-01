let datum;

fetch('/static/json/data.json')
    .then(res => res.json())
    .then(jsondata => datum = jsondata)
    .then(() => {


            const lastDataItem = datum[datum.length - 1];

            //total articles
            const article = lastDataItem.article;
            document.getElementById("article").textContent += ` ${article}`;

            //total text
            const text = lastDataItem.text
            document.getElementById("text").textContent += ` ${text}`;

            //total code
            const code = lastDataItem.code
            document.getElementById("code").textContent += ` ${code}`;

            //total
            const total = lastDataItem.total
            document.getElementById("total").textContent += total;

    })